package getdownloads

import (
    "path/filepath"
    "strings"

    "siem-backend/src/core/boot"
    "siem-backend/src/modules/downloads/domain/services"
    "siem-backend/src/modules/downloads/infrastructure/repositories"
    "siem-backend/src/modules/shared/application"
    "siem-backend/src/modules/shared/infrastructure/repositories/configuration"
)

const (
    DOWNLOADS_FOLDER = "downloads"
    INGEST_PATH      = "/api/v1/ingest-batches"
)

// GetDownloadsService los binarios del siem-client que se pueden descargar desde el panel (los deja ahi "make dist"
// del cliente) y la url de ingesta que tiene que llevar el .env de un agente nuevo. Los ficheros de la carpeta que no
// son binarios del cliente no se listan.
type GetDownloadsService struct {
    environmentReaderRawRepository *configuration.EnvironmentReaderRawRepository
    downloadFileNameParser         *services.DownloadFileNameParser
    downloadsReaderDiskRepository  *repositories.DownloadsReaderDiskRepository
}

func GetDownloadsServiceGetInstance() *GetDownloadsService {
    return &GetDownloadsService{
        environmentReaderRawRepository: configuration.GetInstance(),
        downloadFileNameParser:         services.DownloadFileNameParserGetInstance(),
        downloadsReaderDiskRepository:  repositories.DownloadsReaderDiskRepositoryGetInstance(),
    }
}

// Invoke error: solo uno inesperado (disco). El 401 lo da antes el controller.
func (s *GetDownloadsService) Invoke(getDownloadsDto *GetDownloadsDto) (*GetDownloadsResultDto, error) {
    downloadFiles, err := s.downloadsReaderDiskRepository.GetDownloadFiles(filepath.Join(boot.BootPathsEnum.PATH_PUBLIC, DOWNLOADS_FOLDER))
    if err != nil {
        return nil, err
    }

    downloads := []DownloadPrimitives{}
    for _, downloadFile := range downloadFiles {
        fileName := application.GetTrimmedString(downloadFile, "file_name")
        os, arch, version, isBinary := s.downloadFileNameParser.GetOsArchAndVersion(fileName)
        if !isBinary {
            continue
        }
        downloads = append(downloads, DownloadPrimitives{
            FileName:  fileName,
            Os:        os,
            Arch:      arch,
            Version:   version,
            SizeBytes: application.GetInt(downloadFile, "size_bytes"),
            Sha256:    application.GetTrimmedString(downloadFile, "sha256"),
            Url:       "/" + DOWNLOADS_FOLDER + "/" + fileName,
        })
    }
    return GetDownloadsResultDtoFromPrimitives(map[string]any{
        "endpointUrl": strings.TrimRight(s.environmentReaderRawRepository.GetBaseUrl(), "/") + INGEST_PATH,
        "downloads":   downloads,
    }), nil
}
