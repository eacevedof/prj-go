import { BreadcrumbsTextEnum } from "/js/Modules/Shared/Domain/Enums/BreadcrumbsTextEnum.js";
import { RouteHashEnum } from "/js/Modules/Shared/Domain/Enums/RouteHashEnum.js";
import { ScreenEnum } from "/js/Modules/Shared/Domain/Enums/ScreenEnum.js";
import { SharedEventEnum } from "/js/Modules/Shared/Domain/Enums/SharedEventEnum.js";
import { Dispatcher } from "/js/Modules/Shared/Infrastructure/Components/Dispatcher.js";
import { Htmler } from "/js/Modules/Shared/Infrastructure/Components/Htmler.js";
import { TextFormatter } from "/js/Modules/Shared/Infrastructure/Components/TextFormatter.js";

const CSS_HIDDEN = "hidden";
const CSS_ACTIVE = "is-active";
// pie de la barra: {0} = ultimo pintado (UTC), {1} = segundos entre repintados
const REFRESH_INFO_AUTO = "{0} UTC · cada {1} s";
const REFRESH_INFO_OFF = "{0} UTC · sin refresco";
// migas de cada pantalla al entrar en ella (el detalle de un proyecto las completa con su nombre)
const BREADCRUMBS_BY_SCREEN = Object.freeze({
    [ScreenEnum.DASHBOARD]: Object.freeze([{ label: BreadcrumbsTextEnum.PROJECTS }]),
    [ScreenEnum.PROJECT_DETAIL]: Object.freeze([
        { label: BreadcrumbsTextEnum.PROJECTS, hash: RouteHashEnum.PROJECTS },
        { label: BreadcrumbsTextEnum.PROJECT_LOADING },
    ]),
    [ScreenEnum.EVENTS]: Object.freeze([{ label: BreadcrumbsTextEnum.EVENTS }]),
    [ScreenEnum.DOWNLOADS]: Object.freeze([{ label: BreadcrumbsTextEnum.DOWNLOADS }]),
});
// el detalle de un proyecto cuelga de la portada: su enlace del menu sigue marcado
const NAV_SCREEN_BY_SCREEN = Object.freeze({
    [ScreenEnum.DASHBOARD]: ScreenEnum.DASHBOARD,
    [ScreenEnum.PROJECT_DETAIL]: ScreenEnum.DASHBOARD,
    [ScreenEnum.EVENTS]: ScreenEnum.EVENTS,
    [ScreenEnum.DOWNLOADS]: ScreenEnum.DOWNLOADS,
});

/**
 * Menu y router por hash: #/events?... (visor, la de entrada), #/projects (lista), #/projects/<project_uuid>
 * (detalle), #/downloads (binarios del agente). Cada pantalla es una section con data-screen; solo se ve la del hash y solo con sesion
 * (show/hide). Pintarla es cosa de su controller: aqui solo se pide (SCREEN_REQUESTED). Tambien las migas de pan:
 * las de cada pantalla al entrar; quien sepa mas (el detalle de un proyecto, su nombre) las completa con
 * renderBreadcrumbs. Una sola instancia: guarda las referencias al DOM y la ultima ruta pintada.
 */
export class NavigationView {

    static #instance = null;

    #dispatcher;
    #htmler;
    #textFormatter;
    #elements = {};
    #isVisible = false;
    // "pantalla/pathParam" de la ultima ruta pintada: el auto-refresco de la misma no resetea las migas
    #lastRouteKey = "";

    constructor() {
        this.#dispatcher = Dispatcher.getInstance();
        this.#htmler = Htmler.getInstance();
        this.#textFormatter = TextFormatter.getInstance();
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
            navBreadcrumbs: document.getElementById("nav-breadcrumbs"),
        };
        this.#elements.navBreadcrumbs.setAttribute("aria-label", BreadcrumbsTextEnum.ARIA_LABEL);
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
        this.#elements.navBreadcrumbs.classList.add(CSS_HIDDEN);
        this.#lastRouteKey = "";
    }

    /**
     * Migas de la pantalla visible: todas enlazan menos la ultima (donde se esta).
     * @param {{label: string, hash?: string}[]} breadcrumbs
     */
    renderBreadcrumbs(breadcrumbs) {
        const h = this.#htmler;
        const lastPosition = breadcrumbs.length - 1;
        this.#elements.navBreadcrumbs.innerHTML = `<ol>${breadcrumbs.map((breadcrumb, position) => {
            if (position === lastPosition || !breadcrumb.hash) {
                return `<li${position === lastPosition ? ` aria-current="page"` : ""}>${h.getEscaped(breadcrumb.label)}</li>`;
            }
            return `<li><a href="${h.getEscaped(breadcrumb.hash)}">${h.getEscaped(breadcrumb.label)}</a></li>`;
        }).join("")}</ol>`;
        this.#elements.navBreadcrumbs.classList.toggle(CSS_HIDDEN, breadcrumbs.length === 0);
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
        const routeKey = `${route.screen}/${route.pathParam}`;
        if (routeKey !== this.#lastRouteKey) {
            this.#lastRouteKey = routeKey;
            this.renderBreadcrumbs(BREADCRUMBS_BY_SCREEN[route.screen] ?? []);
        }
        this.#dispatcher.dispatch(SharedEventEnum.SCREEN_REQUESTED, route);
    }

    /**
     * Ir a otra pantalla. Con el mismo hash no hay hashchange: se repinta igual.
     * @param {string} hash p. ej. "#/projects" o "#/events?severity=error"
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
            ? this.#textFormatter.getFormatted(REFRESH_INFO_AUTO, renderedAt, intervalSeconds)
            : this.#textFormatter.getFormatted(REFRESH_INFO_OFF, renderedAt);
    }

    /**
     * "#/projects/abc?x=1" -> {screen: "projects", pathParam: "abc", query: {x: "1"}}; "#/projects" a secas -> la lista
     * (dashboard); hash vacio o desconocido -> la de entrada (eventos).
     * @returns {{screen: string, pathParam: string, query: Object}}
     */
    #getRoute() {
        const [path, queryString = ""] = window.location.hash.replace(/^#\/?/, "").split("?");
        const [firstSegment = "", pathParam = ""] = path.split("/").map((segment) => decodeURIComponent(segment));
        const isKnownScreen = Object.values(ScreenEnum).includes(firstSegment);
        const isProjectsList = firstSegment === ScreenEnum.PROJECT_DETAIL && pathParam === "";
        return {
            screen: isProjectsList ? ScreenEnum.DASHBOARD : (isKnownScreen ? firstSegment : ScreenEnum.DEFAULT),
            pathParam,
            query: Object.fromEntries(new URLSearchParams(queryString)),
        };
    }
}
