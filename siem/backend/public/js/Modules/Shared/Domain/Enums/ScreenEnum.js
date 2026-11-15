/**
 * Pantallas de la SPA: el valor es el primer tramo del hash (#/, #/projects/<uuid>, #/events?..., #/downloads) y el
 * data-screen de su section.
 */
export const ScreenEnum = Object.freeze({
    DASHBOARD: "dashboard",
    PROJECT_DETAIL: "projects",
    EVENTS: "events",
    DOWNLOADS: "downloads",

    DEFAULT: "dashboard",
});
