package server

import (
    "net"
    "os"
    "path/filepath"
)

// ServerIdentifier nombre e ip de la maquina que ejecuta la app.
type ServerIdentifier struct{}

func GetInstance() *ServerIdentifier {
    return &ServerIdentifier{}
}

func (s *ServerIdentifier) GetServerName() string {
    hostname, err := os.Hostname()
    if err != nil {
        return ""
    }
    return hostname
}

// GetServerIp la primera IPv4 del hostname, o el propio nombre si no se resuelve.
func (s *ServerIdentifier) GetServerIp() string {
    serverName := s.GetServerName()
    addresses, err := net.LookupHost(serverName)
    if err != nil {
        return serverName
    }
    for _, address := range addresses {
        if ip := net.ParseIP(address); ip != nil && ip.To4() != nil {
            return address
        }
    }
    return serverName
}

// GetExecutablePath ruta absoluta del binario que se esta ejecutando (symlinks resueltos), o "" si el SO no la da.
// Viaja al backend para que el panel diga donde desinstalar el agente.
func (s *ServerIdentifier) GetExecutablePath() string {
    executable, err := os.Executable()
    if err != nil {
        return ""
    }
    if resolved, err := filepath.EvalSymlinks(executable); err == nil {
        return resolved
    }
    return executable
}
