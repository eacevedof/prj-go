package getmessage

import (
    "notifier-backend/src/core/config"

    "notifier-backend/src/modules/messages/domain/entities"
    "notifier-backend/src/modules/messages/domain/enums"
    "notifier-backend/src/modules/messages/domain/exceptions"
    "notifier-backend/src/modules/messages/infrastructure/repositories"
)

// GetMessageService el estado de un mensaje, con sus adjuntos y sus intentos. Solo con el token que lo encolo: con otro
// es un 404, como si no existiera.
type GetMessageService struct {
    messagesReaderSqliteRepository *repositories.MessagesReaderSqliteRepository
}

func GetMessageServiceGetInstance() *GetMessageService {
    return &GetMessageService{messagesReaderSqliteRepository: repositories.MessagesReaderSqliteRepositoryGetInstance()}
}

// Invoke error: *MessagesException (401, 404) o uno inesperado.
func (s *GetMessageService) Invoke(getMessageDto *GetMessageDto) (*GetMessageResultDto, error) {
    if getMessageDto.GetAuthApiTokenId() < 1 {
        return nil, exceptions.UnauthorizedCustom(config.Tr(enums.MessagesMessageEnum.INVALID_OR_MISSING_API_TOKEN))
    }
    messageUuid := getMessageDto.GetMessageUuid()
    messageEntity, err := s.messagesReaderSqliteRepository.GetMessageByApiTokenIdAndMessageUuid(getMessageDto.GetAuthApiTokenId(), messageUuid)
    if err != nil {
        return nil, err
    }
    if messageEntity == nil {
        return nil, exceptions.NotFoundCustom(config.Tr(enums.MessagesMessageEnum.MESSAGE_X_NOT_FOUND, messageUuid))
    }

    attachmentEntities, err := s.messagesReaderSqliteRepository.GetAttachmentsByMessageId(messageEntity.GetId())
    if err != nil {
        return nil, err
    }
    attemptEntities, err := s.messagesReaderSqliteRepository.GetAttemptsByMessageId(messageEntity.GetId())
    if err != nil {
        return nil, err
    }

    attachments := make([]entities.MessageAttachmentPrimitives, 0, len(attachmentEntities))
    for _, attachmentEntity := range attachmentEntities {
        attachments = append(attachments, attachmentEntity.ToPrimitives())
    }
    attemptsLog := make([]entities.MessageAttemptPrimitives, 0, len(attemptEntities))
    for _, attemptEntity := range attemptEntities {
        attemptsLog = append(attemptsLog, attemptEntity.ToPrimitives())
    }
    return GetMessageResultDtoFromPrimitives(map[string]any{"message": messageEntity.ToPrimitives(attachments, attemptsLog)}), nil
}
