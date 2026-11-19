package enqueuemessage

import (
    "strings"

    "notifier-backend/src/core/config"

    sharedenums "notifier-backend/src/modules/shared/domain/enums"
    "notifier-backend/src/modules/shared/infrastructure/components/datetimer"
    "notifier-backend/src/modules/shared/infrastructure/components/uuid"
    sharedexceptions "notifier-backend/src/modules/shared/infrastructure/exceptions"
    "notifier-backend/src/modules/shared/infrastructure/repositories/configuration"
    "notifier-backend/src/modules/shared/infrastructure/traits"

    sendersenums "notifier-backend/src/modules/senders/domain/enums"

    apitokensrepositories "notifier-backend/src/modules/apitokens/infrastructure/repositories"

    "notifier-backend/src/modules/messages/domain/entities"
    "notifier-backend/src/modules/messages/domain/enums"
    "notifier-backend/src/modules/messages/domain/exceptions"
    "notifier-backend/src/modules/messages/domain/services"
    "notifier-backend/src/modules/messages/infrastructure/repositories"
)

// canales con worker en esta version: el resto existe en el modelo pero aun no se entrega (roadmap)
var channelsWithDelivery = []string{sendersenums.ChannelEnum.EMAIL}

// EnqueueMessageService lo que hace POST /api/v1/messages, y nada mas: valida el mensaje y lo deja en app_messages
// como pending (con sus adjuntos). Enviarlo es de ProcessMessageQueueService (el worker del canal): el cliente
// recibe su 202 enseguida.
type EnqueueMessageService struct {
    traits.ErrorStackTrait

    dateTimer                       *datetimer.DateTimer
    uuidGenerator                   *uuid.UuidGenerator
    messagePayloadNormalizer        *services.MessagePayloadNormalizer
    messageFieldsValidator          *services.MessageFieldsValidator
    messagesReaderSqliteRepository  *repositories.MessagesReaderSqliteRepository
    messagesWriterSqliteRepository  *repositories.MessagesWriterSqliteRepository
    apiTokensWriterSqliteRepository *apitokensrepositories.ApiTokensWriterSqliteRepository

    enqueueMessageDto *EnqueueMessageDto
    messageFields     services.EmailMessageFields
}

func EnqueueMessageServiceGetInstance() *EnqueueMessageService {
    environmentReaderRawRepository := configuration.GetInstance()
    return &EnqueueMessageService{
        dateTimer:                datetimer.GetInstance(),
        uuidGenerator:            uuid.GetInstance(),
        messagePayloadNormalizer: services.MessagePayloadNormalizerGetInstance(),
        messageFieldsValidator: services.MessageFieldsValidatorGetInstance(
            services.AttachmentUrlCheckerFromAllowedHosts(environmentReaderRawRepository.GetAttachmentAllowedHosts()),
        ),
        messagesReaderSqliteRepository:  repositories.MessagesReaderSqliteRepositoryGetInstance(),
        messagesWriterSqliteRepository:  repositories.MessagesWriterSqliteRepositoryGetInstance(),
        apiTokensWriterSqliteRepository: apitokensrepositories.ApiTokensWriterSqliteRepositoryGetInstance(),
    }
}

// Invoke error: *MessagesException (400, 401), *FieldsException (400) o uno inesperado.
func (s *EnqueueMessageService) Invoke(enqueueMessageDto *EnqueueMessageDto) (*EnqueueMessageResultDto, error) {
    s.enqueueMessageDto = enqueueMessageDto

    if err := s.failIfWrongAuthOrChannel(); err != nil {
        return nil, err
    }
    if duplicatedMessageEntity, err := s.getDuplicatedMessage(); err != nil || duplicatedMessageEntity != nil {
        return s.getDuplicateResult(duplicatedMessageEntity, err)
    }
    if err := s.loadMessageFieldsOrFail(); err != nil {
        return nil, err
    }

    now := s.dateTimer.GetNowDtIntoTargetTz()
    messageEntity := s.getNewMessageEntity(now)
    err := s.messagesWriterSqliteRepository.RunInTransaction(func() error {
        messageId, err := s.messagesWriterSqliteRepository.InsertMessage(messageEntity)
        if err != nil {
            return err
        }
        return s.insertAttachments(messageId)
    })
    if err != nil {
        // dos peticiones a la vez con la misma idempotency_key: la segunda choca con el UNIQUE y devuelve la primera
        if s.isIdempotencyConflict(err) {
            return s.getDuplicateResult(s.getDuplicatedMessage())
        }
        return nil, err
    }
    if _, err := s.apiTokensWriterSqliteRepository.UpdateLastUsedAtById(enqueueMessageDto.GetAuthApiTokenId(), now); err != nil {
        return nil, err
    }

    return EnqueueMessageResultDtoFromPrimitives(map[string]any{
        "messageUuid":   messageEntity.GetMessageUuid(),
        "channel":       messageEntity.GetChannel(),
        "messageStatus": enums.MessageStatusEnum.PENDING,
    }), nil
}

func (s *EnqueueMessageService) failIfWrongAuthOrChannel() error {
    if s.enqueueMessageDto.GetAuthApiTokenId() < 1 {
        return exceptions.UnauthorizedCustom(config.Tr(enums.MessagesMessageEnum.INVALID_OR_MISSING_API_TOKEN))
    }
    channel := s.enqueueMessageDto.GetAuthChannel()
    for _, channelWithDelivery := range channelsWithDelivery {
        if channel == channelWithDelivery {
            return nil
        }
    }
    return exceptions.BadRequestCustom(config.Tr(enums.MessagesMessageEnum.CHANNEL_X_NOT_AVAILABLE_YET, channel))
}

func (s *EnqueueMessageService) getDuplicatedMessage() (*entities.MessageEntity, error) {
    idempotencyKey := s.enqueueMessageDto.GetIdempotencyKey()
    if idempotencyKey == "" {
        return nil, nil
    }
    return s.messagesReaderSqliteRepository.GetMessageByApiTokenIdAndIdempotencyKey(s.enqueueMessageDto.GetAuthApiTokenId(), idempotencyKey)
}

func (s *EnqueueMessageService) getDuplicateResult(messageEntity *entities.MessageEntity, err error) (*EnqueueMessageResultDto, error) {
    if err != nil || messageEntity == nil {
        return nil, err
    }
    return EnqueueMessageResultDtoFromPrimitives(map[string]any{
        "messageUuid":   messageEntity.GetMessageUuid(),
        "channel":       messageEntity.GetChannel(),
        "messageStatus": messageEntity.GetMessageStatus(),
        "isDuplicate":   1,
    }), nil
}

func (s *EnqueueMessageService) loadMessageFieldsOrFail() error {
    keys := enums.MessagesKeyEnum
    enqueueMessageDto := s.enqueueMessageDto
    errors := []string{}

    to, isToValid := s.messagePayloadNormalizer.GetEmailList(enqueueMessageDto.GetToJson())
    cc, isCcValid := s.messagePayloadNormalizer.GetEmailList(enqueueMessageDto.GetCcJson())
    bcc, isBccValid := s.messagePayloadNormalizer.GetEmailList(enqueueMessageDto.GetBccJson())
    attachments, areAttachmentsValid := s.messagePayloadNormalizer.GetAttachments(enqueueMessageDto.GetAttachmentsJson())
    emailListChecks := []struct {
        fieldName string
        isValid   bool
    }{{keys.TO, isToValid}, {keys.CC, isCcValid}, {keys.BCC, isBccValid}}
    for _, emailListCheck := range emailListChecks {
        if !emailListCheck.isValid {
            errors = append(errors, emailListCheck.fieldName+": "+config.Tr(enums.MessagesMessageEnum.MUST_BE_A_STRING_OR_A_LIST))
        }
    }
    if !areAttachmentsValid {
        errors = append(errors, keys.ATTACHMENTS+": "+config.Tr(enums.MessagesMessageEnum.MUST_BE_A_LIST))
    }
    // sin "to" en el payload, los destinatarios por defecto del token
    if len(to) == 0 {
        to, _ = s.messagePayloadNormalizer.GetEmailList(enqueueMessageDto.GetAuthDefaultToJson())
    }

    s.messageFields = services.EmailMessageFields{
        To:             to,
        Cc:             cc,
        Bcc:            bcc,
        ReplyTo:        strings.ToLower(enqueueMessageDto.GetReplyTo()),
        Subject:        enqueueMessageDto.GetSubject(),
        Body:           enqueueMessageDto.GetBody(),
        BodyFormat:     s.getBodyFormat(),
        IdempotencyKey: enqueueMessageDto.GetIdempotencyKey(),
        Attachments:    attachments,
    }
    environmentReaderRawRepository := configuration.GetInstance()
    errors = append(errors, s.messageFieldsValidator.GetEmailErrors(s.messageFields, services.MessageFieldsLimits{
        MaxRecipients:      environmentReaderRawRepository.GetMessageMaxRecipients(),
        AttachmentMaxBytes: environmentReaderRawRepository.GetAttachmentMaxBytes(),
    })...)
    if len(errors) == 0 {
        return nil
    }
    return sharedexceptions.FieldsExceptionBadRequestCustom(
        s.GetErrorStackMessage(config.Tr(sharedenums.SharedMessageEnum.FIELDS_VALIDATION_ERRORS), errors),
    )
}

func (s *EnqueueMessageService) getBodyFormat() string {
    if bodyFormat := strings.ToLower(s.enqueueMessageDto.GetBodyFormat()); bodyFormat != "" {
        return bodyFormat
    }
    return enums.BodyFormatEnum.TEXT
}

func (s *EnqueueMessageService) getNewMessageEntity(now string) *entities.MessageEntity {
    keys := enums.MessagesKeyEnum
    return entities.MessageEntityFromPrimitives(map[string]any{
        keys.MESSAGE_UUID: s.uuidGenerator.GetUniqueIdWithPrefix(entities.MESSAGE_UUID_PREFIX),
        keys.API_TOKEN_ID: s.enqueueMessageDto.GetAuthApiTokenId(),
        keys.SENDER_ID:    s.enqueueMessageDto.GetAuthSenderId(),
        keys.CHANNEL:      s.enqueueMessageDto.GetAuthChannel(),
        keys.RECIPIENTS: entities.Recipients{
            To:      s.messageFields.To,
            Cc:      s.messageFields.Cc,
            Bcc:     s.messageFields.Bcc,
            ReplyTo: s.messageFields.ReplyTo,
        },
        keys.SUBJECT:         s.messageFields.Subject,
        keys.BODY:            s.messageFields.Body,
        keys.BODY_FORMAT:     s.messageFields.BodyFormat,
        keys.MAX_ATTEMPTS:    enums.MessagesLimitEnum.MAX_ATTEMPTS,
        keys.IDEMPOTENCY_KEY: s.messageFields.IdempotencyKey,
        keys.REQUEST_IP:      s.enqueueMessageDto.GetUserIpAddress(),
        keys.CREATED_AT:      now,
    })
}

func (s *EnqueueMessageService) insertAttachments(messageId int) error {
    keys := enums.MessagesKeyEnum
    for _, attachment := range s.messageFields.Attachments {
        attachmentType := enums.AttachmentTypeEnum.TEXT
        if attachment.HasUrl {
            attachmentType = enums.AttachmentTypeEnum.URL
        }
        _, err := s.messagesWriterSqliteRepository.InsertAttachment(entities.MessageAttachmentEntityFromPrimitives(map[string]any{
            keys.MESSAGE_ID:      messageId,
            keys.ATTACHMENT_TYPE: attachmentType,
            keys.FILE_NAME:       attachment.FileName,
            keys.CONTENT:         attachment.Content,
            keys.SOURCE_URL:      attachment.Url,
        }))
        if err != nil {
            return err
        }
    }
    return nil
}

func (s *EnqueueMessageService) isIdempotencyConflict(err error) bool {
    return s.enqueueMessageDto.GetIdempotencyKey() != "" && strings.Contains(err.Error(), "UNIQUE constraint failed")
}
