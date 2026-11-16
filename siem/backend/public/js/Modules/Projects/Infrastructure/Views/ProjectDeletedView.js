import { Htmler } from "/js/Modules/Shared/Infrastructure/Components/Htmler.js";
import { TextFormatter } from "/js/Modules/Shared/Infrastructure/Components/TextFormatter.js";

import { ClientOsEnum } from "/js/Modules/Projects/Domain/Enums/ClientOsEnum.js";
import { ProjectsTextEnum } from "/js/Modules/Projects/Domain/Enums/ProjectsTextEnum.js";
import { ClientUninstallCommandBuilder } from "/js/Modules/Projects/Domain/Services/ClientUninstallCommandBuilder.js";

const CSS_HIDDEN = "hidden";
const ACTION_CLOSE = "close";

/**
 * Aviso de la portada tras borrar un proyecto: el borrado no quita los siem-client de sus servidores (solo dejan de
 * ser aceptados), asi que dice, agente a agente, en que servidor esta, la ruta de su ejecutable y como desinstalarlo.
 * Se queda hasta que se pulsa "Entendido": el auto-refresco no lo toca.
 * Una sola instancia: guarda las referencias al DOM.
 */
export class ProjectDeletedView {

    static #instance = null;

    #htmler;
    #textFormatter;
    #clientUninstallCommandBuilder;
    #elements = {};

    constructor() {
        this.#htmler = Htmler.getInstance();
        this.#textFormatter = TextFormatter.getInstance();
        this.#clientUninstallCommandBuilder = ClientUninstallCommandBuilder.getInstance();
    }

    /**
     * @returns {ProjectDeletedView}
     */
    static getInstance() {
        ProjectDeletedView.#instance ??= new ProjectDeletedView();
        return ProjectDeletedView.#instance;
    }

    init() {
        this.#elements = {
            divProjectDeleted: document.getElementById("div-project-deleted"),
        };
        this.#elements.divProjectDeleted.addEventListener("click", (event) => {
            if (event.target.closest(`button[data-action="${ACTION_CLOSE}"]`)) this.hide();
        });
    }

    /**
     * @param {string} projectName
     * @param {{agentName: string, hostName: string, hostOs: string, hostIp: string, clientPath: string}[]} agents
     */
    renderDeletedProject(projectName, agents) {
        const h = this.#htmler;
        this.#elements.divProjectDeleted.innerHTML = `
            <p><strong>${h.getEscaped(this.#textFormatter.getFormatted(ProjectsTextEnum.DELETED_TITLE, projectName))}</strong></p>
            ${agents.length ? this.#getHtmlAgents(agents) : `<p>${h.getEscaped(ProjectsTextEnum.DELETED_WITHOUT_AGENTS)}</p>`}
            <p><button type="button" class="button button-small" data-action="${ACTION_CLOSE}">${h.getEscaped(ProjectsTextEnum.DELETED_CLOSE)}</button></p>`;
        this.#elements.divProjectDeleted.classList.remove(CSS_HIDDEN);
    }

    hide() {
        this.#elements.divProjectDeleted.classList.add(CSS_HIDDEN);
        this.#elements.divProjectDeleted.innerHTML = "";
    }

    #getHtmlAgents(agents) {
        const h = this.#htmler;
        const uninstallText = agents.length === 1
            ? ProjectsTextEnum.DELETED_UNINSTALL_ONE
            : this.#textFormatter.getFormatted(ProjectsTextEnum.DELETED_UNINSTALL_MANY, agents.length);
        return `<p>${h.getEscaped(uninstallText)}</p>
            <ol class="uninstall-list">${agents.map((agent) => this.#getHtmlAgent(agent)).join("")}</ol>`;
    }

    #getHtmlAgent(agent) {
        const h = this.#htmler;
        // nunca envio nada: ni SO ni ruta; solo donde lo habria dejado el instalador, sin comando a ciegas
        if (!agent.hostOs) {
            return `<li>${this.#getHtmlAgentTitle(agent.agentName, ProjectsTextEnum.DELETED_HOST_UNKNOWN)}
                <div class="muted">${h.getEscaped(this.#textFormatter.getFormatted(
                    ProjectsTextEnum.DELETED_NEVER_SEEN, ClientOsEnum.DEFAULT_PATH_LINUX, ClientOsEnum.DEFAULT_PATH_WINDOWS
                ))}</div>
            </li>`;
        }
        const commandTitle = agent.hostOs === ClientOsEnum.WINDOWS
            ? ProjectsTextEnum.DELETED_COMMAND_WINDOWS
            : ProjectsTextEnum.DELETED_COMMAND_UNIX;
        const pathHint = agent.clientPath ? "" : ` <span class="muted">(${h.getEscaped(ProjectsTextEnum.DELETED_PATH_DEFAULT)})</span>`;
        return `<li>${this.#getHtmlAgentTitle(agent.agentName, [agent.hostName, agent.hostOs, agent.hostIp].filter(Boolean).join(" · "))}
            <div>${h.getCopyableValue(this.#clientUninstallCommandBuilder.getClientPathOrDefault(agent.hostOs, agent.clientPath))}${pathHint}</div>
            ${h.getCodeBlock(this.#clientUninstallCommandBuilder.getUninstallCommand(agent.hostOs, agent.clientPath), commandTitle)}
        </li>`;
    }

    #getHtmlAgentTitle(agentName, hostLabel) {
        return `<div><strong>${this.#htmler.getEscaped(
            this.#textFormatter.getFormatted(ProjectsTextEnum.DELETED_AGENT_ON_HOST, agentName, hostLabel)
        )}</strong></div>`;
    }
}
