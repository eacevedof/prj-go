export const SessionEventEnum = Object.freeze({
    // {dashboardToken}: el formulario de entrada o, al recargar, el token guardado (se revalida contra la API)
    LOG_IN_REQUESTED: "session.log-in-requested",
    // {dashboardToken}: el token vale; StartSessionController entra
    SESSION_STARTED: "session.session-started",
    LOG_OUT_REQUESTED: "session.log-out-requested",
    // la API respondio 401 (token cambiado en el .env del backend): cualquier modulo lo dispara
    SESSION_EXPIRED: "session.session-expired",

    USER_LOGGED_IN: "session.user-logged-in",
    USER_LOGGED_OUT: "session.user-logged-out",
});
