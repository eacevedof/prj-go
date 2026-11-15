package repositories

import (
    "io/fs"
    "os"
    "path/filepath"
    "sort"
)

// WatchedFilesReaderDiskRepository que ficheros hay ahora mismo en una carpeta vigilada.
type WatchedFilesReaderDiskRepository struct{}

func WatchedFilesReaderDiskRepositoryGetInstance() *WatchedFilesReaderDiskRepository {
    return &WatchedFilesReaderDiskRepository{}
}

func (r *WatchedFilesReaderDiskRepository) DoesFolderExist(folderPath string) bool {
    fileInfo, err := os.Stat(folderPath)
    return err == nil && fileInfo.IsDir()
}

// GetFilePathsByPattern rutas absolutas de los ficheros cuyo NOMBRE casa con pattern (glob: *.log, app-*.txt),
// ordenadas. Recursivo: tambien en subcarpetas (una subcarpeta sin permiso de lectura se salta).
func (r *WatchedFilesReaderDiskRepository) GetFilePathsByPattern(folderPath string, pattern string, isRecursive bool) ([]string, error) {
    filePaths := []string{}
    if !isRecursive {
        dirEntries, err := os.ReadDir(folderPath)
        if err != nil {
            return nil, err
        }
        for _, dirEntry := range dirEntries {
            if dirEntry.Type().IsRegular() && r.isMatch(pattern, dirEntry.Name()) {
                filePaths = append(filePaths, filepath.Join(folderPath, dirEntry.Name()))
            }
        }
        return filePaths, nil
    }

    err := filepath.WalkDir(folderPath, func(path string, dirEntry fs.DirEntry, walkError error) error {
        if walkError != nil {
            if path == folderPath {
                return walkError
            }
            return fs.SkipDir
        }
        if dirEntry.Type().IsRegular() && r.isMatch(pattern, dirEntry.Name()) {
            filePaths = append(filePaths, path)
        }
        return nil
    })
    if err != nil {
        return nil, err
    }
    sort.Strings(filePaths)
    return filePaths, nil
}

// GetFileSize el tamaño actual (start_at=end: se empieza a leer desde aqui).
func (r *WatchedFilesReaderDiskRepository) GetFileSize(filePath string) (int64, error) {
    fileInfo, err := os.Stat(filePath)
    if err != nil {
        return 0, err
    }
    return fileInfo.Size(), nil
}

func (r *WatchedFilesReaderDiskRepository) isMatch(pattern string, fileName string) bool {
    isMatch, err := filepath.Match(pattern, fileName)
    return err == nil && isMatch
}
