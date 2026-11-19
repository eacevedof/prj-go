package processmessagequeue

import (
    "mime"
    "path"
    "slices"
    "strings"
    "time"

    "notifier-backend/src/core/config"

    sharedenums "notifier-backend/src/modules/shared/domain/enums"
    "notifier-backend/src/modules/shared/infrastructure/components/datetimer"
    "notifier-backend/src/modules/shared/infrastructure/components/secretcipher"
    sharedexceptions "notifier-backend/src/modules/shared/infrastructure/exceptions"
    "notifier-backend/src/modules/shared/infrastructure/repositories/configuration"
    "notifier-backend/src/modules/shared/infrastructure/traits"

    sendersentities "notifier-backend/src/modules/senders/domain/entities"
    sendersenums "notifier-backend/src/modules/senders/domain/enums"
    sendersrepositories "notifier-backend/src/modules/senders/infrastructure/repositories"

    messagesentities "notifier-backend/src/modules/messages/domain/entities"
    messagesenums "notifier-backend/src/modules/messages/domain/enums"
    messagesservices "notifier-backend/src/modules/messages/domain/services"
    messagesrepositories "notifier-backend/src/modules/messages/infrastructure/repositories"

    "notifier-backend/src/modules/delivery/domain/entities"
    "notifier-backend/src/modules/delivery/domain/enums"
    "notifier-backend/src/modules/delivery/domain/exceptions"
    "notifier-backend/src/modules/delivery/domain/services"
    "notifier-backend/src/modules/delivery/infrastructure/components/attachmentdownloader"
    "notifier-backend/src/modules/delivery/infrastructure/messagesenders"
)

// canales con worker en esta version (fase 1: email)
var ChannelsWithWorker = []string{sendersenums.ChannelEnum.EMAIL}

// ProcessMessageQueueService vacia la cola de un canal: cada mensaje pendiente con turno (FIFO) se reclama (sending),
// se entrega por el proveedor de su perfil y queda sent, vuelve a pending con su siguiente turno o queda failed.
// Cada intento deja su fila en app_message_attempts. Lo llaman el worker de la API y ./console message-queue-process.
// Un perfil deshabilitado no se toca (sus mensajes esperan) y uno que llego a su limite aplaza el mensaje sin gastar
// intento.
type ProcessMessageQueueService struct {
    traits.ErrorStackTrait

    dateTimer                      *datetimer.DateTimer
    secretCipher                   *secretcipher.SecretCipher
    retryScheduler                 *services.RetryScheduler
    smtpErrorClassifier            *services.SmtpErrorClassifier
    emailSenderFactory             *messagesenders.EmailSenderFactory
    attachmentDownloader           *attachmentdownloader.AttachmentDownloader
    sendersReaderSqliteRepository  *sendersrepositories.SendersReaderSqliteRepository
    sendersWriterSqliteRepository  *sendersrepositories.SendersWriterSqliteRepository
    messagesReaderSqliteRepository *messagesrepositories.MessagesReaderSqliteRepository
    messagesWriterSqliteRepository *messagesrepositories.MessagesWriterSqliteRepository

    processMessageQueueDto *ProcessMessageQueueDto
    senderEntitiesById     map[int]*sendersentities.SenderEntity
    sentMessages           int
    retriedMessages        int
    failedMessages         int
    postponedMessages      int
}

func ProcessMessageQueueServiceGetInstance() *ProcessMessageQueueService {
    environmentReaderRawRepository := configuration.GetInstance()
    return &ProcessMessageQueueService{
        dateTimer:           datetimer.GetInstance(),
        secretCipher:        secretcipher.GetInstance(),
        retryScheduler:      services.RetrySchedulerGetInstance(),
        smtpErrorClassifier: services.SmtpErrorClassifierGetInstance(),
        emailSenderFactory:  messagesenders.EmailSenderFactoryGetInstance(),
        attachmentDownloader: attachmentdownloader.AttachmentDownloaderGetInstance(
            messagesservices.AttachmentUrlCheckerFromAllowedHosts(environmentReaderRawRepository.GetAttachmentAllowedHosts()),
            environmentReaderRawRepository.GetAttachmentMaxBytes(),
            environmentReaderRawRepository.GetAttachmentTimeoutSeconds(),
        ),
        sendersReaderSqliteRepository:  sendersrepositories.SendersReaderSqliteRepositoryGetInstance(),
        sendersWriterSqliteRepository:  sendersrepositories.SendersWriterSqliteRepositoryGetInstance(),
        messagesReaderSqliteRepository: messagesrepositories.MessagesReaderSqliteRepositoryGetInstance(),
        messagesWriterSqliteRepository: messagesrepositories.MessagesWriterSqliteRepositoryGetInstance(),
        senderEntitiesById:             map[int]*sendersentities.SenderEntity{},
    }
}

// Invoke error: *FieldsException (400) o uno inesperado (la BD). Un mensaje que falla no es un error del servicio:
// se cuenta y se sigue con el siguiente.
func (s *ProcessMessageQueueService) Invoke(processMessageQueueDto *ProcessMessageQueueDto) (*ProcessMessageQueueResultDto, error) {
    s.processMessageQueueDto = processMessageQueueDto
    if err := s.failIfWrongInput(); err != nil {
        return nil, err
    }

    orphansRecovered, err := s.messagesWriterSqliteRepository.UpdateOrphanSendingAsPending(
        s.dateTimer.GetNowDtMinusSeconds(enums.DeliveryLimitEnum.ORPHAN_SENDING_SECONDS), s.dateTimer.GetNowDtIntoTargetTz(),
    )
    if err != nil {
        return nil, err
    }

    dueMessageEntities, err := s.messagesReaderSqliteRepository.GetDueMessagesByChannel(
        processMessageQueueDto.GetChannel(), s.dateTimer.GetNowDtIntoTargetTz(), processMessageQueueDto.GetMaxMessages(),
    )
    if err != nil {
        return nil, err
    }
    for _, messageEntity := range dueMessageEntities {
        if err := s.processMessage(messageEntity); err != nil {
            return nil, err
        }
    }

    return ProcessMessageQueueResultDtoFromPrimitives(map[string]any{
        "channel":           processMessageQueueDto.GetChannel(),
        "sentMessages":      s.sentMessages,
        "retriedMessages":   s.retriedMessages,
        "failedMessages":    s.failedMessages,
        "postponedMessages": s.postponedMessages,
        "orphansRecovered":  orphansRecovered,
        "dueMessages":       len(dueMessageEntities),
    }), nil
}

func (s *ProcessMessageQueueService) failIfWrongInput() error {
    errors := []string{}
    maxMessages := s.processMessageQueueDto.GetMaxMessages()
    if maxMessages < enums.DeliveryLimitEnum.MAX_MESSAGES_MIN || maxMessages > enums.DeliveryLimitEnum.MAX_MESSAGES_MAX {
        errors = append(errors, enums.DeliveryMessageEnum.MAX_MESSAGES_MUST_BE_1_TO_500)
    }
    if !slices.Contains(ChannelsWithWorker, s.processMessageQueueDto.GetChannel()) {
        errors = append(errors, config.Tr(enums.DeliveryMessageEnum.CHANNEL_X_HAS_NO_WORKER, s.processMessageQueueDto.GetChannel()))
    }
    if len(errors) == 0 {
        return nil
    }
    return sharedexceptions.FieldsExceptionBadRequestCustom(
        s.GetErrorStackMessage(config.Tr(sharedenums.SharedMessageEnum.FIELDS_VALIDATION_ERRORS), errors),
    )
}

// processMessage solo devuelve error si falla la BD; lo del proveedor queda anotado en el mensaje.
func (s *ProcessMessageQueueService) processMessage(messageEntity *messagesentities.MessageEntity) error {
    senderEntity, err := s.getSenderEntity(messageEntity.GetSenderId())
    if err != nil {
        return err
    }
    if senderEntity != nil && !senderEntity.IsEnabled() {
        return nil
    }
    if senderEntity != nil {
        if isPostponed, err := s.postponeIfSenderAtLimit(messageEntity, senderEntity); err != nil || isPostponed {
            return err
        }
    }

    now := s.dateTimer.GetNowDtIntoTargetTz()
    claimedMessages, err := s.messagesWriterSqliteRepository.UpdatePendingAsSendingById(messageEntity.GetId(), now)
    if err != nil || claimedMessages == 0 {
        return err
    }

    startedAt := time.Now()
    providerResponse, deliveryErr := s.deliverMessage(messageEntity, senderEntity)
    durationMs := int(time.Since(startedAt).Milliseconds())
    if deliveryErr == nil {
        return s.markAsSent(messageEntity, providerResponse, durationMs)
    }
    return s.markAfterFailure(messageEntity, s.smtpErrorClassifier.GetDeliveryException(deliveryErr), durationMs)
}

// getSenderEntity cada perfil se lee una vez por pasada. nil si ya no existe (el mensaje fallara).
func (s *ProcessMessageQueueService) getSenderEntity(senderId int) (*sendersentities.SenderEntity, error) {
    if senderEntity, isLoaded := s.senderEntitiesById[senderId]; isLoaded {
        return senderEntity, nil
    }
    senderEntity, err := s.sendersReaderSqliteRepository.GetSenderById(senderId)
    if err != nil {
        return nil, err
    }
    s.senderEntitiesById[senderId] = senderEntity
    return senderEntity, nil
}

// postponeIfSenderAtLimit los limites cuentan los envios correctos del perfil en el ultimo minuto y en las ultimas 24 h.
func (s *ProcessMessageQueueService) postponeIfSenderAtLimit(messageEntity *messagesentities.MessageEntity, senderEntity *sendersentities.SenderEntity) (bool, error) {
    postponeSeconds := 0
    if minuteLimit := senderEntity.GetEffectiveMinuteLimit(); minuteLimit > 0 {
        sentLastMinute, err := s.messagesReaderSqliteRepository.CountSentBySenderIdSince(senderEntity.GetId(), s.dateTimer.GetNowDtMinusSeconds(60))
        if err != nil {
            return false, err
        }
        if sentLastMinute >= minuteLimit {
            postponeSeconds = enums.DeliveryLimitEnum.MINUTE_LIMIT_POSTPONE_SECONDS
        }
    }
    if dailyLimit := senderEntity.GetEffectiveDailyLimit(); postponeSeconds == 0 && dailyLimit > 0 {
        sentLastDay, err := s.messagesReaderSqliteRepository.CountSentBySenderIdSince(senderEntity.GetId(), s.dateTimer.GetNowDtMinusSeconds(86400))
        if err != nil {
            return false, err
        }
        if sentLastDay >= dailyLimit {
            postponeSeconds = enums.DeliveryLimitEnum.DAILY_LIMIT_POSTPONE_SECONDS
        }
    }
    if postponeSeconds == 0 {
        return false, nil
    }

    _, err := s.messagesWriterSqliteRepository.UpdateNextAttemptAtById(
        messageEntity.GetId(), s.dateTimer.GetNowDtPlusSeconds(postponeSeconds), s.dateTimer.GetNowDtIntoTargetTz(),
    )
    if err == nil {
        s.postponedMessages++
    }
    return true, err
}

// deliverMessage monta el email (secreto descifrado, adjuntos descargados) y lo entrega. Devuelve la respuesta del
// proveedor o el error tal cual (lo clasifica SmtpErrorClassifier).
func (s *ProcessMessageQueueService) deliverMessage(messageEntity *messagesentities.MessageEntity, senderEntity *sendersentities.SenderEntity) (string, error) {
    if senderEntity == nil {
        return "", exceptions.PermanentCustom(enums.DeliveryMessageEnum.SENDER_NOT_FOUND)
    }
    recipients := messageEntity.GetRecipients()
    if len(recipients.GetEnvelopeRecipients()) == 0 {
        return "", exceptions.PermanentCustom(enums.DeliveryMessageEnum.NO_RECIPIENTS)
    }
    emailSender, err := s.emailSenderFactory.GetEmailSender(senderEntity.GetProviderCode())
    if err != nil {
        return "", err
    }
    smtpSettings, err := s.getSmtpSettings(senderEntity)
    if err != nil {
        return "", err
    }
    attachments, err := s.getOutgoingAttachments(messageEntity.GetId())
    if err != nil {
        return "", err
    }

    keys := sendersenums.SendersKeyEnum
    fromEmail := senderEntity.GetConfigString(keys.FROM_EMAIL)
    return emailSender.Send(entities.OutgoingEmail{
        FromEmail:   fromEmail,
        FromName:    senderEntity.GetConfigString(keys.FROM_NAME),
        To:          recipients.To,
        Cc:          recipients.Cc,
        Bcc:         recipients.Bcc,
        ReplyTo:     recipients.ReplyTo,
        Subject:     messageEntity.GetSubject(),
        Body:        messageEntity.GetBody(),
        IsHtml:      messageEntity.GetBodyFormat() == messagesenums.BodyFormatEnum.HTML,
        MessageId:   s.getEmailMessageId(messageEntity, fromEmail),
        Date:        time.Now(),
        Attachments: attachments,
    }, smtpSettings)
}

// getSmtpSettings un secreto que no se descifra (cambio APP_SECRETS_KEY) es permanente: hay que volver a meterlo.
func (s *ProcessMessageQueueService) getSmtpSettings(senderEntity *sendersentities.SenderEntity) (entities.SmtpSettings, error) {
    password := ""
    if secretEncrypted := senderEntity.GetSecretEncrypted(); secretEncrypted != "" {
        decryptedSecret, err := s.secretCipher.GetDecrypted(secretEncrypted)
        if err != nil {
            return entities.SmtpSettings{}, exceptions.PermanentCustom(err.Error())
        }
        password = decryptedSecret
    }
    keys := sendersenums.SendersKeyEnum
    return entities.SmtpSettings{
        Host:     senderEntity.GetConfigString(keys.SMTP_HOST),
        Port:     senderEntity.GetConfigInt(keys.SMTP_PORT),
        Security: senderEntity.GetConfigString(keys.SMTP_SECURITY),
        Username: senderEntity.GetConfigString(keys.USERNAME),
        Password: password,
    }, nil
}

func (s *ProcessMessageQueueService) getOutgoingAttachments(messageId int) ([]entities.OutgoingAttachment, error) {
    attachmentEntities, err := s.messagesReaderSqliteRepository.GetAttachmentsByMessageId(messageId)
    if err != nil {
        return nil, err
    }
    outgoingAttachments := make([]entities.OutgoingAttachment, 0, len(attachmentEntities))
    for i, attachmentEntity := range attachmentEntities {
        content, mimeType := []byte(attachmentEntity.GetContent()), s.getTextMimeType(attachmentEntity.GetFileName())
        if attachmentEntity.IsUrl() {
            if content, mimeType, err = s.attachmentDownloader.GetContent(attachmentEntity.GetSourceUrl(), attachmentEntity.GetFileName(), i+1); err != nil {
                return nil, err
            }
        }
        if _, err := s.messagesWriterSqliteRepository.UpdateAttachmentMetadataById(attachmentEntity.GetId(), mimeType, len(content)); err != nil {
            return nil, err
        }
        outgoingAttachments = append(outgoingAttachments, entities.OutgoingAttachment{
            FileName: attachmentEntity.GetFileName(),
            MimeType: mimeType,
            Content:  content,
        })
    }
    return outgoingAttachments, nil
}

// getTextMimeType un adjunto de texto: por su extension (nota.csv -> text/csv) y, si no dice nada, text/plain.
func (s *ProcessMessageQueueService) getTextMimeType(fileName string) string {
    if mimeType := mime.TypeByExtension(path.Ext(fileName)); mimeType != "" {
        mediaType, _, _ := mime.ParseMediaType(mimeType)
        return mediaType
    }
    return "text/plain"
}

// getEmailMessageId <message_uuid@dominio-del-remitente>: unico, y se guarda como provider_message_id (SMTP no
// devuelve otro).
func (s *ProcessMessageQueueService) getEmailMessageId(messageEntity *messagesentities.MessageEntity, fromEmail string) string {
    domain := "notifier.local"
    if atPosition := strings.LastIndex(fromEmail, "@"); atPosition >= 0 {
        domain = fromEmail[atPosition+1:]
    }
    return messageEntity.GetMessageUuid() + "@" + domain
}

func (s *ProcessMessageQueueService) markAsSent(messageEntity *messagesentities.MessageEntity, providerResponse string, durationMs int) error {
    now := s.dateTimer.GetNowDtIntoTargetTz()
    if _, err := s.messagesWriterSqliteRepository.UpdateAsSentById(messageEntity.GetId(), s.getEmailMessageId(messageEntity, s.getFromEmail(messageEntity)), now); err != nil {
        return err
    }
    if err := s.insertAttempt(messageEntity, messagesenums.AttemptStatusEnum.SENT, providerResponse, durationMs, now); err != nil {
        return err
    }
    if _, err := s.sendersWriterSqliteRepository.UpdateLastUsedAtById(messageEntity.GetSenderId(), now); err != nil {
        return err
    }
    s.sentMessages++
    return nil
}

func (s *ProcessMessageQueueService) markAfterFailure(messageEntity *messagesentities.MessageEntity, deliveryException *exceptions.DeliveryException, durationMs int) error {
    now := s.dateTimer.GetNowDtIntoTargetTz()
    attemptsDone := messageEntity.GetAttempts() + 1
    isFinal := deliveryException.IsPermanent() || s.retryScheduler.IsLastAttempt(attemptsDone, messageEntity.GetMaxAttempts())

    messageStatus, attemptStatus, nextAttemptAt := messagesenums.MessageStatusEnum.FAILED, messagesenums.AttemptStatusEnum.PERMANENT_ERROR, now
    if !deliveryException.IsPermanent() {
        attemptStatus = messagesenums.AttemptStatusEnum.TEMPORARY_ERROR
    }
    if !isFinal {
        messageStatus = messagesenums.MessageStatusEnum.PENDING
        nextAttemptAt = s.dateTimer.GetNowDtPlusSeconds(s.retryScheduler.GetDelaySeconds(attemptsDone))
    }

    if _, err := s.messagesWriterSqliteRepository.UpdateAfterFailureById(messageEntity.GetId(), messageStatus, deliveryException.Error(), nextAttemptAt, now); err != nil {
        return err
    }
    if err := s.insertAttempt(messageEntity, attemptStatus, deliveryException.Error(), durationMs, now); err != nil {
        return err
    }
    if _, err := s.sendersWriterSqliteRepository.UpdateLastErrorById(messageEntity.GetSenderId(), deliveryException.Error(), now); err != nil {
        return err
    }
    if isFinal {
        s.failedMessages++
        return nil
    }
    s.retriedMessages++
    return nil
}

func (s *ProcessMessageQueueService) insertAttempt(messageEntity *messagesentities.MessageEntity, attemptStatus string, providerResponse string, durationMs int, attemptedAt string) error {
    keys := messagesenums.MessagesKeyEnum
    _, err := s.messagesWriterSqliteRepository.InsertAttempt(messagesentities.MessageAttemptEntityFromPrimitives(map[string]any{
        keys.MESSAGE_ID:        messageEntity.GetId(),
        keys.ATTEMPT_NUMBER:    messageEntity.GetAttempts() + 1,
        keys.ATTEMPT_STATUS:    attemptStatus,
        keys.PROVIDER_RESPONSE: providerResponse,
        keys.DURATION_MS:       durationMs,
        keys.ATTEMPTED_AT:      attemptedAt,
    }))
    return err
}

func (s *ProcessMessageQueueService) getFromEmail(messageEntity *messagesentities.MessageEntity) string {
    if senderEntity := s.senderEntitiesById[messageEntity.GetSenderId()]; senderEntity != nil {
        return senderEntity.GetConfigString(sendersenums.SendersKeyEnum.FROM_EMAIL)
    }
    return ""
}
