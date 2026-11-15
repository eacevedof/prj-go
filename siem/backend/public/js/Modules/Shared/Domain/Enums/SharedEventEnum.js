export const SharedEventEnum = Object.freeze({
    // {screen, pathParam, query}: la url (hash) pide una pantalla; el controller de esa pantalla la pinta
    SCREEN_REQUESTED: "shared.screen-requested",
    // {isAutoRefreshed, intervalSeconds?}: la pantalla ya esta pintada; RefreshScreenController programa el siguiente
    // repintado (sin intervalSeconds, el de RefreshEnum)
    SCREEN_RENDERED: "shared.screen-rendered",
});
