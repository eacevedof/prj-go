package repositories

import (
    "bytes"
    "io"
    "os"
    "strings"
)

// LogLinesReaderDiskRepository lee las lineas nuevas de un fichero desde un offset (tail -f sin quedarse abierto).
type LogLinesReaderDiskRepository struct{}

func LogLinesReaderDiskRepositoryGetInstance() *LogLinesReaderDiskRepository {
    return &LogLinesReaderDiskRepository{}
}

// GetLinesFromOffset lineas completas desde offset (como mucho maxBytes por llamada) y el offset siguiente.
//   - una linea a medio escribir (sin \n al final) se deja para la proxima pasada.
//   - si el fichero ahora es mas corto que offset, se trunco o roto: se vuelve a leer desde 0.
//   - una sola linea mas larga que maxBytes se entrega troceada (si no, el offset no avanzaria nunca).
//
// Las lineas vacias no se devuelven.
func (r *LogLinesReaderDiskRepository) GetLinesFromOffset(filePath string, offset int64, maxBytes int) ([]string, int64, error) {
    logFile, err := os.Open(filePath)
    if err != nil {
        return nil, offset, err
    }
    defer logFile.Close()

    fileInfo, err := logFile.Stat()
    if err != nil {
        return nil, offset, err
    }
    fileSize := fileInfo.Size()
    if fileSize < offset {
        offset = 0
    }
    if fileSize == offset {
        return []string{}, offset, nil
    }

    if _, err := logFile.Seek(offset, io.SeekStart); err != nil {
        return nil, offset, err
    }
    buffer := make([]byte, min(fileSize-offset, int64(maxBytes)))
    readBytes, err := io.ReadFull(logFile, buffer)
    if err != nil && err != io.ErrUnexpectedEOF {
        return nil, offset, err
    }
    buffer = buffer[:readBytes]

    lastNewLineIndex := bytes.LastIndexByte(buffer, '\n')
    if lastNewLineIndex < 0 {
        if readBytes < maxBytes {
            return []string{}, offset, nil
        }
        lastNewLineIndex = readBytes - 1
    }

    lines := []string{}
    for _, line := range strings.Split(string(buffer[:lastNewLineIndex+1]), "\n") {
        line = strings.TrimRight(line, "\r")
        if strings.TrimSpace(line) != "" {
            lines = append(lines, line)
        }
    }
    return lines, offset + int64(lastNewLineIndex+1), nil
}
