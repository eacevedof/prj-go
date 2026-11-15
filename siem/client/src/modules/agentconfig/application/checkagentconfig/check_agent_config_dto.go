package checkagentconfig

// CheckAgentConfigDto sin campos: la configuracion sale de APP_CONFIG_PATH y del .env. Existe para que el caso de uso
// tenga la misma forma que el resto (y crezca sin cambiar su firma).
type CheckAgentConfigDto struct{}

func CheckAgentConfigDtoFromPrimitives(primitives map[string]any) *CheckAgentConfigDto {
    return &CheckAgentConfigDto{}
}
