package enums

// DownloadsKeyEnum claves de cada fichero que lee DownloadsReaderDiskRepository.
var DownloadsKeyEnum = struct {
    FILE_NAME  string
    SIZE_BYTES string
    SHA256     string
}{
    FILE_NAME:  "file_name",
    SIZE_BYTES: "size_bytes",
    SHA256:     "sha256",
}
