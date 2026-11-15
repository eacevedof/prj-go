package enums

// IngestMessageEnum mensajes de error de ingest ({0}... los rellena config.Tr).
var IngestMessageEnum = struct {
    INVALID_OR_MISSING_API_TOKEN     string
    BATCH_UUID_IS_NOT_VALID          string
    EVENTS_MUST_BE_AN_ARRAY          string
    EVENTS_MAX_PER_BATCH_IS_2000     string
    MAX_BATCHES_MUST_BE_1_TO_1000    string
    BATCH_PAYLOAD_IS_NOT_A_JSON_LIST string
}{
    INVALID_OR_MISSING_API_TOKEN:     "Invalid or missing api token",
    BATCH_UUID_IS_NOT_VALID:          "batch_uuid: required, 8-60 chars of letters, digits and -",
    EVENTS_MUST_BE_AN_ARRAY:          "events: must be an array",
    EVENTS_MAX_PER_BATCH_IS_2000:     "events: max 2000 per batch",
    MAX_BATCHES_MUST_BE_1_TO_1000:    "max_batches: must be between 1 and 1000",
    BATCH_PAYLOAD_IS_NOT_A_JSON_LIST: "batch payload is not a JSON list of events",
}
