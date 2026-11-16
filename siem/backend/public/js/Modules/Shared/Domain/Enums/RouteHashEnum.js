/**
 * Hash de entrada de cada pantalla del menu (el detalle de un proyecto es PROJECTS + "/<project_uuid>"). La de
 * entrada es Eventos: lo del dia a dia es revisar eventos; "#/" y un hash desconocido llevan ahi (ScreenEnum.DEFAULT).
 */
export const RouteHashEnum = Object.freeze({
    EVENTS: "#/events",
    PROJECTS: "#/projects",
    DOWNLOADS: "#/downloads",
});
