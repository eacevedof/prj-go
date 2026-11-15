const SESSION_KEY = "siem-dashboard-session";

/**
 * Sesion del dashboard en sessionStorage (se borra al cerrar la pestaña): solo el token del dashboard
 * (APP_DASHBOARD_TOKEN del backend). Como todo componente no captura: si el navegador bloquea sessionStorage el
 * error sube al controller.
 */
export class SessionStorer {

    /**
     * @returns {SessionStorer}
     */
    static getInstance() {
        return new SessionStorer();
    }

    /**
     * @returns {{dashboardToken: string}|null}
     */
    getSession() {
        const session = JSON.parse(sessionStorage.getItem(SESSION_KEY) ?? "null");
        return session?.dashboardToken ? session : null;
    }

    /**
     * @returns {string}
     */
    getDashboardToken() {
        return this.getSession()?.dashboardToken ?? "";
    }

    /**
     * @param {{dashboardToken: string}} session
     */
    saveSession(session) {
        sessionStorage.setItem(SESSION_KEY, JSON.stringify(session));
    }

    removeSession() {
        sessionStorage.removeItem(SESSION_KEY);
    }
}
