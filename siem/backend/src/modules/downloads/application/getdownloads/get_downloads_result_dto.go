package getdownloads

import "siem-backend/src/modules/shared/application"

// DownloadPrimitives un binario del siem-client. url: publica (sin token), relativa al origen del panel.
type DownloadPrimitives struct {
    FileName  string `json:"file_name"`
    Os        string `json:"os"`
    Arch      string `json:"arch"`
    Version   string `json:"version"`
    SizeBytes int    `json:"size_bytes"`
    Sha256    string `json:"sha256"`
    Url       string `json:"url"`
}

// InstallerPrimitives un instalador asistido (install.sh / install.ps1 / install.bat). os: "linux,macos" o "windows".
type InstallerPrimitives struct {
    FileName  string `json:"file_name"`
    Os        string `json:"os"`
    SizeBytes int    `json:"size_bytes"`
    Sha256    string `json:"sha256"`
    Url       string `json:"url"`
}

// GetDownloadsPrimitives lo que devuelve la API. endpoint_url: lo que va en APP_SIEM_ENDPOINT_URL del agente.
type GetDownloadsPrimitives struct {
    EndpointUrl string               `json:"endpoint_url"`
    Downloads   []DownloadPrimitives  `json:"downloads"`
    Installers  []InstallerPrimitives `json:"installers"`
}

type GetDownloadsResultDto struct {
    endpointUrl string
    downloads   []DownloadPrimitives
    installers  []InstallerPrimitives
}

func GetDownloadsResultDtoFromPrimitives(primitives map[string]any) *GetDownloadsResultDto {
    downloads, _ := primitives["downloads"].([]DownloadPrimitives)
    installers, _ := primitives["installers"].([]InstallerPrimitives)
    if installers == nil {
        installers = []InstallerPrimitives{}
    }
    if downloads == nil {
        downloads = []DownloadPrimitives{}
    }
    return &GetDownloadsResultDto{
        endpointUrl: application.GetTrimmedString(primitives, "endpointUrl"),
        downloads:   downloads,
        installers:  installers,
    }
}

func (d *GetDownloadsResultDto) ToArray() GetDownloadsPrimitives {
    return GetDownloadsPrimitives{EndpointUrl: d.endpointUrl, Downloads: d.downloads, Installers: d.installers}
}
