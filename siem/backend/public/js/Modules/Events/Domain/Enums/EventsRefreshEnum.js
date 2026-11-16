/**
 * Auto-refresco del visor de eventos (parametro refresh del hash): apagado o cada N segundos.
 */
export const EventsRefreshEnum = Object.freeze({
    OFF: "off",
    SECONDS_10: "10",
    SECONDS_30: "30",
    SECONDS_60: "60",

    DEFAULT: "10",
});

export const EVENTS_REFRESH_VALUES = Object.freeze([
    EventsRefreshEnum.OFF, EventsRefreshEnum.SECONDS_10, EventsRefreshEnum.SECONDS_30, EventsRefreshEnum.SECONDS_60,
]);

export const EVENTS_REFRESH_LABELS = Object.freeze({
    [EventsRefreshEnum.OFF]: "Refresco: apagado",
    [EventsRefreshEnum.SECONDS_10]: "Refresco: 10 s",
    [EventsRefreshEnum.SECONDS_30]: "Refresco: 30 s",
    [EventsRefreshEnum.SECONDS_60]: "Refresco: 60 s",
});
