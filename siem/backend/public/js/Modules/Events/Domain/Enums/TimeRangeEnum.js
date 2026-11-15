/**
 * Atajos del selector de rango de tiempo. En el hash viajan como relativos (range=24h) y se traducen a from/to en cada
 * peticion: la url compartida sigue siendo "las ultimas 24 h", no un dia fijo. Un rango a medida va como from/to.
 */
export const TimeRangeEnum = Object.freeze({
    MINUTES_15: "15m",
    HOUR_1: "1h",
    HOURS_6: "6h",
    HOURS_24: "24h",
    DAYS_7: "7d",
    DAYS_30: "30d",

    DEFAULT: "24h",
    // valor del <select> para "rango a medida" (no viaja en el hash: van from y to)
    CUSTOM: "custom",
});

export const TIME_RANGE_SECONDS = Object.freeze({
    [TimeRangeEnum.MINUTES_15]: 15 * 60,
    [TimeRangeEnum.HOUR_1]: 60 * 60,
    [TimeRangeEnum.HOURS_6]: 6 * 60 * 60,
    [TimeRangeEnum.HOURS_24]: 24 * 60 * 60,
    [TimeRangeEnum.DAYS_7]: 7 * 24 * 60 * 60,
    [TimeRangeEnum.DAYS_30]: 30 * 24 * 60 * 60,
});

export const TIME_RANGE_LABELS = Object.freeze({
    [TimeRangeEnum.MINUTES_15]: "Últimos 15 min",
    [TimeRangeEnum.HOUR_1]: "Última hora",
    [TimeRangeEnum.HOURS_6]: "Últimas 6 h",
    [TimeRangeEnum.HOURS_24]: "Últimas 24 h",
    [TimeRangeEnum.DAYS_7]: "Últimos 7 días",
    [TimeRangeEnum.DAYS_30]: "Últimos 30 días",
    [TimeRangeEnum.CUSTOM]: "Rango a medida (UTC)",
});
