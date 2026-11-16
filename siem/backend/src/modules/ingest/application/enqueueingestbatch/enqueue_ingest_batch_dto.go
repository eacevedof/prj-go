package enqueueingestbatch

import "siem-backend/src/modules/shared/application"

// EnqueueIngestBatchDto un lote del siem-client. eventsJson es el array "events" del body ya serializado (el dto solo
// lleva primitivos) y eventsCount cuantos trae (-1 si "events" no era un array). La ip del host es la de la peticion
// (AbstractInputDto.GetUserIpAddress).
type EnqueueIngestBatchDto struct {
    application.AbstractInputDto
    authAgentId   int
    batchUuid     string
    hostName      string
    hostOs        string
    clientVersion string
    clientPath    string
    eventsJson    string
    eventsCount   int
}

func EnqueueIngestBatchDtoFromPrimitives(primitives map[string]any) *EnqueueIngestBatchDto {
    return &EnqueueIngestBatchDto{
        AbstractInputDto: application.NewAbstractInputDto(primitives),
        authAgentId:      application.GetInt(primitives, "authAgentId"),
        batchUuid:        application.GetTrimmedString(primitives, "batchUuid"),
        hostName:         application.GetTrimmedString(primitives, "hostName"),
        hostOs:           application.GetTrimmedString(primitives, "hostOs"),
        clientVersion:    application.GetTrimmedString(primitives, "clientVersion"),
        clientPath:       application.GetTrimmedString(primitives, "clientPath"),
        eventsJson:       application.GetTrimmedString(primitives, "eventsJson"),
        eventsCount:      application.GetInt(primitives, "eventsCount"),
    }
}

func (d *EnqueueIngestBatchDto) GetAuthAgentId() int {
    return d.authAgentId
}

func (d *EnqueueIngestBatchDto) GetBatchUuid() string {
    return d.batchUuid
}

func (d *EnqueueIngestBatchDto) GetHostName() string {
    return d.hostName
}

func (d *EnqueueIngestBatchDto) GetHostOs() string {
    return d.hostOs
}

func (d *EnqueueIngestBatchDto) GetClientVersion() string {
    return d.clientVersion
}

func (d *EnqueueIngestBatchDto) GetEventsJson() string {
    return d.eventsJson
}

func (d *EnqueueIngestBatchDto) GetEventsCount() int {
    return d.eventsCount
}

func (d *EnqueueIngestBatchDto) GetClientPath() string {
    return d.clientPath
}
