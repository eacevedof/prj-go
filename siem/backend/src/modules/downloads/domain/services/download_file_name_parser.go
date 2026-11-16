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

// instaladores asistidos que copia "make dist" (client/install/) y para que sistemas sirve cada uno
var osByInstallerFileName = map[string]string{
    "install.sh":  "linux,macos",
    "install.ps1": "windows",
    "install.bat": "windows",
}

// DownloadFileNameParser saca el SO, la arquitectura y la version del nombre de un binario del siem-client, y que
// ficheros son instaladores.
type DownloadFileNameParser struct{}

func DownloadFileNameParserGetInstance() *DownloadFileNameParser {
    return &DownloadFileNameParser{}
}

// GetOsArchAndVersion isBinary false si el fichero no es un binario del cliente (no se lista como binario).
func (p *DownloadFileNameParser) GetOsArchAndVersion(fileName string) (os string, arch string, version string, isBinary bool) {
    matches := downloadFileNamePattern.FindStringSubmatch(fileName)
    if matches == nil {
        return "", "", "", false
    }
    return matches[1], archsByOs[matches[1]], matches[2], true
}

// GetInstallerOs para que sistemas es un instalador ("linux,macos" o "windows"); isInstaller false si no lo es.
func (p *DownloadFileNameParser) GetInstallerOs(fileName string) (os string, isInstaller bool) {
    os, isInstaller = osByInstallerFileName[fileName]
    return os, isInstaller
}
