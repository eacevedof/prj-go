package entities

import (
    "siem-backend/src/modules/events/domain/enums"
    "siem-backend/src/modules/shared/application"
)

// EventEntity un evento ya normalizado (EventNormalizer), listo para insertar. sourceType + sourceName identifican
// su source dentro del agente; el id de esa source lo resuelve quien inserta.
type EventEntity struct {
    projectId   int
    agentId     int
    eventType   string
    severity    string
    message     string
    sourceType  string
    sourceName  string
    filePath    string
    exitCode    int
    durationMs  int
    // JSON compacto de los attributes del evento ("" = sin attributes)
    attributes  string
    occurredAt  string
    receivedAt  string
}

func EventEntityFromPrimitives(primitives map[string]any) *EventEntity {
    return &EventEntity{
        projectId:  application.GetInt(primitives, enums.EventsKeyEnum.PROJECT_ID),
        agentId:    application.GetInt(primitives, enums.EventsKeyEnum.AGENT_ID),
        eventType:  application.GetTrimmedString(primitives, enums.EventsKeyEnum.EVENT_TYPE),
        severity:   application.GetTrimmedString(primitives, enums.EventsKeyEnum.SEVERITY),
        message:    application.GetTrimmedString(primitives, enums.EventsKeyEnum.MESSAGE),
        sourceType: application.GetTrimmedString(primitives, enums.EventsKeyEnum.SOURCE_TYPE),
        sourceName: application.GetTrimmedString(primitives, enums.EventsKeyEnum.SOURCE_NAME),
        filePath:   application.GetTrimmedString(primitives, enums.EventsKeyEnum.FILE_PATH),
        exitCode:   application.GetInt(primitives, enums.EventsKeyEnum.EXIT_CODE),
        durationMs: application.GetInt(primitives, enums.EventsKeyEnum.DURATION_MS),
        attributes: application.GetTrimmedString(primitives, enums.EventsKeyEnum.ATTRIBUTES),
        occurredAt: application.GetTrimmedString(primitives, enums.EventsKeyEnum.OCCURRED_AT),
        receivedAt: application.GetTrimmedString(primitives, enums.EventsKeyEnum.RECEIVED_AT),
    }
}

func (e *EventEntity) GetProjectId() int {
    return e.projectId
}

func (e *EventEntity) GetAgentId() int {
    return e.agentId
}

func (e *EventEntity) GetEventType() string {
    return e.eventType
}

// HasExitCode el resultado de un comando y, de Docker, el estado de un contenedor o un evento del daemon (die).
func (e *EventEntity) HasExitCode() bool {
    return e.eventType == enums.EventTypeEnum.COMMAND_RESULT || e.eventType == enums.EventTypeEnum.DOCKER_STATE ||
        e.eventType == enums.EventTypeEnum.DOCKER_EVENT
}

func (e *EventEntity) IsCommandResult() bool {
    return e.eventType == enums.EventTypeEnum.COMMAND_RESULT
}

func (e *EventEntity) GetSeverity() string {
    return e.severity
}

func (e *EventEntity) GetMessage() string {
    return e.message
}

func (e *EventEntity) GetSourceType() string {
    return e.sourceType
}

func (e *EventEntity) GetSourceName() string {
    return e.sourceName
}

func (e *EventEntity) GetFilePath() string {
    return e.filePath
}

func (e *EventEntity) GetExitCode() int {
    return e.exitCode
}

func (e *EventEntity) GetDurationMs() int {
    return e.durationMs
}

func (e *EventEntity) GetOccurredAt() string {
    return e.occurredAt
}

func (e *EventEntity) GetReceivedAt() string {
    return e.receivedAt
}

func (e *EventEntity) GetAttributes() string {
    return e.attributes
}
