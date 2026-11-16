package enums

// DownloadsConstantsEnum carpeta publica de los binarios del siem-client y ruta de ingesta que va en el .env de un
// agente (APP_SIEM_ENDPOINT_URL = url base + INGEST_PATH).
var DownloadsConstantsEnum = struct {
    DOWNLOADS_FOLDER string
    INGEST_PATH      string
    URL_SEPARATOR    string
}{
    DOWNLOADS_FOLDER: "downloads",
    INGEST_PATH:      "/api/v1/ingest-batches",
    URL_SEPARATOR:    "/",
}
