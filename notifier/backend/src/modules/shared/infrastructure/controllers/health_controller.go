package controllers

import (
    "database/sql"
    "runtime"

    "notifier-backend/src/modules/shared/domain/enums"
    "notifier-backend/src/modules/shared/infrastructure/components/httprequester"
    "notifier-backend/src/modules/shared/infrastructure/components/httpresponder"
)

// HealthController GET /health
type HealthController struct {
    AbstractRestController
}

func HealthControllerGetInstance(httpRequester *httprequester.HttpRequester) *HealthController {
    return &HealthController{AbstractRestController: NewAbstractRestController(httpRequester)}
}

func (c *HealthController) Invoke() *httpresponder.HttpResponder {
    return c.RespondOk(map[string]any{
        enums.HealthKeyEnum.APP:         enums.HealthKeyEnum.APP_NAME,
        enums.HealthKeyEnum.GO:          runtime.Version(),
        enums.HealthKeyEnum.SQL_DRIVERS: sql.Drivers(),
    })
}
