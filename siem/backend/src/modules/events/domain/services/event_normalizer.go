package services

import (
    "encoding/json"
    "slices"
    "time"

    "siem-backend/src/modules/events/domain/entities"
    "siem-backend/src/modules/events/domain/enums"
    "siem-backend/src/modules/shared/application"
)

const (
    dateTimeFormat         = "2006-01-02 15:04:05"
    MESSAGE_MAX_LENGTH     = 8000
    SOURCE_NAME_MAX_LENGTH = 255
    FILE_PATH_MAX_LENGTH   = 500
    // attributes son unos pocos datos estructurados (estado, exit_code, % de cpu...): mas que esto no es un attribute
    ATTRIBUTES_MAX_BYTES = 4000
    EMPTY_COMMAND_OUTPUT = "(no output)"
)

// EventNormalizer de un evento crudo del siem-client (un objeto JSON del lote) a un EventEntity valido: recorta
// longitudes, deduce tipo de source y gravedad, y corrige la fecha. Lo que no se puede salvar se descarta (nil): un
// evento malo no tumba el lote.
type EventNormalizer struct {
    eventSeverityResolver *EventSeverityResolver
}

func EventNormalizerGetInstance() *EventNormalizer {
    return &EventNormalizer{eventSeverityResolver: EventSeverityResolverGetInstance()}
}

// GetEventEntityOrNil nil si el evento no tiene un event_type conocido, source_name o mensaje (salvo un comando sin
// salida, que se anota con EMPTY_COMMAND_OUTPUT).
func (n *EventNormalizer) GetEventEntityOrNil(rawEvent map[string]any, projectId int, agentId int, receivedAt string) *entities.EventEntity {
    eventType := application.GetTrimmedString(rawEvent, enums.EventsKeyEnum.EVENT_TYPE)
    if !slices.Contains(enums.EventTypeEnum.VALUES, eventType) {
        return nil
    }

    sourceName := application.GetTrimmedString(rawEvent, enums.EventsKeyEnum.SOURCE_NAME)
    if sourceName == "" || len([]rune(sourceName)) > SOURCE_NAME_MAX_LENGTH {
        return nil
    }

    message := n.getTruncated(application.GetTrimmedString(rawEvent, enums.EventsKeyEnum.MESSAGE), MESSAGE_MAX_LENGTH)
    if message == "" && eventType == enums.EventTypeEnum.COMMAND_RESULT {
        message = EMPTY_COMMAND_OUTPUT
    }
    if message == "" {
        return nil
    }

    attributes := n.getAttributesOrNil(rawEvent[enums.EventsKeyEnum.ATTRIBUTES])
    exitCode := application.GetInt(rawEvent, enums.EventsKeyEnum.EXIT_CODE)
    return entities.EventEntityFromPrimitives(map[string]any{
        enums.EventsKeyEnum.PROJECT_ID:  projectId,
        enums.EventsKeyEnum.AGENT_ID:    agentId,
        enums.EventsKeyEnum.EVENT_TYPE:  eventType,
        enums.EventsKeyEnum.SEVERITY:    n.eventSeverityResolver.GetSeverity(eventType, message, exitCode, application.GetTrimmedString(rawEvent, enums.EventsKeyEnum.SEVERITY), attributes),
        enums.EventsKeyEnum.MESSAGE:     message,
        enums.EventsKeyEnum.SOURCE_TYPE: enums.SourceTypeByEventType[eventType],
        enums.EventsKeyEnum.SOURCE_NAME: sourceName,
        enums.EventsKeyEnum.FILE_PATH:   n.getTruncated(application.GetTrimmedString(rawEvent, enums.EventsKeyEnum.FILE_PATH), FILE_PATH_MAX_LENGTH),
        enums.EventsKeyEnum.EXIT_CODE:   exitCode,
        enums.EventsKeyEnum.DURATION_MS: max(0, application.GetInt(rawEvent, enums.EventsKeyEnum.DURATION_MS)),
        enums.EventsKeyEnum.ATTRIBUTES:  n.getAttributesJson(attributes),
        enums.EventsKeyEnum.OCCURRED_AT: n.getOccurredAtOrReceivedAt(application.GetTrimmedString(rawEvent, enums.EventsKeyEnum.OCCURRED_AT), receivedAt),
        enums.EventsKeyEnum.RECEIVED_AT: receivedAt,
    })
}

// getAttributesOrNil solo un objeto JSON que, compacto, no pase de ATTRIBUTES_MAX_BYTES; otra cosa se ignora (el
// evento se guarda igual, sin attributes).
func (n *EventNormalizer) getAttributesOrNil(rawAttributes any) map[string]any {
    attributes, isObject := rawAttributes.(map[string]any)
    if !isObject || len(attributes) == 0 {
        return nil
    }
    if len(n.getAttributesJson(attributes)) > ATTRIBUTES_MAX_BYTES {
        return nil
    }
    return attributes
}

func (n *EventNormalizer) getAttributesJson(attributes map[string]any) string {
    if attributes == nil {
        return ""
    }
    attributesJson, err := json.Marshal(attributes)
    if err != nil {
        return ""
    }
    return string(attributesJson)
}

// getOccurredAtOrReceivedAt la hora del cliente si es una fecha valida y no viene del futuro (reloj descolocado).
func (n *EventNormalizer) getOccurredAtOrReceivedAt(occurredAt string, receivedAt string) string {
    occurredTime, err := time.Parse(dateTimeFormat, occurredAt)
    if err != nil {
        return receivedAt
    }
    receivedTime, err := time.Parse(dateTimeFormat, receivedAt)
    if err == nil && occurredTime.After(receivedTime.Add(5*time.Minute)) {
        return receivedAt
    }
    return occurredAt
}

func (n *EventNormalizer) getTruncated(text string, maxLength int) string {
    runes := []rune(text)
    if len(runes) <= maxLength {
        return text
    }
    return string(runes[:maxLength])
}
