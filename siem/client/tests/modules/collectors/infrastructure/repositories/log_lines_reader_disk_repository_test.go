package repositories

import (
    "os"
    "path/filepath"
    "testing"

    "siem-client/src/modules/collectors/infrastructure/repositories"
    "siem-client/tests"
)

func TestReadsOnlyCompleteLinesAndResumes(t *testing.T) {
    logFilePath := filepath.Join(t.TempDir(), "app.log")
    writeFile(t, logFilePath, "first\r\n\nsecond\npartial")
    logLinesReaderDiskRepository := repositories.LogLinesReaderDiskRepositoryGetInstance()

    lines, nextOffset, err := logLinesReaderDiskRepository.GetLinesFromOffset(logFilePath, 0, 1<<20)
    tests.AssertNil(t, err)
    tests.AssertEquals(t, []string{"first", "second"}, lines)
    tests.AssertEquals(t, int64(len("first\r\n\nsecond\n")), nextOffset)

    appendFile(t, logFilePath, " done\nthird\n")
    lines, _, err = logLinesReaderDiskRepository.GetLinesFromOffset(logFilePath, nextOffset, 1<<20)
    tests.AssertNil(t, err)
    tests.AssertEquals(t, []string{"partial done", "third"}, lines)
}

func TestTruncatedFileIsReadFromTheStart(t *testing.T) {
    logFilePath := filepath.Join(t.TempDir(), "app.log")
    writeFile(t, logFilePath, "new\n")

    lines, nextOffset, err := repositories.LogLinesReaderDiskRepositoryGetInstance().GetLinesFromOffset(logFilePath, 500, 1<<20)
    tests.AssertNil(t, err)
    tests.AssertEquals(t, []string{"new"}, lines)
    tests.AssertEquals(t, int64(4), nextOffset)
}

func TestLineLongerThanMaxBytesAdvances(t *testing.T) {
    logFilePath := filepath.Join(t.TempDir(), "app.log")
    writeFile(t, logFilePath, "0123456789")

    lines, nextOffset, err := repositories.LogLinesReaderDiskRepositoryGetInstance().GetLinesFromOffset(logFilePath, 0, 4)
    tests.AssertNil(t, err)
    tests.AssertEquals(t, []string{"0123"}, lines)
    tests.AssertEquals(t, int64(4), nextOffset)
}

func writeFile(t *testing.T, filePath string, content string) {
    t.Helper()
    if err := os.WriteFile(filePath, []byte(content), 0o644); err != nil {
        t.Fatal(err)
    }
}

func appendFile(t *testing.T, filePath string, content string) {
    t.Helper()
    file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0o644)
    if err != nil {
        t.Fatal(err)
    }
    defer file.Close()
    if _, err := file.WriteString(content); err != nil {
        t.Fatal(err)
    }
}
