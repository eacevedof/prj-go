package entities

import "strings"

// getIntOrDefault en siem-client.json un 0 (o la clave ausente) = el valor por defecto.
func getIntOrDefault(value int, defaultValue int) int {
    if value == 0 {
        return defaultValue
    }
    return value
}

func getStringOrDefault(value string, defaultValue string) string {
    if trimmedValue := strings.TrimSpace(value); trimmedValue != "" {
        return trimmedValue
    }
    return defaultValue
}

func trim(value string) string {
    return strings.TrimSpace(value)
}
