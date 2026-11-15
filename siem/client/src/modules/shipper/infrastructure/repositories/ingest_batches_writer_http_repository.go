package repositories

import (
    "siem-client/src/modules/shared/infrastructure/components/jsonposter"
    "siem-client/src/modules/shared/infrastructure/repositories/configuration"
)

// IngestBatchesWriterHttpRepository el backend como fuente de datos: escribe un lote en POST /api/v1/ingest-batches.
type IngestBatchesWriterHttpRepository struct {
    jsonPoster                     *jsonposter.JsonPoster
    environmentReaderRawRepository *configuration.EnvironmentReaderRawRepository
}

func IngestBatchesWriterHttpRepositoryGetInstance() *IngestBatchesWriterHttpRepository {
    return &IngestBatchesWriterHttpRepository{
        jsonPoster:                     jsonposter.GetInstance(),
        environmentReaderRawRepository: configuration.GetInstance(),
    }
}

// SendIngestBatch codigo HTTP y body de la respuesta. error = no hubo respuesta (red, DNS, timeout).
func (r *IngestBatchesWriterHttpRepository) SendIngestBatch(endpointUrl string, apiToken string, ingestBatch map[string]any) (int, string, error) {
    userAgent := r.environmentReaderRawRepository.GetAppName() + "/" + r.environmentReaderRawRepository.GetAppVersion()
    return r.jsonPoster.PostJson(endpointUrl, apiToken, userAgent, ingestBatch)
}
