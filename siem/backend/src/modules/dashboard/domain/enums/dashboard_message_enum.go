package enums

// DashboardMessageEnum mensajes de error del dashboard ({0}... los rellena config.Tr).
var DashboardMessageEnum = struct {
    INVALID_OR_MISSING_DASHBOARD_TOKEN string
    PROJECT_X_NOT_FOUND                string
}{
    INVALID_OR_MISSING_DASHBOARD_TOKEN: "Invalid or missing dashboard token",
    PROJECT_X_NOT_FOUND:                "Project {0} not found",
}
