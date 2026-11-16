import { Dispatcher } from "/js/Modules/Shared/Infrastructure/Components/Dispatcher.js";
import { Htmler } from "/js/Modules/Shared/Infrastructure/Components/Htmler.js";
import { NumberFormatter } from "/js/Modules/Shared/Infrastructure/Components/NumberFormatter.js";

import { DownloadsEventEnum } from "/js/Modules/Downloads/Domain/Enums/DownloadsEventEnum.js";
import { DownloadsOsEnum } from "/js/Modules/Downloads/Domain/Enums/DownloadsOsEnum.js";
import { DownloadsTextEnum } from "/js/Modules/Downloads/Domain/Enums/DownloadsTextEnum.js";
import { InstallInstructionsBuilder } from "/js/Modules/Downloads/Domain/Services/InstallInstructionsBuilder.js";

const CSS_HIDDEN = "hidden";
const PATH_INGEST = "/api/v1/ingest-batches";
const AGENT_MODE_NEW = "new";
const AGENT_MODE_EXISTING = "existing";
const MACOS_SHELL = "zsh";

/**
 * Descargas: asistente para instalar un agente en 3 pasos. 1 - sistema del servidor (Linux, Windows, macOS);
 * 2 - el agente: crearlo aqui mismo (su token entra directo en los comandos) o pegar el token de uno que ya existe;
 * 3 - instalar: la instalacion asistida en una linea (install.sh / install.ps1) y, plegado, el paso a paso manual con
 * "copiar todo". El token solo vive en esta vista: no se guarda ni se envia. Una sola instancia.
 */
export class DownloadsView {

    static #instance = null;

    #dispatcher;
    #htmler;
    #numberFormatter;
    #installInstructionsBuilder;
    #elements = {};
    #endpointUrl = "";
    #downloads = [];
    #installers = [];
    #selectedOs = DownloadsOsEnum.LINUX;
    #agentMode = AGENT_MODE_NEW;
    // el token del agente creado en este asistente (y su nombre): manda sobre el campo de pegar
    #createdAgent = null;
    #isManualOpen = false;

    constructor() {
        this.#dispatcher = Dispatcher.getInstance();
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
            divDownloadsOrigin: document.getElementById("div-downloads-origin"),
            divDownloadsOs: document.getElementById("div-downloads-os"),
            buttonAgentNew: document.getElementById("button-downloads-agent-new"),
            buttonAgentExisting: document.getElementById("button-downloads-agent-existing"),
            formDownloadsAgent: document.getElementById("form-downloads-agent"),
            selectDownloadsProject: document.getElementById("select-downloads-project"),
            inputDownloadsAgentName: document.getElementById("input-downloads-agent-name"),
            buttonDownloadsAgentCreate: document.getElementById("button-downloads-agent-create"),
            pDownloadsAgentErrors: document.getElementById("p-downloads-agent-errors"),
            divDownloadsAgentExisting: document.getElementById("div-downloads-agent-existing"),
            inputDownloadsToken: document.getElementById("input-downloads-token"),
            divDownloadsTokenStatus: document.getElementById("div-downloads-token-status"),
            divDownloadsInstall: document.getElementById("div-downloads-install"),
        };
        this.#elements.divDownloadsOs.addEventListener("click", (event) => this.#onOsClick(event));
        [this.#elements.buttonAgentNew, this.#elements.buttonAgentExisting].forEach((button) => {
            button.addEventListener("click", () => this.#setAgentMode(button.dataset.agentMode));
        });
        this.#elements.formDownloadsAgent.addEventListener("submit", (event) => this.#onCreateAgentSubmit(event));
        this.#elements.inputDownloadsToken.addEventListener("input", () => this.#renderTokenAndInstall());
        // <details> del paso a paso: que un repintado (cambiar de sistema, pegar el token) no lo cierre
        this.#elements.divDownloadsInstall.addEventListener("toggle", (event) => {
            if (event.target instanceof HTMLDetailsElement && event.target.dataset.manual) {
                this.#isManualOpen = event.target.open;
            }
        }, true);
    }

    /**
     * @param {string} endpointUrl la url de ingesta segun el backend; vacia = la de esta misma origin
     * @param {Object[]} downloads binarios tal cual los devuelve la API
     * @param {Object[]} installers instaladores asistidos tal cual los devuelve la API
     */
    renderDownloads(endpointUrl, downloads, installers) {
        this.#endpointUrl = endpointUrl || `${window.location.origin}${PATH_INGEST}`;
        this.#downloads = downloads;
        this.#installers = installers;
        if (!downloads.some((download) => download.os === this.#selectedOs) && downloads.length) {
            this.#selectedOs = downloads[0].os;
        }
        this.renderErrors([]);
        this.#renderOrigin();
        this.#renderOsButtons();
        this.#setAgentMode(this.#agentMode);
    }

    /**
     * @param {{project_uuid: string, project_name: string}[]} projects para crear el agente en uno de ellos
     */
    renderProjects(projects) {
        const h = this.#htmler;
        const selectedProjectUuid = this.#elements.selectDownloadsProject.value;
        this.#elements.selectDownloadsProject.innerHTML = projects.length
            ? projects.map((project) => `<option value="${h.getEscaped(project.project_uuid)}">${h.getEscaped(project.project_name)}</option>`).join("")
            : `<option value="">${h.getEscaped(DownloadsTextEnum.NO_PROJECTS)}</option>`;
        if (projects.some((project) => project.project_uuid === selectedProjectUuid)) {
            this.#elements.selectDownloadsProject.value = selectedProjectUuid;
        }
    }

    /**
     * @param {{agentName: string, apiToken: string, endpointUrl: string}} agent el recien creado
     */
    renderCreatedAgent(agent) {
        this.#createdAgent = { agentName: agent.agentName, apiToken: agent.apiToken };
        if (agent.endpointUrl) this.#endpointUrl = agent.endpointUrl;
        this.#elements.formDownloadsAgent.reset();
        this.renderAgentErrors([]);
        this.#renderTokenAndInstall();
    }

    /**
     * @param {boolean} isBusy
     */
    renderAgentBusy(isBusy) {
        this.#elements.buttonDownloadsAgentCreate.disabled = isBusy;
    }

    /**
     * @param {string[]} errors
     */
    renderAgentErrors(errors) {
        this.#elements.pDownloadsAgentErrors.textContent = errors.join(" · ");
        this.#elements.pDownloadsAgentErrors.classList.toggle(CSS_HIDDEN, errors.length === 0);
    }

    /**
     * @param {string[]} errors
     */
    renderErrors(errors) {
        this.#elements.pDownloadsErrors.textContent = errors.join(" · ");
        this.#elements.pDownloadsErrors.classList.toggle(CSS_HIDDEN, errors.length === 0);
    }

    /**
     * Al salir de sesion: fuera los tokens (el creado y el pegado) y los comandos que los llevaban.
     */
    reset() {
        this.#createdAgent = null;
        this.#elements.inputDownloadsToken.value = "";
        this.#elements.formDownloadsAgent.reset();
        this.#elements.divDownloadsInstall.innerHTML = "";
        this.#elements.divDownloadsTokenStatus.innerHTML = "";
        this.#downloads = [];
        this.#installers = [];
    }

    #onOsClick(event) {
        const button = event.target.closest("[data-os]");
        if (!button) return;
        this.#selectedOs = button.dataset.os;
        this.#renderOsButtons();
        this.#renderInstall();
    }

    #setAgentMode(agentMode) {
        this.#agentMode = agentMode === AGENT_MODE_EXISTING ? AGENT_MODE_EXISTING : AGENT_MODE_NEW;
        const isNew = this.#agentMode === AGENT_MODE_NEW;
        this.#elements.buttonAgentNew.setAttribute("aria-checked", String(isNew));
        this.#elements.buttonAgentExisting.setAttribute("aria-checked", String(!isNew));
        this.#elements.formDownloadsAgent.classList.toggle(CSS_HIDDEN, !isNew);
        this.#elements.divDownloadsAgentExisting.classList.toggle(CSS_HIDDEN, isNew);
        this.#renderTokenAndInstall();
    }

    #onCreateAgentSubmit(event) {
        event.preventDefault();
        this.renderAgentErrors([]);
        this.#dispatcher.dispatch(DownloadsEventEnum.CREATE_AGENT_REQUESTED, {
            projectUuid: this.#elements.selectDownloadsProject.value,
            agentName: this.#elements.inputDownloadsAgentName.value,
        });
    }

    /**
     * El token que va en los comandos: el del agente creado aqui (modo nuevo) o el pegado (modo existente).
     * @returns {string}
     */
    #getApiToken() {
        if (this.#agentMode === AGENT_MODE_NEW) return this.#createdAgent?.apiToken ?? "";
        return this.#elements.inputDownloadsToken.value.trim();
    }

    #renderTokenAndInstall() {
        this.#renderTokenStatus();
        this.#renderInstall();
    }

    #renderTokenStatus() {
        const h = this.#htmler;
        const status = this.#elements.divDownloadsTokenStatus;
        if (this.#agentMode === AGENT_MODE_NEW && this.#createdAgent) {
            status.innerHTML = `<p class="notice"><strong>Agente ${h.getEscaped(this.#createdAgent.agentName)} creado.</strong>
                Su token ya va en los comandos del paso 3. Es la única vez que se ve: si cierras esta página sin instalarlo,
                tendrás que borrar el agente y crear otro.</p>
                <p>${h.getCopyableValue(this.#createdAgent.apiToken)}</p>`;
            return;
        }
        if (this.#getApiToken()) {
            status.innerHTML = `<p class="status status-ok">${h.getEscaped(DownloadsTextEnum.TOKEN_READY)}</p>`;
            return;
        }
        status.innerHTML = `<p class="status">${h.getEscaped(DownloadsTextEnum.TOKEN_MISSING)}</p>`;
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

    #renderOrigin() {
        const h = this.#htmler;
        const serverOrigin = this.#getServerOrigin();
        this.#elements.divDownloadsOrigin.innerHTML = serverOrigin === window.location.origin ? "" :
            `<p class="warning-note">Estás viendo el panel en <code>${h.getEscaped(window.location.origin)}</code>, pero el backend está
                configurado como <code>${h.getEscaped(serverOrigin)}</code> (<code>APP_BASE_URL</code> de su <code>.env</code>). Los
                comandos usan esa dirección: comprueba que los servidores llegan a ella.</p>`;
    }

    #renderOsButtons() {
        const h = this.#htmler;
        if (this.#downloads.length === 0) {
            this.#elements.divDownloadsOs.innerHTML = "";
            return;
        }
        this.#elements.divDownloadsOs.innerHTML = this.#downloads.map((download) => {
            const isSelected = download.os === this.#selectedOs;
            return `<button type="button" class="segment" role="radio" aria-checked="${isSelected}" data-os="${h.getEscaped(download.os)}">
                ${h.getEscaped(DownloadsOsEnum.LABELS[download.os] ?? download.os)} <span class="segment-sub">${h.getEscaped(download.arch)}</span></button>`;
        }).join("");
    }

    #renderInstall() {
        const h = this.#htmler;
        if (this.#downloads.length === 0) {
            this.#elements.divDownloadsInstall.innerHTML = `<p class="empty">Aún no hay binarios publicados. Se generan con
                <code>cd siem/client &amp;&amp; make dist</code>, que los deja en <code>backend/public/downloads/</code> junto a los
                instaladores.</p>`;
            return;
        }
        const download = this.#downloads.find((item) => item.os === this.#selectedOs) ?? this.#downloads[0];
        const serverOrigin = this.#getServerOrigin();
        const apiToken = this.#getApiToken();
        const installer = this.#installInstructionsBuilder.getInstallerCommand(download.os, serverOrigin, apiToken);
        const steps = this.#installInstructionsBuilder.getSteps(download, serverOrigin, this.#endpointUrl, apiToken);
        const installerFiles = this.#installers.filter((item) => String(item.os).split(",").includes(download.os));

        this.#elements.divDownloadsInstall.innerHTML = `
            <section class="panel install-panel">
                <h3>3.1 Instalación asistida <span class="badge badge-ok">recomendada</span></h3>
                <p>Copia y pega en el servidor${download.os === DownloadsOsEnum.WINDOWS ? ", en un PowerShell abierto como administrador" : ""}:</p>
                ${h.getCodeBlock(installer.code, installer.language)}
                <ol class="does">${DownloadsTextEnum.INSTALLER_DOES.map((item) => `<li>${h.getEscaped(item)}</li>`).join("")}</ol>
                ${installerFiles.length ? `<p class="muted">O descarga ${download.os === DownloadsOsEnum.WINDOWS
                    ? "los dos ficheros, déjalos juntos y haz doble clic en <code>install.bat</code>"
                    : "el instalador y lánzalo con <code>sudo sh install.sh</code>"}:
                    ${installerFiles.map((item) => `<a href="${h.getEscaped(item.url)}" download>${h.getEscaped(item.file_name)}</a>`).join(" · ")}</p>` : ""}
            </section>

            <section class="panel">
                <div class="panel-head">
                    <h3>3.2 El binario</h3>
                    <span class="badge badge-version">v${h.getEscaped(download.version)}</span>
                    <a class="button button-small panel-actions" href="${h.getEscaped(download.url)}" download>Descargar ${h.getEscaped(download.file_name)}</a>
                </div>
                <dl class="facts">
                    <dt>Tamaño</dt><dd>${this.#numberFormatter.getBytes(download.size_bytes)}</dd>
                    <dt>sha256</dt><dd>${h.getCopyableValue(download.sha256)}</dd>
                </dl>
            </section>

            <details class="panel manual" data-manual="1" ${this.#isManualOpen ? "open" : ""}>
                <summary><h3>3.3 Paso a paso, sin instalador (${steps.length} pasos)</h3></summary>
                ${h.getCodeBlock(this.#installInstructionsBuilder.getFullScript(steps), `todo junto · ${steps[0]?.language ?? ""}`)}
                <ol class="steps">${steps.map((step, index) => `<li>
                    <h3>3.3.${index + 1} ${h.getEscaped(step.title)}</h3>
                    <p class="muted">${h.getEscaped(step.description)}</p>
                    ${h.getCodeBlock(step.code, step.language)}
                </li>`).join("")}</ol>
            </details>
            ${download.os === DownloadsOsEnum.MACOS ? this.#getHtmlMacosTroubleshooting() : ""}`;
    }

    // macOS: binario compilado sin Mac, sin firmar ni probar en uno: los fallos esperables, a la vista
    #getHtmlMacosTroubleshooting() {
        const h = this.#htmler;
        return `<section class="panel troubleshooting">
            <h3>3.4 ${h.getEscaped(DownloadsTextEnum.MACOS_TROUBLESHOOTING_TITLE)}</h3>
            <p class="muted">${h.getEscaped(DownloadsTextEnum.MACOS_TROUBLESHOOTING_INTRO)}</p>
            <ol class="steps">${DownloadsTextEnum.MACOS_TROUBLESHOOTING.map((item, index) => `<li>
                <h3>3.4.${index + 1} ${h.getEscaped(item.problem)}</h3>
                <p>${h.getEscaped(item.fix)}</p>
                ${item.code ? h.getCodeBlock(item.code, MACOS_SHELL) : ""}
            </li>`).join("")}</ol>
        </section>`;
    }
}
