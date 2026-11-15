package services

import (
    "testing"

    "siem-backend/src/modules/downloads/domain/services"
    "siem-backend/tests"
)

func TestBinaryNamesAreParsed(t *testing.T) {
    downloadFileNameParser := services.DownloadFileNameParserGetInstance()
    for fileName, expected := range map[string][3]string{
        "siem-client-linux-0.3.0":       {"linux", "amd64", "0.3.0"},
        "siem-client-windows-0.3.0.exe": {"windows", "amd64", "0.3.0"},
        "siem-client-macos-1.12.4":      {"macos", "arm64", "1.12.4"},
    } {
        os, arch, version, isBinary := downloadFileNameParser.GetOsArchAndVersion(fileName)
        tests.AssertEquals(t, true, isBinary, fileName)
        tests.AssertEquals(t, expected, [3]string{os, arch, version}, fileName)
    }
}

func TestOtherFilesAreNotListed(t *testing.T) {
    downloadFileNameParser := services.DownloadFileNameParserGetInstance()
    for _, fileName := range []string{"readme.txt", "siem-client-freebsd-0.3.0", "siem-client-linux-latest", "siem-client-linux-0.3.0.sig", "../siem-client-linux-0.3.0"} {
        _, _, _, isBinary := downloadFileNameParser.GetOsArchAndVersion(fileName)
        tests.AssertEquals(t, false, isBinary, fileName)
    }
}
