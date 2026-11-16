package services

import (
    "regexp"

    "siem-backend/src/modules/agents/domain/enums"
)

// rutas absolutas sin nada que un shell pueda interpretar: ni comillas (tampoco las tipograficas), ni $, `, ;, &, |...
// Letras y digitos de cualquier idioma (C:\Users\José) y los separadores habituales en nombres de carpeta.
var (
    windowsClientPathPattern = regexp.MustCompile(`^[A-Za-z]:\\[\p{L}\p{N} _.()+\-\\]*$`)
    unixClientPathPattern    = regexp.MustCompile(`^/[\p{L}\p{N} _.+@\-/]*$`)
)

// AgentHostDataSanitizer los datos de host que el siem-client informa de si mismo al ingestar. No son de fiar (los
// manda quien tenga el token del agente) y el panel los usa para montar el comando de desinstalar que el admin pega
// en un shell elevado: lo que no encaja se guarda vacio y el panel cae en la ruta por defecto del instalador.
type AgentHostDataSanitizer struct{}

func AgentHostDataSanitizerGetInstance() *AgentHostDataSanitizer {
    return &AgentHostDataSanitizer{}
}

// GetHostOs el host_os si es uno de AgentHostOsEnum; "" si no.
func (s *AgentHostDataSanitizer) GetHostOs(hostOs string) string {
    switch hostOs {
    case enums.AgentHostOsEnum.LINUX, enums.AgentHostOsEnum.WINDOWS, enums.AgentHostOsEnum.MACOS:
        return hostOs
    }
    return ""
}

// GetClientPath la ruta del ejecutable si es una ruta absoluta limpia para el SO del host; "" si no.
func (s *AgentHostDataSanitizer) GetClientPath(hostOs string, clientPath string) string {
    clientPathPattern := unixClientPathPattern
    if s.GetHostOs(hostOs) == enums.AgentHostOsEnum.WINDOWS {
        clientPathPattern = windowsClientPathPattern
    }
    if !clientPathPattern.MatchString(clientPath) {
        return ""
    }
    return clientPath
}
