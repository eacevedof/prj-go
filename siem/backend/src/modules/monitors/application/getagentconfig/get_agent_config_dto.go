package getagentconfig

import "siem-backend/src/modules/shared/application"

// GetAgentConfigDto GET /api/v1/agent-config: el agente autenticado por su token (authAgentId + authProjectId).
type GetAgentConfigDto struct {
    application.AbstractInputDto
    authAgentId   int
    authProjectId int
}

func GetAgentConfigDtoFromPrimitives(primitives map[string]any) *GetAgentConfigDto {
    return &GetAgentConfigDto{
        AbstractInputDto: application.NewAbstractInputDto(primitives),
        authAgentId:      application.GetInt(primitives, "authAgentId"),
        authProjectId:    application.GetInt(primitives, "authProjectId"),
    }
}

func (d *GetAgentConfigDto) GetAuthAgentId() int {
    return d.authAgentId
}

func (d *GetAgentConfigDto) GetAuthProjectId() int {
    return d.authProjectId
}
