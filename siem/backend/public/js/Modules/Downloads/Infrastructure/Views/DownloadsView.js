import { Htmler } from "/js/Modules/Shared/Infrastructure/Components/Htmler.js";
import { NumberFormatter } from "/js/Modules/Shared/Infrastructure/Components/NumberFormatter.js";

import { API_TOKEN_PLACEHOLDER, InstallInstructionsBuilder } from "/js/Modules/Downloads/Domain/Services/InstallInstructionsBuilder.js";

const CSS_HIDDEN = "hidden";
const PATH_INGEST = "/api/v1/ingest-batches";
const OS_LABELS = Object.freeze({ linux: "Linux", windows: "Windows", macos: "macOS" });

/**
 * Descargas: los binarios del agente (SO, version, tamaño, sha256 copiable, enlace) y, por cada uno, como instalarlo
 * y configurarlo con los comandos en bloques copiables. El token del agente se puede pegar en un campo para que los
 * comandos salgan ya completos: solo vive en esta vista (no se guarda ni se envia) y por defecto va el marcador
 * <API_TOKEN>. Una sola instancia: guarda las referencias al DOM, lo ultimo pintado y que instrucciones estan abiertas.
 */
export class DownloadsView {

    static #instance = null;

    #htmler;
    #numberFormatter;
    #installInstructionsBuilder;
    #elements = {};
    #endpointUrl = "";
    #downloads = [];
    #openFileNames = new Set();
    #hasRendered = false;

    constructor() {
        this.#htmler = Htmler.getInstance();
        this.#numberFormatter = NumberFormatter.getInstance();
        this.#installInstructionsBuilder = InstallInstructionsBuilder.getInstance();
    }

    /**
     * @returns {DownloadsView}
     */
    static getInstance() {
        DownloadsView.#instance ??= new DownloadsView();
        return DownloadsView.#instance;
    }

    init() {
        this.#elements = {
            pDownloadsErrors: document.getElementById("p-downloads-errors"),
            inputDownloadsToken: document.getElementById("input-downloads-token"),
            divDownloadsList: document.getElementById("div-downloads-list"),
        };
        this.#elements.inputDownloadsToken.addEventListener("input", () => this.#renderList());
        // toggle no burbujea: en captura
        this.#elements.divDownloadsList.addEventListener("toggle", (event) => this.#onDetailsToggle(event), true);
    }

    /**
     * @param {string} endpointUrl la url de ingesta segun el backend; vacia = la de esta misma origin
     * @param {Object[]} downloads tal cual los devuelve la API
     */
    renderDownloads(endpointUrl, downloads) {
        this.#endpointUrl = endpointUrl || `${window.location.origin}${PATH_INGEST}`;
        this.#downloads = downloads;
        // la primera vez, las instrucciones del primer binario abiertas: se ve de un vistazo que hay y como es
        if (!this.#hasRendered && downloads.length) {
            this.#openFileNames.add(downloads[0].file_name);
        }
        this.#hasRendered = true;
        this.renderErrors([]);
        this.#renderList();
    }

    /**
     * @param {string[]} errors
     */
    renderErrors(errors) {
        this.#elements.pDownloadsErrors.textContent = errors.join(" · ");
        this.#elements.pDownloadsErrors.classList.toggle(CSS_HIDDEN, errors.length === 0);
    }

    /**
     * Al salir de sesion: fuera el token pegado y los comandos que lo llevaban.
     */
    reset() {
        this.#elements.inputDownloadsToken.value = "";
        this.#elements.divDownloadsList.innerHTML = "";
        this.#downloads = [];
    }

    #renderList() {
        if (this.#downloads.length === 0) {
            this.#elements.divDownloadsList.innerHTML = `<p class="empty">Aún no hay binarios publicados. Se generan con
                <code>cd siem/client &amp;&amp; make dist</code>, que deja una copia en <code>backend/public/downloads/</code>.</p>`;
            return;
        }
        this.#elements.divDownloadsList.innerHTML = this.#getServerOriginNoticeHtml() + this.#downloads
            .map((download, index) => this.#getDownloadHtml(download, index + 1)).join("");
    }

    /**
     * La direccion que usaran los servidores: la del backend (APP_BASE_URL, de donde sale la url de ingesta). La
     * descarga y el .env usan la misma, aunque el panel se este viendo por otra (un tunel, localhost).
     * @returns {string} p. ej. https://siem.dominio.com
     */
    #getServerOrigin() {
        try {
            return new URL(this.#endpointUrl).origin;
        }
        catch (error) {
            return window.location.origin;
        }
    }

    #getServerOriginNoticeHtml() {
        const h = this.#htmler;
        const serverOrigin = this.#getServerOrigin();
        if (serverOrigin === window.location.origin) return "";
        return `<p class="notice">Estás viendo el panel en <code>${h.getEscaped(window.location.origin)}</code>, pero el backend
            está configurado como <code>${h.getEscaped(serverOrigin)}</code> (<code>APP_BASE_URL</code> de su <code>.env</code>).
            Los comandos de abajo usan esa dirección: comprueba que los servidores donde instales el agente llegan a ella.</p>`;
    }

    #getDownloadHtml(download, position) {
        const h = this.#htmler;
        const fileName = String(download.file_name ?? "");
        const steps = this.#installInstructionsBuilder.getSteps(
            download, this.#getServerOrigin(), this.#endpointUrl, this.#elements.inputDownloadsToken.value
        );
        return `<section class="panel download">
            <div class="panel-head">
                <span class="download-name">1.${position} ${h.getEscaped(OS_LABELS[download.os] ?? download.os)}
                    <span class="muted">${h.getEscaped(download.arch)}</span></span>
                <span class="badge badge-version">v${h.getEscaped(download.version)}</span>
                <a class="button button-primary button-small" href="${h.getEscaped(download.url)}" download>Descargar</a>
            </div>
            <dl class="facts">
                <dt>Fichero</dt><dd class="mono">${h.getEscaped(fileName)}</dd>
                <dt>Tamaño</dt><dd>${this.#numberFormatter.getBytes(download.size_bytes)}</dd>
                <dt>sha256</dt><dd>${h.getCopyableValue(download.sha256)}</dd>
            </dl>
            <details class="install" data-file-name="${h.getEscaped(fileName)}" ${this.#openFileNames.has(fileName) ? "open" : ""}>
                <summary>Cómo instalarlo y configurarlo (${steps.length} pasos)</summary>
                ${steps.length === 0
                    ? `<p class="empty">Sin instrucciones para este sistema operativo.</p>`
                    : `<ol class="steps">${steps.map((step, stepIndex) => `<li>
                        <h3>1.${position}.${stepIndex + 1} ${h.getEscaped(step.title)}</h3>
                        <p class="muted">${h.getEscaped(step.description)}</p>
                        ${h.getCodeBlock(step.code, step.language)}
                    </li>`).join("")}</ol>`}
                ${this.#elements.inputDownloadsToken.value.trim() === ""
                    ? `<p class="muted">Sustituye <code>${h.getEscaped(API_TOKEN_PLACEHOLDER)}</code> por el token del agente, o pégalo arriba para que los comandos salgan completos.</p>`
                    : ""}
            </details>
        </section>`;
    }

    #onDetailsToggle(event) {
        const details = event.target;
        if (!(details instanceof HTMLDetailsElement) || !details.dataset.fileName) return;
        if (details.open) {
            this.#openFileNames.add(details.dataset.fileName);
            return;
        }
        this.#openFileNames.delete(details.dataset.fileName);
    }
}
