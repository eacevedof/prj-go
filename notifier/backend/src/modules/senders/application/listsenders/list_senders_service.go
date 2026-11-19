package listsenders

import (
    "notifier-backend/src/core/config"

    usersrepositories "notifier-backend/src/modules/users/infrastructure/repositories"

    "notifier-backend/src/modules/senders/domain/entities"
    "notifier-backend/src/modules/senders/domain/enums"
    "notifier-backend/src/modules/senders/domain/exceptions"
    "notifier-backend/src/modules/senders/infrastructure/repositories"
)

// ListSendersService los perfiles de envio de un usuario, o todos (./console sender-list).
type ListSendersService struct {
    usersReaderSqliteRepository   *usersrepositories.UsersReaderSqliteRepository
    sendersReaderSqliteRepository *repositories.SendersReaderSqliteRepository
}

func ListSendersServiceGetInstance() *ListSendersService {
    return &ListSendersService{
        usersReaderSqliteRepository:   usersrepositories.UsersReaderSqliteRepositoryGetInstance(),
        sendersReaderSqliteRepository: repositories.SendersReaderSqliteRepositoryGetInstance(),
    }
}

// Invoke error: *SendersException (404) o uno inesperado.
func (s *ListSendersService) Invoke(listSendersDto *ListSendersDto) (*ListSendersResultDto, error) {
    userId := 0
    if userUuid := listSendersDto.GetUserUuid(); userUuid != "" {
        userEntity, err := s.usersReaderSqliteRepository.GetUserByUserUuid(userUuid)
        if err != nil {
            return nil, err
        }
        if userEntity == nil {
            return nil, exceptions.NotFoundCustom(config.Tr(enums.SendersMessageEnum.USER_X_NOT_FOUND, userUuid))
        }
        userId = userEntity.GetId()
    }

    senderEntities, err := s.sendersReaderSqliteRepository.GetSenders(userId)
    if err != nil {
        return nil, err
    }
    senders := make([]entities.SenderPrimitives, 0, len(senderEntities))
    for _, senderEntity := range senderEntities {
        senders = append(senders, senderEntity.ToPrimitives())
    }
    return ListSendersResultDtoFromPrimitives(map[string]any{"senders": senders}), nil
}
