package controllers

import (
    "errors"

    agentsexceptions "siem-backend/src/modules/agents/domain/exceptions"
    "siem-backend/src/modules/ingest/application/enqueueingestbatch"
    "siem-backend/src/modules/ingest/domain/enums"
    "siem-backend/src/modules/ingest/domain/exceptions"
    "siem-backend/src/modules/shared/application"
    "siem-backend/src/modules/shared/infrastructure/components/httprequester"
    "siem-backend/src/modules/shared/infrastructure/components/httpresponder"
    sharedcontrollers "siem-backend/src/modules/shared/infrastructure/controllers"
    sharedexceptions "siem-backend/src/modules/shared/infrastructure/exceptions"
)

// EnqueueIngestBatchRestController POST /api/v1/ingest-batches  (Bearer: api_token del agente)
//
//    body: {batch_uuid, agent: {host_name, host_os, client_version, client_path}, events: [{event_type, source_name, message, ...}]}
//    202: {batch_uuid, queued_events, is_duplicate}. "events" vacio = solo latido.
type EnqueueIngestBatchRestController struct {
    sharedcontrollers.AbstractRestController
}

func EnqueueIngestBatchRestControllerGetInstance(httpRequester *httprequester.HttpRequester) *EnqueueIngestBatchRestController {
    return &EnqueueIngestBatchRestController{AbstractRestController: sharedcontrollers.NewAbstractRestController(httpRequester)}
}

func (c *EnqueueIngestBatchRestController) Invoke() *httpresponder.HttpResponder {
    enqueueIngestBatchResultDto, err := c.getEnqueueIngestBatchResultDto()
    if err == nil {
        return c.RespondAccepted(enqueueIngestBatchResultDto.ToArray())
    }

    var agentsException *agentsexceptions.AgentsException
    var ingestException *exceptions.IngestException
    var fieldsException *sharedexceptions.FieldsException
    switch {
    case errors.As(err, &agentsException):
        return c.RespondDomainError(agentsException)
    case errors.As(err, &ingestException):
        return c.RespondDomainError(ingestException)
    case errors.As(err, &fieldsException):
        return c.RespondDomainError(fieldsException)
    }
    return c.RespondUnexpectedError(err, "EnqueueIngestBatchRestController.Invoke")
}

func (c *EnqueueIngestBatchRestController) getEnqueueIngestBatchResultDto() (*enqueueingestbatch.EnqueueIngestBatchResultDto, error) {
    authAgentPrimitives, err := c.GetAuthAgentOrFail()
    if err != nil {
        return nil, err
    }

    httpRequest := c.GetHttpRequest()
    agent, _ := httpRequest.GetPostValue(enums.IngestKeyEnum.AGENT, map[string]any{}).(map[string]any)
    eventsJson, eventsCount := c.getEventsJsonAndCount()
    return enqueueingestbatch.EnqueueIngestBatchServiceGetInstance().Invoke(
        enqueueingestbatch.EnqueueIngestBatchDtoFromPrimitives(application.MergePrimitives(
            httpRequest.GetMetadata(),
            authAgentPrimitives,
            map[string]any{
                "batchUuid":     httpRequest.GetPostString(enums.IngestKeyEnum.BATCH_UUID),
                "hostName":      application.GetTrimmedString(agent, enums.IngestKeyEnum.HOST_NAME),
                "hostOs":        application.GetTrimmedString(agent, enums.IngestKeyEnum.HOST_OS),
                "clientVersion": application.GetTrimmedString(agent, enums.IngestKeyEnum.CLIENT_VERSION),
                "clientPath":    application.GetTrimmedString(agent, enums.IngestKeyEnum.CLIENT_PATH),
                "eventsJson":    eventsJson,
                "eventsCount":   eventsCount,
            },
        )),
    )
}

// getEventsJsonAndCount el array "events" del body tal cual, serializado (un dto solo lleva primitivos), y cuantos
// trae. Sin "events" = 0 eventos (latido); si no es un array, -1 (lo rechaza el servicio con un 400).
func (c *EnqueueIngestBatchRestController) getEventsJsonAndCount() (string, int) {
    rawEvents := c.GetHttpRequest().GetPostValue(enums.IngestKeyEnum.EVENTS, []any{})
    events, isArray := rawEvents.([]any)
    if !isArray {
        return "[]", -1
    }
    return sharedcontrollers.GetJsonEncoded(events), len(events)
}
