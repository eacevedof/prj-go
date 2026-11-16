/**
 * Pantallas de la SPA: el valor es el data-screen de su section. El hash elige la pantalla por su primer tramo:
 * #/events?... (visor, la de entrada), #/projects (lista de proyectos = DASHBOARD), #/projects/<uuid> (detalle =
 * PROJECT_DETAIL), #/downloads. "#/" o un tramo desconocido = DEFAULT.
 */
export const ScreenEnum = Object.freeze({
    DASHBOARD: "dashboard",
    PROJECT_DETAIL: "projects",
    EVENTS: "events",
    DOWNLOADS: "downloads",

    DEFAULT: "events",
});
