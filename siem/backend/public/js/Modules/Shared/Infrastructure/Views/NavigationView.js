import { ScreenEnum } from "/js/Modules/Shared/Domain/Enums/ScreenEnum.js";
import { SharedEventEnum } from "/js/Modules/Shared/Domain/Enums/SharedEventEnum.js";
import { Dispatcher } from "/js/Modules/Shared/Infrastructure/Components/Dispatcher.js";

const CSS_HIDDEN = "hidden";
const CSS_ACTIVE = "is-active";
// el detalle de un proyecto cuelga de la portada: su enlace del menu sigue marcado
const NAV_SCREEN_BY_SCREEN = Object.freeze({
    [ScreenEnum.DASHBOARD]: ScreenEnum.DASHBOARD,
    [ScreenEnum.PROJECT_DETAIL]: ScreenEnum.DASHBOARD,
    [ScreenEnum.EVENTS]: ScreenEnum.EVENTS,
    [ScreenEnum.DOWNLOADS]: ScreenEnum.DOWNLOADS,
});

/**
 * Menu y router por hash: #/ (portada), #/projects/<project_uuid> (detalle), #/events?... (visor), #/downloads
 * (binarios del agente). Cada pantalla es una section con data-screen; solo se ve la del hash y solo con sesion
 * (show/hide). Pintarla es cosa de su controller: aqui solo se pide (SCREEN_REQUESTED). Una sola instancia: guarda
 * las referencias al DOM.
 */
export class NavigationView {

    static #instance = null;

    #dispatcher;
    #elements = {};
    #isVisible = false;

    constructor() {
        this.#dispatcher = Dispatcher.getInstance();
    }

    /**
     * @returns {NavigationView}
     */
    static getInstance() {
        NavigationView.#instance ??= new NavigationView();
        return NavigationView.#instance;
    }

    init() {
        this.#elements = {
            navScreens: document.getElementById("nav-screens"),
            screenLinks: [...document.querySelectorAll("#nav-screens a[data-screen]")],
            screens: [...document.querySelectorAll("section[data-screen]")],
            spanRefreshInfo: document.getElementById("span-refresh-info"),
        };
        window.addEventListener("hashchange", () => this.renderCurrentScreen());
    }

    /**
     * Con sesion: menu visible y la pantalla del hash (o la de por defecto).
     */
    show() {
        this.#isVisible = true;
        this.#elements.navScreens.classList.remove(CSS_HIDDEN);
        this.renderCurrentScreen();
    }

    hide() {
        this.#isVisible = false;
        this.#elements.navScreens.classList.add(CSS_HIDDEN);
        this.#elements.screens.forEach((screen) => screen.classList.add(CSS_HIDDEN));
        this.#elements.spanRefreshInfo.textContent = "";
    }

    /**
     * Muestra la section del hash, marca su enlace y pide a su controller que la pinte (tambien en cada refresco y
     * despues de un alta o una baja).
     */
    renderCurrentScreen() {
        if (!this.#isVisible) return;

        const route = this.#getRoute();
        this.#elements.screens.forEach((screen) => screen.classList.toggle(CSS_HIDDEN, screen.dataset.screen !== route.screen));
        this.#elements.screenLinks.forEach((link) => link.classList.toggle(CSS_ACTIVE, link.dataset.screen === NAV_SCREEN_BY_SCREEN[route.screen]));
        this.#dispatcher.dispatch(SharedEventEnum.SCREEN_REQUESTED, route);
    }

    /**
     * Ir a otra pantalla. Con el mismo hash no hay hashchange: se repinta igual.
     * @param {string} hash p. ej. "#/" o "#/events?severity=error"
     */
    goToHash(hash) {
        if (window.location.hash === hash) {
            this.renderCurrentScreen();
            return;
        }
        window.location.hash = hash;
    }

    /**
     * @param {boolean} isAutoRefreshed
     * @param {string} renderedAt "YYYY-MM-DD HH:MM:SS" UTC del ultimo pintado
     * @param {number} intervalSeconds
     */
    renderRefreshInfo(isAutoRefreshed, renderedAt, intervalSeconds) {
        this.#elements.spanRefreshInfo.textContent = isAutoRefreshed
            ? `${renderedAt} UTC · cada ${intervalSeconds} s`
            : `${renderedAt} UTC · sin refresco`;
    }

    /**
     * "#/projects/abc?x=1" -> {screen: "projects", pathParam: "abc", query: {x: "1"}}; hash vacio o desconocido -> portada.
     * @returns {{screen: string, pathParam: string, query: Object}}
     */
    #getRoute() {
        const [path, queryString = ""] = window.location.hash.replace(/^#\/?/, "").split("?");
        const [firstSegment = "", pathParam = ""] = path.split("/").map((segment) => decodeURIComponent(segment));
        const isKnownScreen = Object.values(ScreenEnum).includes(firstSegment);
        return {
            screen: isKnownScreen ? firstSegment : ScreenEnum.DEFAULT,
            pathParam,
            query: Object.fromEntries(new URLSearchParams(queryString)),
        };
    }
}
