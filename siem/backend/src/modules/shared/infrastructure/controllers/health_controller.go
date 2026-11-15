package controllers

import (
    "database/sql"
    "runtime"

    "siem-backend/src/modules/shared/domain/enums"
    "siem-backend/src/modules/shared/infrastructure/components/httprequester"
    "siem-backend/src/modules/shared/infrastructure/components/httpresponder"
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
