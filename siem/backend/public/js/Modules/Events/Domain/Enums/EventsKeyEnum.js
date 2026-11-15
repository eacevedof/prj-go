/**
 * Claves de la API de eventos: filtros de GET /api/v1/events, /histogram y /facets (y del hash #/events?...) y sus
 * respuestas. Los campos filtrables (incluir: campo=a,b · excluir: not_campo=a,b) estan en EventFilterFieldEnum.
 */
export const EventsKeyEnum = Object.freeze({
    FROM: "from",
    TO: "to",
    SEARCH: "search",
    PAGE: "page",
    PAGE_SIZE: "page_size",
    NOT_PREFIX: "not_",

    // solo del hash (el front los traduce): rango relativo y auto-refresco
    RANGE: "range",
    REFRESH: "refresh",

    EVENTS: "events",
    HAS_MORE: "has_more",
    BUCKET_SECONDS: "bucket_seconds",
    BUCKETS: "buckets",
    FACETS: "facets",
});
