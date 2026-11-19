package loginuser

import "notifier-backend/src/modules/shared/application"

type LogInUserDto struct {
    application.AbstractInputDto
    userEmail string
    password  string
}

func LogInUserDtoFromPrimitives(primitives map[string]any) *LogInUserDto {
    return &LogInUserDto{
        AbstractInputDto: application.NewAbstractInputDto(primitives),
        userEmail:        application.GetTrimmedString(primitives, "userEmail"),
        password:         application.GetTrimmedString(primitives, "password"),
    }
}

func (d *LogInUserDto) GetUserEmail() string {
    return d.userEmail
}

func (d *LogInUserDto) GetPassword() string {
    return d.password
}
