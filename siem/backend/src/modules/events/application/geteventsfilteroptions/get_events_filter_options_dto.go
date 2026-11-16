package geteventsfilteroptions

import "siem-backend/src/modules/shared/application"

// GetEventsFilterOptionsDto GET /api/v1/events/filter-options: sin parametros, solo la sesion del panel.
type GetEventsFilterOptionsDto struct {
    application.AbstractInputDto
    isDashboardAuthenticated bool
}

func GetEventsFilterOptionsDtoFromPrimitives(primitives map[string]any) *GetEventsFilterOptionsDto {
    return &GetEventsFilterOptionsDto{
        AbstractInputDto:         application.NewAbstractInputDto(primitives),
        isDashboardAuthenticated: application.GetInt(primitives, "isDashboardAuthenticated") == 1,
    }
}

func (d *GetEventsFilterOptionsDto) IsDashboardAuthenticated() bool {
    return d.isDashboardAuthenticated
}
