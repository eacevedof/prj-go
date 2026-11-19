package createsender

import "notifier-backend/src/modules/shared/application"

type CreateSenderDto struct {
    application.AbstractInputDto
    userUuid     string
    senderName   string
    providerCode string
    fromEmail    string
    fromName     string
    smtpHost     string
    smtpPort     int
    smtpSecurity string
    username     string
    secret       string
    dailyLimit   int
    minuteLimit  int
    tokenName    string
}

func CreateSenderDtoFromPrimitives(primitives map[string]any) *CreateSenderDto {
    return &CreateSenderDto{
        AbstractInputDto: application.NewAbstractInputDto(primitives),
        userUuid:         application.GetTrimmedString(primitives, "userUuid"),
        senderName:       application.GetTrimmedString(primitives, "senderName"),
        providerCode:     application.GetTrimmedString(primitives, "providerCode"),
        fromEmail:        application.GetTrimmedString(primitives, "fromEmail"),
        fromName:         application.GetTrimmedString(primitives, "fromName"),
        smtpHost:         application.GetTrimmedString(primitives, "smtpHost"),
        smtpPort:         application.GetInt(primitives, "smtpPort"),
        smtpSecurity:     application.GetTrimmedString(primitives, "smtpSecurity"),
        username:         application.GetTrimmedString(primitives, "username"),
        // el secreto no se recorta: una app password puede llevar espacios (las de Gmail se muestran en 4 grupos)
        secret:      application.ToString(primitives["secret"]),
        dailyLimit:  application.GetInt(primitives, "dailyLimit"),
        minuteLimit: application.GetInt(primitives, "minuteLimit"),
        tokenName:   application.GetTrimmedString(primitives, "tokenName"),
    }
}

func (d *CreateSenderDto) GetUserUuid() string {
    return d.userUuid
}

func (d *CreateSenderDto) GetSenderName() string {
    return d.senderName
}

func (d *CreateSenderDto) GetProviderCode() string {
    return d.providerCode
}

func (d *CreateSenderDto) GetFromEmail() string {
    return d.fromEmail
}

func (d *CreateSenderDto) GetFromName() string {
    return d.fromName
}

func (d *CreateSenderDto) GetSmtpHost() string {
    return d.smtpHost
}

func (d *CreateSenderDto) GetSmtpPort() int {
    return d.smtpPort
}

func (d *CreateSenderDto) GetSmtpSecurity() string {
    return d.smtpSecurity
}

func (d *CreateSenderDto) GetUsername() string {
    return d.username
}

func (d *CreateSenderDto) GetSecret() string {
    return d.secret
}

func (d *CreateSenderDto) GetDailyLimit() int {
    return d.dailyLimit
}

func (d *CreateSenderDto) GetMinuteLimit() int {
    return d.minuteLimit
}

func (d *CreateSenderDto) GetTokenName() string {
    return d.tokenName
}
