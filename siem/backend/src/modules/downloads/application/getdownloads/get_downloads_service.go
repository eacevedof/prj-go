package getdownloads

import (
    "path/filepath"
    "strings"

    "siem-backend/src/core/boot"

    "siem-backend/src/modules/shared/application"
    "siem-backend/src/modules/shared/infrastructure/repositories/configuration"

    "siem-backend/src/modules/downloads/domain/enums"
    "siem-backend/src/modules/downloads/domain/services"
    "siem-backend/src/modules/downloads/infrastructure/repositories"
)

// GetDownloadsService los binarios del siem-client y sus instaladores asistidos que se pueden descargar desde el panel
// (los deja ahi "make dist" del cliente) y la url de ingesta que tiene que llevar el .env de un agente nuevo. Los
// ficheros de la carpeta que no son binarios del cliente no se listan.
type GetDownloadsService struct {
    environmentReaderRawRepository *configuration.EnvironmentReaderRawRepository
    downloadsReaderDiskRepository  *repositories.DownloadsReaderDiskRepository
    downloadFileNameParser         *services.DownloadFileNameParser
}

func GetDownloadsServiceGetInstance() *GetDownloadsService {
    return &GetDownloadsService{
        environmentReaderRawRepository: configuration.GetInstance(),
        downloadsReaderDiskRepository:  repositories.DownloadsReaderDiskRepositoryGetInstance(),
        downloadFileNameParser:         services.DownloadFileNameParserGetInstance(),
    }
}

// Invoke error: solo uno inesperado (disco). El 401 lo da antes el controller.
func (s *GetDownloadsService) Invoke(getDownloadsDto *GetDownloadsDto) (*GetDownloadsResultDto, error) {
    downloadFiles, err := s.downloadsReaderDiskRepository.GetDownloadFiles(filepath.Join(boot.BootPathsEnum.PATH_PUBLIC, enums.DownloadsConstantsEnum.DOWNLOADS_FOLDER))
    if err != nil {
        return nil, err
    }

    downloads := []DownloadPrimitives{}
    installers := []InstallerPrimitives{}
    for _, downloadFile := range downloadFiles {
        fileName := application.GetTrimmedString(downloadFile, enums.DownloadsKeyEnum.FILE_NAME)
        if installerOs, isInstaller := s.downloadFileNameParser.GetInstallerOs(fileName); isInstaller {
            installers = append(installers, InstallerPrimitives{
                FileName:  fileName,
                Os:        installerOs,
                SizeBytes: application.GetInt(downloadFile, enums.DownloadsKeyEnum.SIZE_BYTES),
                Sha256:    application.GetTrimmedString(downloadFile, enums.DownloadsKeyEnum.SHA256),
                Url:       s.getDownloadUrl(fileName),
            })
            continue
        }
        os, arch, version, isBinary := s.downloadFileNameParser.GetOsArchAndVersion(fileName)
        if !isBinary {
            continue
        }
        downloads = append(downloads, DownloadPrimitives{
            FileName:  fileName,
            Os:        os,
            Arch:      arch,
            Version:   version,
            SizeBytes: application.GetInt(downloadFile, enums.DownloadsKeyEnum.SIZE_BYTES),
            Sha256:    application.GetTrimmedString(downloadFile, enums.DownloadsKeyEnum.SHA256),
            Url:       s.getDownloadUrl(fileName),
        })
    }
    return GetDownloadsResultDtoFromPrimitives(map[string]any{
        "endpointUrl": strings.TrimRight(s.environmentReaderRawRepository.GetBaseUrl(), enums.DownloadsConstantsEnum.URL_SEPARATOR) + enums.DownloadsConstantsEnum.INGEST_PATH,
        "downloads":   downloads,
        "installers":  installers,
    }), nil
}

// getDownloadUrl publica (sin token), relativa al origen del panel: /downloads/<fichero>.
func (s *GetDownloadsService) getDownloadUrl(fileName string) string {
    urlSeparator := enums.DownloadsConstantsEnum.URL_SEPARATOR
    return urlSeparator + enums.DownloadsConstantsEnum.DOWNLOADS_FOLDER + urlSeparator + fileName
}
