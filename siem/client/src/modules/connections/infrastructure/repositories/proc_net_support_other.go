//go:build !linux

package repositories

// isProcNetSupported fuera de Linux no hay /proc/net: el modulo compila pero no muestrea, y el agente lo avisa una vez
// al backend (evento is_unavailable). Windows (GetExtendedTcpTable) y macOS quedan para mas adelante.
const isProcNetSupported = false
