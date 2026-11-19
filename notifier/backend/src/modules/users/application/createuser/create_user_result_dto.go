package createuser

import "notifier-backend/src/modules/users/domain/entities"

type CreateUserResultDto struct {
    user entities.UserPrimitives
}

func CreateUserResultDtoFromPrimitives(primitives map[string]any) *CreateUserResultDto {
    user, _ := primitives["user"].(entities.UserPrimitives)
    return &CreateUserResultDto{user: user}
}

func (d *CreateUserResultDto) ToArray() entities.UserPrimitives {
    return d.user
}
