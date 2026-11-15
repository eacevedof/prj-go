package services

import "regexp"

// siem-client-<os>-<version>[.exe], como los nombra "make dist" del cliente
var downloadFileNamePattern = regexp.MustCompile(`^siem-client-(linux|windows|macos)-(\d+\.\d+\.\d+)(\.exe)?$`)

// arquitectura para la que compila "make dist" cada SO
var archsByOs = map[string]string{
    "linux":   "amd64",
    "windows": "amd64",
    "macos":   "arm64",
}

// DownloadFileNameParser saca el SO, la arquitectura y la version del nombre de un binario del siem-client.
type DownloadFileNameParser struct{}

func DownloadFileNameParserGetInstance() *DownloadFileNameParser {
    return &DownloadFileNameParser{}
}

// GetOsArchAndVersion isBinary false si el fichero no es un binario del cliente (no se lista).
func (p *DownloadFileNameParser) GetOsArchAndVersion(fileName string) (os string, arch string, version string, isBinary bool) {
    matches := downloadFileNamePattern.FindStringSubmatch(fileName)
    if matches == nil {
        return "", "", "", false
    }
    return matches[1], archsByOs[matches[1]], matches[2], true
}
