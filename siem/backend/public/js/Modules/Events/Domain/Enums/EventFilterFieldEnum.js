/**
 * Campos por los que la API deja filtrar eventos (incluir y excluir, varios valores con coma).
 */
export const EventFilterFieldEnum = Object.freeze({
    PROJECT_UUID: "project_uuid",
    AGENT_UUID: "agent_uuid",
    SEVERITY: "severity",
    EVENT_TYPE: "event_type",
    SOURCE_TYPE: "source_type",
    SOURCE_NAME: "source_name",
});

export const EVENT_FILTER_FIELDS = Object.freeze(Object.values(EventFilterFieldEnum));

// como se llama cada campo en las pastillas de filtro
export const EVENT_FILTER_FIELD_LABELS = Object.freeze({
    [EventFilterFieldEnum.PROJECT_UUID]: "proyecto",
    [EventFilterFieldEnum.AGENT_UUID]: "agente",
    [EventFilterFieldEnum.SEVERITY]: "gravedad",
    [EventFilterFieldEnum.EVENT_TYPE]: "tipo",
    [EventFilterFieldEnum.SOURCE_TYPE]: "tipo de fuente",
    [EventFilterFieldEnum.SOURCE_NAME]: "fuente",
});

/**
 * Facetas de GET /api/v1/events/facets, en el orden del panel lateral: su titulo y el campo por el que filtran. Las
 * de nombre (agent_name, project_name) filtran por el uuid: el valor de la faceta es el nombre.
 */
export const EVENT_FACETS = Object.freeze([
    { facet: "severity", label: "Gravedad", filterField: EventFilterFieldEnum.SEVERITY, isName: false },
    { facet: "event_type", label: "Tipo de evento", filterField: EventFilterFieldEnum.EVENT_TYPE, isName: false },
    { facet: "source_type", label: "Tipo de fuente", filterField: EventFilterFieldEnum.SOURCE_TYPE, isName: false },
    { facet: "source_name", label: "Fuente", filterField: EventFilterFieldEnum.SOURCE_NAME, isName: false },
    { facet: "agent_name", label: "Agente", filterField: EventFilterFieldEnum.AGENT_UUID, isName: true },
    { facet: "project_name", label: "Proyecto", filterField: EventFilterFieldEnum.PROJECT_UUID, isName: true },
]);
