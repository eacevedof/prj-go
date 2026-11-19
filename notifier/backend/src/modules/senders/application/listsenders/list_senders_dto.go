package listsenders

import "notifier-backend/src/modules/shared/application"

type ListSendersDto struct {
    application.AbstractInputDto
    userUuid string
}

func ListSendersDtoFromPrimitives(primitives map[string]any) *ListSendersDto {
    return &ListSendersDto{
        AbstractInputDto: application.NewAbstractInputDto(primitives),
        userUuid:         application.GetTrimmedString(primitives, "userUuid"),
    }
}

// GetUserUuid "" = los de todos los usuarios.
func (d *ListSendersDto) GetUserUuid() string {
    return d.userUuid
}
