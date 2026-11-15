package getdashboard

// GetDashboardPrimitives la respuesta de GET /api/v1/dashboard (claves en este orden).
type GetDashboardPrimitives struct {
    Server   map[string]any   `json:"server"`
    Projects []map[string]any `json:"projects"`
}

type GetDashboardResultDto struct {
    // estado del servidor (claves de DashboardKeyEnum; solo primitivos)
    server map[string]any
    // un map por proyecto con su semaforo y contadores (solo primitivos)
    projects []map[string]any
}

func GetDashboardResultDtoFromPrimitives(primitives map[string]any) *GetDashboardResultDto {
    server, _ := primitives["server"].(map[string]any)
    if server == nil {
        server = map[string]any{}
    }
    projects, _ := primitives["projects"].([]map[string]any)
    if projects == nil {
        projects = []map[string]any{}
    }
    return &GetDashboardResultDto{server: server, projects: projects}
}

func (d *GetDashboardResultDto) GetServer() map[string]any {
    return d.server
}

func (d *GetDashboardResultDto) GetProjects() []map[string]any {
    return d.projects
}

func (d *GetDashboardResultDto) ToArray() GetDashboardPrimitives {
    return GetDashboardPrimitives{Server: d.server, Projects: d.projects}
}
