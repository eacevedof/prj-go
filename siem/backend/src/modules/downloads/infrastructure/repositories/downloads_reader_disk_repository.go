package repositories

import (
    "crypto/sha256"
    "encoding/hex"
    "io"
    "os"
    "path/filepath"
    "sort"

    "siem-backend/src/modules/downloads/domain/enums"
)

// DownloadsReaderDiskRepository los binarios del siem-client que deja "make dist" del cliente en public/downloads/
// (el backend los sirve como estaticos en /downloads/<fichero>).
type DownloadsReaderDiskRepository struct{}

func DownloadsReaderDiskRepositoryGetInstance() *DownloadsReaderDiskRepository {
    return &DownloadsReaderDiskRepository{}
}

// GetDownloadFiles los ficheros de la carpeta (no las subcarpetas ni los ocultos), por nombre. Sin carpeta: vacio.
// Por fichero: file_name, size_bytes y sha256 (se calcula aqui: son tres binarios de ~7 MB).
func (r *DownloadsReaderDiskRepository) GetDownloadFiles(pathDownloads string) ([]map[string]any, error) {
    dirEntries, err := os.ReadDir(pathDownloads)
    if os.IsNotExist(err) {
        return []map[string]any{}, nil
    }
    if err != nil {
        return nil, err
    }

    downloadFiles := []map[string]any{}
    for _, dirEntry := range dirEntries {
        if dirEntry.IsDir() || dirEntry.Name()[0] == '.' {
            continue
        }
        fileInfo, err := dirEntry.Info()
        if err != nil {
            continue
        }
        fileSha256, err := r.getSha256(filepath.Join(pathDownloads, dirEntry.Name()))
        if err != nil {
            continue
        }
        downloadFiles = append(downloadFiles, map[string]any{
            enums.DownloadsKeyEnum.FILE_NAME:  dirEntry.Name(),
            enums.DownloadsKeyEnum.SIZE_BYTES: int(fileInfo.Size()),
            enums.DownloadsKeyEnum.SHA256:     fileSha256,
        })
    }
    sort.Slice(downloadFiles, func(i int, j int) bool {
        return downloadFiles[i][enums.DownloadsKeyEnum.FILE_NAME].(string) < downloadFiles[j][enums.DownloadsKeyEnum.FILE_NAME].(string)
    })
    return downloadFiles, nil
}

func (r *DownloadsReaderDiskRepository) getSha256(pathFile string) (string, error) {
    file, err := os.Open(pathFile)
    if err != nil {
        return "", err
    }
    defer file.Close()

    hash := sha256.New()
    if _, err := io.Copy(hash, file); err != nil {
        return "", err
    }
    return hex.EncodeToString(hash.Sum(nil)), nil
}
