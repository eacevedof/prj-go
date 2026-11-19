package loginuser

import (
    "notifier-backend/src/modules/shared/application"

    "notifier-backend/src/modules/users/domain/entities"
)

// LogInUserPrimitives la respuesta del log-in (claves en este orden).
type LogInUserPrimitives struct {
    User     entities.UserPrimitives `json:"user"`
    ApiToken string                  `json:"api_token"`
}

type LogInUserResultDto struct {
    user     entities.UserPrimitives
    apiToken string
}

func LogInUserResultDtoFromPrimitives(primitives map[string]any) *LogInUserResultDto {
    user, _ := primitives["user"].(entities.UserPrimitives)
    return &LogInUserResultDto{user: user, apiToken: application.GetTrimmedString(primitives, "apiToken")}
}

func (d *LogInUserResultDto) ToArray() LogInUserPrimitives {
    return LogInUserPrimitives{User: d.user, ApiToken: d.apiToken}
}
