package services

import "siem-backend/src/modules/dashboard/domain/enums"

// ProjectHealthResolver el semaforo de un proyecto a partir de sus agentes y de sus eventos de las ultimas 24 h:
//   - empty: sin agentes habilitados.
//   - critical: algun evento critico, o ningun agente online.
//   - warning: algun agente caido, o algun error.
//   - ok: el resto.
type ProjectHealthResolver struct{}

func ProjectHealthResolverGetInstance() *ProjectHealthResolver {
    return &ProjectHealthResolver{}
}

func (r *ProjectHealthResolver) GetProjectHealth(enabledAgentsTotal int, agentsOnline int, errors24h int, criticals24h int) string {
    switch {
    case enabledAgentsTotal == 0:
        return enums.ProjectHealthEnum.EMPTY
    case criticals24h > 0 || agentsOnline == 0:
        return enums.ProjectHealthEnum.CRITICAL
    case agentsOnline < enabledAgentsTotal || errors24h > 0:
        return enums.ProjectHealthEnum.WARNING
    }
    return enums.ProjectHealthEnum.OK
}
