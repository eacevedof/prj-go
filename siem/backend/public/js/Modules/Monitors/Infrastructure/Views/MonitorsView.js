import { DateTimer } from "/js/Modules/Shared/Infrastructure/Components/DateTimer.js";
import { Dispatcher } from "/js/Modules/Shared/Infrastructure/Components/Dispatcher.js";
import { Htmler } from "/js/Modules/Shared/Infrastructure/Components/Htmler.js";
import { TextFormatter } from "/js/Modules/Shared/Infrastructure/Components/TextFormatter.js";

import { MonitorDuplicateReasonEnum } from "/js/Modules/Monitors/Domain/Enums/MonitorDuplicateReasonEnum.js";
import { MonitorsEventEnum } from "/js/Modules/Monitors/Domain/Enums/MonitorsEventEnum.js";
import { MonitorsTextEnum } from "/js/Modules/Monitors/Domain/Enums/MonitorsTextEnum.js";
import {
    METRIC_MONITOR_TYPES, MONITOR_SEVERITIES, MONITOR_TYPE_LABELS, MONITOR_TYPES, MonitorDefaultEnum, MonitorTypeEnum,
    PERCENT_MONITOR_TYPES,
} from "/js/Modules/Monitors/Domain/Enums/MonitorTypeEnum.js";
import { MonitorDuplicateFinder } from "/js/Modules/Monitors/Domain/Services/MonitorDuplicateFinder.js";

const CSS_HIDDEN = "hidden";
const ACTION_EDIT = "edit";
const ACTION_TOGGLE = "toggle";
const ACTION_ASK_DELETE = "ask-delete";
const ACTION_CONFIRM_DELETE = "confirm-delete";
const ACTION_CANCEL_DELETE = "cancel-delete";
const UNIT_PERCENT = "%";
const UNIT_MBIT = "Mbit/s";
// ejemplo de ruta del campo path en los monitores de lineas (en disco, la ruta por defecto)
const LINES_PATH_PLACEHOLDER = "/var/log/nginx";
// que campos de config enseña el formulario para cada tipo
const CONFIG_FIELDS_BY_TYPE = Object.freeze({
    [MonitorTypeEnum.DISK]: ["path"],
    [MonitorTypeEnum.NETWORK]: ["interface"],
    [MonitorTypeEnum.FOLDER]: ["path", "pattern"],
    [MonitorTypeEnum.LOG_PATTERN]: ["path", "pattern", "match"],
});

/**
 * "Monitores" del detalle de un proyecto: los recursos a vigilar (cpu, memoria, disco, gpu, red, carpetas, patrones
 * de log). Tabla con su monitor_uuid (copiable), ambito, umbrales y ultimo valor por agente; y un modal (<dialog>, el
 * mismo para alta y edicion) cuyo bloque de config cambia con el tipo. El ambito lo fija desde donde se abre: la fila
 * de un agente ("Crear monitor"), el boton de "todos los agentes" o "Editar" (el que ya tenia). Mientras se escribe
 * avisa si el monitor estaria repetido (MonitorDuplicateFinder) y no deja enviarlo. El refresco solo repinta la tabla:
 * lo que se este escribiendo no se pierde. Borrar pide un segundo clic en la propia fila.
 * Una sola instancia: guarda las referencias al DOM, el proyecto y los monitores pintados.
 */
export class MonitorsView {

    static #instance = null;

    #dateTimer;
    #dispatcher;
    #htmler;
    #textFormatter;
    #monitorDuplicateFinder;
    #elements = {};
    #projectUuid = "";
    #monitorsByUuid = new Map();
    #pendingDeleteMonitorUuid = "";
    #isBusy = false;
    #hasDuplicate = false;

    constructor() {
        this.#dateTimer = DateTimer.getInstance();
        this.#dispatcher = Dispatcher.getInstance();
        this.#htmler = Htmler.getInstance();
        this.#textFormatter = TextFormatter.getInstance();
        this.#monitorDuplicateFinder = MonitorDuplicateFinder.getInstance();
    }

    /**
     * @returns {MonitorsView}
     */
    static getInstance() {
        MonitorsView.#instance ??= new MonitorsView();
        return MonitorsView.#instance;
    }

    init() {
        this.#loadElements();
        this.#loadSelectOptions();
        this.#renderFormByType();
        this.#elements.selectMonitorType.addEventListener("change", () => this.#renderFormByType());
        this.#elements.formMonitor.addEventListener("input", () => this.#renderDuplicate());
        this.#elements.formMonitor.addEventListener("change", () => this.#renderDuplicate());
        this.#elements.formMonitor.addEventListener("submit", (event) => this.#onSubmit(event));
        this.#elements.buttonMonitorCancel.addEventListener("click", () => this.closeDialog());
        this.#elements.buttonMonitorCreateAll.addEventListener("click", () => this.openCreateDialog("", ""));
        this.#elements.divMonitorsTable.addEventListener("click", (event) => this.#onTableClick(event));
    }

    /**
     * El proyecto que se esta viendo. Al cambiar de proyecto, tabla y modal se vacian; en un refresco, no.
     * @param {string} projectUuid
     */
    renderProject(projectUuid) {
        if (projectUuid === this.#projectUuid) return;
        this.#projectUuid = projectUuid;
        this.#pendingDeleteMonitorUuid = "";
        this.#monitorsByUuid = new Map();
        this.#elements.divMonitorsTable.innerHTML = "";
        this.closeDialog();
        this.renderMessage("");
        this.renderErrors([]);
    }

    /**
     * Modal de alta. agentUuid vacio = para todos los agentes del proyecto.
     * @param {string} agentUuid
     * @param {string} agentName
     */
    openCreateDialog(agentUuid, agentName) {
        this.resetForm();
        this.#renderScope(agentUuid, agentName);
        this.#elements.h2MonitorForm.textContent = agentUuid
            ? this.#textFormatter.getFormatted(MonitorsTextEnum.FORM_TITLE_CREATE_IN_AGENT, agentName)
            : MonitorsTextEnum.FORM_TITLE_CREATE_FOR_ALL_AGENTS;
        this.renderMessage("");
        this.#showDialog();
        this.#elements.selectMonitorType.focus();
    }

    closeDialog() {
        if (typeof this.#elements.dialogMonitor.close === "function") {
            this.#elements.dialogMonitor.close();
            return;
        }
        this.#elements.dialogMonitor.removeAttribute("open");
    }

    /**
     * @param {Object[]} monitors tal cual los devuelve la API (claves snake_case, con states[])
     */
    renderMonitors(monitors) {
        this.#monitorsByUuid = new Map(monitors.map((monitor) => [monitor.monitor_uuid, monitor]));
        if (!this.#monitorsByUuid.has(this.#pendingDeleteMonitorUuid)) {
            this.#pendingDeleteMonitorUuid = "";
        }
        this.renderErrors([]);
        this.#renderTable();
        // otro pudo crear o borrar uno mientras el modal estaba abierto
        if (this.#elements.dialogMonitor.open) this.#renderDuplicate();
    }

    /**
     * Aviso del camino feliz (creado, guardado), encima de la tabla. Vacio lo oculta.
     * @param {string} message
     */
    renderMessage(message) {
        this.#renderText(this.#elements.pMonitorsMessage, message);
    }

    /**
     * @param {string[]} errors de la lista (no se pudo cargar, o fallo activar / borrar)
     */
    renderErrors(errors) {
        this.#renderText(this.#elements.pMonitorsErrors, errors.join(" · "));
    }

    /**
     * @param {string[]} errors del formulario (los de la API, como el 409 de repetido, incluidos)
     */
    renderFormErrors(errors) {
        this.#renderText(this.#elements.pMonitorFormErrors, errors.join(" · "));
    }

    /**
     * Guardado: se cierra el modal y se avisa encima de la tabla.
     * @param {string} message
     */
    renderSaved(message) {
        this.closeDialog();
        this.resetForm();
        this.renderMessage(message);
    }

    /**
     * @param {boolean} isBusy
     */
    renderBusy(isBusy) {
        this.#isBusy = isBusy;
        this.#renderSubmitState();
    }

    resetForm() {
        this.#elements.formMonitor.reset();
        this.#elements.hiddenMonitorUuid.value = "";
        this.#elements.selectMonitorType.disabled = false;
        this.#elements.inputMonitorInterval.value = String(MonitorDefaultEnum.INTERVAL_SECONDS);
        this.#elements.checkboxMonitorEnabled.checked = true;
        this.#elements.h2MonitorForm.textContent = MonitorsTextEnum.FORM_TITLE_CREATE;
        this.#elements.buttonMonitorSubmit.textContent = MonitorsTextEnum.SUBMIT_CREATE;
        this.#renderScope("", "");
        this.renderFormErrors([]);
        this.#renderFormByType();
    }

    #loadElements() {
        this.#elements = {
            buttonMonitorCreateAll: document.getElementById("button-monitor-create-all"),
            pMonitorsMessage: document.getElementById("p-monitors-message"),
            pMonitorsErrors: document.getElementById("p-monitors-errors"),
            divMonitorsTable: document.getElementById("div-monitors-table"),
            dialogMonitor: document.getElementById("dialog-monitor"),
            h2MonitorForm: document.getElementById("h2-monitor-form"),
            formMonitor: document.getElementById("form-monitor"),
            hiddenMonitorUuid: document.getElementById("hidden-monitor-uuid"),
            hiddenMonitorAgentUuid: document.getElementById("hidden-monitor-agent-uuid"),
            strongMonitorScope: document.getElementById("strong-monitor-scope"),
            selectMonitorType: document.getElementById("select-monitor-type"),
            inputMonitorName: document.getElementById("input-monitor-name"),
            checkboxMonitorEnabled: document.getElementById("checkbox-monitor-enabled"),
            divMonitorMetricFields: document.getElementById("div-monitor-metric-fields"),
            inputMonitorInterval: document.getElementById("input-monitor-interval"),
            inputMonitorWarning: document.getElementById("input-monitor-warning"),
            inputMonitorCritical: document.getElementById("input-monitor-critical"),
            spanMonitorThresholdUnit: document.getElementById("span-monitor-threshold-unit"),
            divMonitorFieldPath: document.getElementById("div-monitor-field-path"),
            inputMonitorPath: document.getElementById("input-monitor-path"),
            divMonitorFieldInterface: document.getElementById("div-monitor-field-interface"),
            inputMonitorInterface: document.getElementById("input-monitor-interface"),
            divMonitorFieldPattern: document.getElementById("div-monitor-field-pattern"),
            inputMonitorPattern: document.getElementById("input-monitor-pattern"),
            checkboxMonitorRecursive: document.getElementById("checkbox-monitor-recursive"),
            divMonitorFieldMatch: document.getElementById("div-monitor-field-match"),
            inputMonitorMatch: document.getElementById("input-monitor-match"),
            selectMonitorSeverity: document.getElementById("select-monitor-severity"),
            pMonitorDuplicate: document.getElementById("p-monitor-duplicate"),
            buttonMonitorSubmit: document.getElementById("button-monitor-submit"),
            buttonMonitorCancel: document.getElementById("button-monitor-cancel"),
            pMonitorFormErrors: document.getElementById("p-monitor-form-errors"),
        };
    }

    #loadSelectOptions() {
        this.#elements.selectMonitorType.innerHTML = MONITOR_TYPES
            .map((monitorType) => `<option value="${monitorType}">${this.#htmler.getEscaped(MONITOR_TYPE_LABELS[monitorType])}</option>`).join("");
        this.#elements.selectMonitorSeverity.innerHTML = `<option value="">la que deduzca el backend</option>`
            + MONITOR_SEVERITIES.map((severity) => `<option value="${severity}">${severity}</option>`).join("");
    }

    #showDialog() {
        // <dialog>: modal nativo (foco atrapado, Esc cierra); sin soporte, al menos visible
        if (typeof this.#elements.dialogMonitor.showModal === "function") {
            if (!this.#elements.dialogMonitor.open) this.#elements.dialogMonitor.showModal();
            return;
        }
        this.#elements.dialogMonitor.setAttribute("open", "");
    }

    #renderScope(agentUuid, agentName) {
        this.#elements.hiddenMonitorAgentUuid.value = agentUuid;
        this.#elements.strongMonitorScope.textContent = agentUuid ? (agentName || agentUuid) : MonitorsTextEnum.ALL_AGENTS;
    }

    /**
     * Solo los campos que usa el tipo elegido: intervalo y umbrales en los de metricas, y la config propia de cada uno.
     */
    #renderFormByType() {
        const monitorType = this.#elements.selectMonitorType.value;
        const configFields = CONFIG_FIELDS_BY_TYPE[monitorType] ?? [];
        this.#elements.divMonitorMetricFields.classList.toggle(CSS_HIDDEN, !METRIC_MONITOR_TYPES.includes(monitorType));
        this.#elements.divMonitorFieldPath.classList.toggle(CSS_HIDDEN, !configFields.includes("path"));
        this.#elements.divMonitorFieldInterface.classList.toggle(CSS_HIDDEN, !configFields.includes("interface"));
        this.#elements.divMonitorFieldPattern.classList.toggle(CSS_HIDDEN, !configFields.includes("pattern"));
        this.#elements.divMonitorFieldMatch.classList.toggle(CSS_HIDDEN, !configFields.includes("match"));
        this.#elements.spanMonitorThresholdUnit.textContent = this.#getUnit(monitorType);
        this.#elements.inputMonitorPath.placeholder = monitorType === MonitorTypeEnum.DISK ? MonitorDefaultEnum.DISK_PATH : LINES_PATH_PLACEHOLDER;
        this.#renderDuplicate();
    }

    /**
     * Si lo que hay en el modal chocaria con un monitor ya creado: aviso con cual y por que, y el boton desactivado.
     */
    #renderDuplicate() {
        const formPrimitives = this.#getFormPrimitives();
        const duplicate = this.#monitorDuplicateFinder.getDuplicate([...this.#monitorsByUuid.values()], {
            monitorUuid: formPrimitives.monitorUuid,
            monitorType: formPrimitives.monitorType,
            monitorName: formPrimitives.monitorName,
            agentUuid: formPrimitives.agentUuid,
            config: {
                path: formPrimitives.path,
                interface: formPrimitives.interface,
                pattern: formPrimitives.pattern,
                match: formPrimitives.match,
            },
        });
        this.#hasDuplicate = duplicate !== null;
        this.#renderText(this.#elements.pMonitorDuplicate, duplicate ? this.#getDuplicateText(duplicate) : "");
        this.#renderSubmitState();
    }

    #getDuplicateText(duplicate) {
        const monitor = duplicate.monitor;
        if (duplicate.reason === MonitorDuplicateReasonEnum.NAME) {
            return `Repetido: ya hay un monitor que se llama ${monitor.monitor_name} en este proyecto. Elige otro nombre.`;
        }
        const scope = monitor.agent_uuid ? `el agente ${monitor.agent_name || monitor.agent_uuid}` : MonitorsTextEnum.ALL_AGENTS;
        const configLabel = this.#getConfigLabel(monitor);
        return `Repetido: ${monitor.monitor_name} (${monitor.monitor_type}${configLabel ? ` ${configLabel}` : ""}) ya vigila lo mismo en ${scope}.`
            + " Edita ese monitor en vez de crear otro.";
    }

    #renderSubmitState() {
        this.#elements.buttonMonitorSubmit.disabled = this.#isBusy || this.#hasDuplicate;
    }

    #renderTable() {
        const monitors = [...this.#monitorsByUuid.values()];
        if (monitors.length === 0) {
            this.#elements.divMonitorsTable.innerHTML = `<p class="empty">Sin monitores. Créalo desde la fila de un agente («Crear monitor») o para todos con el botón de arriba: los agentes lo aplican solos en menos de un minuto.</p>`;
            return;
        }
        this.#elements.divMonitorsTable.innerHTML = `<div class="table-wrap"><table class="table">
            <thead><tr>
                <th class="num">#</th><th>Tipo</th><th>Nombre y configuración</th><th>monitor_uuid</th><th>Ámbito</th>
                <th class="num">Cada</th><th>Umbrales</th><th>Activo</th><th>Último valor por agente</th><th>Acciones</th>
            </tr></thead>
            <tbody>${monitors.map((monitor, index) => this.#getRowHtml(monitor, index + 1)).join("")}</tbody>
        </table></div>`;
    }

    #getRowHtml(monitor, position) {
        const h = this.#htmler;
        const isMetric = METRIC_MONITOR_TYPES.includes(monitor.monitor_type);
        return `<tr class="${monitor.is_enabled ? "" : "is-disabled"}">
            <td class="num">${position}</td>
            <td class="mono">${h.getEscaped(monitor.monitor_type)}</td>
            <td><strong>${h.getEscaped(monitor.monitor_name)}</strong><br><span class="muted mono">${h.getEscaped(this.#getConfigLabel(monitor))}</span></td>
            <td class="nowrap">${h.getCopyableValue(monitor.monitor_uuid)}</td>
            <td>${h.getEscaped(monitor.agent_uuid ? (monitor.agent_name || monitor.agent_uuid) : MonitorsTextEnum.ALL_AGENTS)}</td>
            <td class="num nowrap">${isMetric ? `${h.getEscaped(monitor.interval_seconds)} s` : "cada pasada"}</td>
            <td class="nowrap">${h.getEscaped(this.#getThresholdsLabel(monitor))}</td>
            <td>${h.getBadge(monitor.is_enabled ? "enabled" : "disabled")}</td>
            <td>${this.#getStatesHtml(monitor.states ?? [])}</td>
            <td class="nowrap">${this.#getActionsHtml(monitor)}</td>
        </tr>`;
    }

    #getConfigLabel(monitor) {
        const config = monitor.config ?? {};
        switch (monitor.monitor_type) {
            case MonitorTypeEnum.DISK:
                return config.path || MonitorDefaultEnum.DISK_PATH;
            case MonitorTypeEnum.NETWORK:
                return config.interface ? `interfaz ${config.interface}` : "todas las interfaces";
            case MonitorTypeEnum.FOLDER:
                return `${config.path ?? ""}/${config.pattern || MonitorDefaultEnum.FILE_PATTERN}${config.is_recursive ? " (recursivo)" : ""}`;
            case MonitorTypeEnum.LOG_PATTERN:
                return `${config.path ?? ""}/${config.pattern || MonitorDefaultEnum.FILE_PATTERN}${config.is_recursive ? " (recursivo)" : ""}`
                    + ` ~ /${config.match ?? ""}/${config.severity ? ` → ${config.severity}` : ""}`;
            default:
                return "";
        }
    }

    #getThresholdsLabel(monitor) {
        if (!METRIC_MONITOR_TYPES.includes(monitor.monitor_type)) return "—";
        const unit = this.#getUnit(monitor.monitor_type);
        const hasWarning = monitor.warning_threshold !== null && monitor.warning_threshold !== undefined;
        const hasCritical = monitor.critical_threshold !== null && monitor.critical_threshold !== undefined;
        if (!hasWarning && !hasCritical) return "sin umbral";
        return [
            hasWarning ? `warning ≥ ${monitor.warning_threshold} ${unit}` : "",
            hasCritical ? `critical ≥ ${monitor.critical_threshold} ${unit}` : "",
        ].filter(Boolean).join(" · ");
    }

    #getUnit(monitorType) {
        return PERCENT_MONITOR_TYPES.includes(monitorType) ? UNIT_PERCENT : UNIT_MBIT;
    }

    #getStatesHtml(states) {
        const h = this.#htmler;
        if (states.length === 0) {
            return `<span class="muted">aún sin datos</span>`;
        }
        return states.map((state) => `<div class="monitor-state">
            <span class="mono">${h.getEscaped(state.agent_name)}</span> ${state.last_severity ? h.getBadge(state.last_severity) : ""}
            <span class="muted mono">${h.getEscaped(this.#dateTimer.getUtcDateTime(state.last_event_at))}</span><br>
            <span class="message">${h.getEscaped(state.last_message)}</span>
        </div>`).join("");
    }

    #getActionsHtml(monitor) {
        const monitorUuid = this.#htmler.getEscaped(monitor.monitor_uuid);
        if (monitor.monitor_uuid === this.#pendingDeleteMonitorUuid) {
            return `<span class="confirm">¿Borrar?
                <button type="button" class="button button-small button-danger" data-action="${ACTION_CONFIRM_DELETE}" data-monitor-uuid="${monitorUuid}">Sí, borrar</button>
                <button type="button" class="button button-small" data-action="${ACTION_CANCEL_DELETE}">Cancelar</button></span>`;
        }
        return `<button type="button" class="button button-small" data-action="${ACTION_EDIT}" data-monitor-uuid="${monitorUuid}">Editar</button>
            <button type="button" class="button button-small" data-action="${ACTION_TOGGLE}" data-monitor-uuid="${monitorUuid}">${monitor.is_enabled ? "Desactivar" : "Activar"}</button>
            <button type="button" class="button button-small button-ghost" data-action="${ACTION_ASK_DELETE}" data-monitor-uuid="${monitorUuid}">Borrar</button>`;
    }

    #onTableClick(event) {
        const button = event.target.closest("button[data-action]");
        if (!button) return;

        const monitor = this.#monitorsByUuid.get(button.dataset.monitorUuid ?? "");
        switch (button.dataset.action) {
            case ACTION_EDIT:
                if (monitor) this.#openEditDialog(monitor);
                return;
            case ACTION_TOGGLE:
                if (!monitor) return;
                this.#dispatcher.dispatch(MonitorsEventEnum.UPDATE_MONITOR_REQUESTED, {
                    ...this.#getFormPrimitivesOfMonitor(monitor),
                    isEnabled: !monitor.is_enabled,
                    isFromTable: true,
                });
                return;
            case ACTION_CONFIRM_DELETE:
                this.#pendingDeleteMonitorUuid = "";
                this.#dispatcher.dispatch(MonitorsEventEnum.DELETE_MONITOR_REQUESTED, { monitorUuid: button.dataset.monitorUuid });
                return;
            case ACTION_ASK_DELETE:
                this.#pendingDeleteMonitorUuid = button.dataset.monitorUuid;
                break;
            default:
                this.#pendingDeleteMonitorUuid = "";
        }
        this.#renderTable();
    }

    /**
     * Un monitor de la API como los primitivos del formulario (los que entiende CreateMonitorDto / UpdateMonitorDto).
     */
    #getFormPrimitivesOfMonitor(monitor) {
        const config = monitor.config ?? {};
        return {
            monitorUuid: monitor.monitor_uuid,
            monitorType: monitor.monitor_type,
            monitorName: monitor.monitor_name,
            agentUuid: monitor.agent_uuid ?? "",
            isEnabled: monitor.is_enabled === true,
            intervalSeconds: String(monitor.interval_seconds ?? ""),
            warningThreshold: String(monitor.warning_threshold ?? ""),
            criticalThreshold: String(monitor.critical_threshold ?? ""),
            path: config.path ?? "",
            interface: config.interface ?? "",
            pattern: config.pattern ?? "",
            isRecursive: config.is_recursive === true,
            match: config.match ?? "",
            severity: config.severity ?? "",
        };
    }

    #openEditDialog(monitor) {
        const primitives = this.#getFormPrimitivesOfMonitor(monitor);
        this.resetForm();
        this.#elements.hiddenMonitorUuid.value = primitives.monitorUuid;
        this.#elements.selectMonitorType.value = primitives.monitorType;
        // el tipo y el ambito de un monitor no cambian: se borra y se crea otro
        this.#elements.selectMonitorType.disabled = true;
        this.#renderScope(primitives.agentUuid, monitor.agent_name ?? "");
        this.#elements.inputMonitorName.value = primitives.monitorName;
        this.#elements.checkboxMonitorEnabled.checked = primitives.isEnabled;
        this.#elements.inputMonitorInterval.value = primitives.intervalSeconds;
        this.#elements.inputMonitorWarning.value = primitives.warningThreshold;
        this.#elements.inputMonitorCritical.value = primitives.criticalThreshold;
        this.#elements.inputMonitorPath.value = primitives.path;
        this.#elements.inputMonitorInterface.value = primitives.interface;
        this.#elements.inputMonitorPattern.value = primitives.pattern;
        this.#elements.checkboxMonitorRecursive.checked = primitives.isRecursive;
        this.#elements.inputMonitorMatch.value = primitives.match;
        this.#elements.selectMonitorSeverity.value = primitives.severity;
        this.#elements.h2MonitorForm.textContent = this.#textFormatter.getFormatted(MonitorsTextEnum.FORM_TITLE_UPDATE, primitives.monitorName);
        this.#elements.buttonMonitorSubmit.textContent = MonitorsTextEnum.SUBMIT_UPDATE;
        this.#renderFormByType();
        this.renderMessage("");
        this.#showDialog();
        this.#elements.inputMonitorName.focus();
    }

    /**
     * Lo que hay ahora en el modal, como los primitivos de CreateMonitorDto / UpdateMonitorDto.
     */
    #getFormPrimitives() {
        return {
            projectUuid: this.#projectUuid,
            monitorUuid: this.#elements.hiddenMonitorUuid.value,
            monitorType: this.#elements.selectMonitorType.value,
            monitorName: this.#elements.inputMonitorName.value,
            agentUuid: this.#elements.hiddenMonitorAgentUuid.value,
            isEnabled: this.#elements.checkboxMonitorEnabled.checked,
            intervalSeconds: this.#elements.inputMonitorInterval.value,
            warningThreshold: this.#elements.inputMonitorWarning.value,
            criticalThreshold: this.#elements.inputMonitorCritical.value,
            path: this.#elements.inputMonitorPath.value,
            interface: this.#elements.inputMonitorInterface.value,
            pattern: this.#elements.inputMonitorPattern.value,
            isRecursive: this.#elements.checkboxMonitorRecursive.checked,
            match: this.#elements.inputMonitorMatch.value,
            severity: this.#elements.selectMonitorSeverity.value,
        };
    }

    #onSubmit(event) {
        event.preventDefault();
        this.renderFormErrors([]);
        this.#renderDuplicate();
        if (this.#hasDuplicate) return;

        const primitives = this.#getFormPrimitives();
        this.#dispatcher.dispatch(
            primitives.monitorUuid ? MonitorsEventEnum.UPDATE_MONITOR_REQUESTED : MonitorsEventEnum.CREATE_MONITOR_REQUESTED,
            primitives
        );
    }

    #renderText(element, text) {
        element.textContent = text;
        element.classList.toggle(CSS_HIDDEN, text === "");
    }
}
