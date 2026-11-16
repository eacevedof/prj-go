package logger

import (
    "os"
    "path/filepath"
    "regexp"
    "strings"
    "sync"
    "time"

    "siem-backend/src/modules/shared/infrastructure/components/datetimer"
    "siem-backend/src/modules/shared/infrastructure/repositories/configuration"
)

// mismo fichero desde varias peticiones a la vez: las escrituras se serializan
var writeLogMutex sync.Mutex

var sqlMarginPattern = regexp.MustCompile(`(?m)^ {8}`)

// Logger escribe en <logsPath>/YYYY-MM-DD-<nombre>.<ext> (ext "sql" si el nombre contiene "sql").
type Logger struct{}

func GetInstance() *Logger {
    return &Logger{}
}

func (l *Logger) WriteLog(filePath string, content string, logsPath string) {
    if logsPath == "" {
        logsPath = configuration.GetInstance().GetLogPath()
    }
    if logsPath == "" {
        return
    }
    if err := os.MkdirAll(logsPath, 0o777); err != nil {
        return
    }

    logsFolderPath, err := filepath.Abs(logsPath)
    if err != nil {
        return
    }
    logFilePath := filepath.Join(logsFolderPath, filePath)
    finalDirectoryPath := filepath.Dir(logFilePath)
    _ = os.MkdirAll(finalDirectoryPath, 0o777)

    fileExtension := strings.TrimPrefix(filepath.Ext(logFilePath), ".")
    fileName := strings.TrimSuffix(filepath.Base(logFilePath), filepath.Ext(logFilePath))
    if fileExtension == "" {
        fileExtension = "log"
    }
    if strings.Contains(filePath, "sql") {
        fileExtension = "sql"
    }

    // la app trabaja en UTC (date_default_timezone_set("UTC") en task-php)
    now := time.Now().UTC()
    logFilePath = filepath.Join(finalDirectoryPath, now.Format(datetimer.FORMAT_ONLY_DATE)+"-"+fileName+"."+fileExtension)
    content = strings.TrimSpace(content)
    if fileExtension == "sql" {
        content = l.getNormalizedMarginForSql(content)
    }

    writeLogMutex.Lock()
    defer writeLogMutex.Unlock()
    _, statErr := os.Stat(logFilePath)
    isNewFile := os.IsNotExist(statErr)
    logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o666)
    if err != nil {
        return
    }
    defer logFile.Close()
    _, _ = logFile.WriteString("\n[" + now.Format(datetimer.FORMAT_FULL_DT) + "]\n" + content)
    if isNewFile {
        _ = os.Chmod(logFilePath, 0o666)
    }
}

func (l *Logger) getNormalizedMarginForSql(content string) string {
    return sqlMarginPattern.ReplaceAllString(content+"\n", "")
}
