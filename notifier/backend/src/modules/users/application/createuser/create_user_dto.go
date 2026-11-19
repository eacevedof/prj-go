package createuser

import "notifier-backend/src/modules/shared/application"

type CreateUserDto struct {
    application.AbstractInputDto
    userName  string
    userEmail string
    userRole  string
    password  string
}

func CreateUserDtoFromPrimitives(primitives map[string]any) *CreateUserDto {
    return &CreateUserDto{
        AbstractInputDto: application.NewAbstractInputDto(primitives),
        userName:         application.GetTrimmedString(primitives, "userName"),
        userEmail:        application.GetTrimmedString(primitives, "userEmail"),
        userRole:         application.GetTrimmedString(primitives, "userRole"),
        password:         application.GetTrimmedString(primitives, "password"),
    }
}

func (d *CreateUserDto) GetUserName() string {
    return d.userName
}

func (d *CreateUserDto) GetUserEmail() string {
    return d.userEmail
}

func (d *CreateUserDto) GetUserRole() string {
    return d.userRole
}

func (d *CreateUserDto) GetPassword() string {
    return d.password
}
