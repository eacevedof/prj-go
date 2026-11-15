import { DateTimer } from "/js/Modules/Shared/Infrastructure/Components/DateTimer.js";
import { Dispatcher } from "/js/Modules/Shared/Infrastructure/Components/Dispatcher.js";
import { Htmler } from "/js/Modules/Shared/Infrastructure/Components/Htmler.js";

import { MonitorsEventEnum } from "/js/Modules/Monitors/Domain/Enums/MonitorsEventEnum.js";
import {
    METRIC_MONITOR_TYPES, MONITOR_SEVERITIES, MONITOR_TYPE_LABELS, MONITOR_TYPES, MonitorDefaultEnum, MonitorTypeEnum,
    PERCENT_MONITOR_TYPES,
} from "/js/Modules/Monitors/Domain/Enums/MonitorTypeEnum.js";

const CSS_HIDDEN = "hidden";
const ACTION_EDIT = "edit";
const ACTION_TOGGLE = "toggle";
const ACTION_ASK_DELETE = "ask-delete";
const ACTION_CONFIRM_DELETE = "confirm-delete";
const ACTION_CANCEL_DELETE = "cancel-delete";
const ALL_AGENTS_LABEL = "todos los agentes";
const UNIT_PERCENT = "%";
const UNIT_MBIT = "Mbit/s";
const TITLE_CREATE = "5 - Nuevo monitor";
const SUBMIT_CREATE = "Crear monitor";
const SUBMIT_UPDATE = "Guardar cambios";
// que campos de config enseña el formulario para cada tipo
const CONFIG_FIELDS_BY_TYPE = Object.freeze({
    [MonitorTypeEnum.DISK]: ["path"],
    [MonitorTypeEnum.NETWORK]: ["interface"],
    [MonitorTypeEnum.FOLDER]: ["path", "pattern"],
    [MonitorTypeEnum.LOG_PATTERN]: ["path", "pattern", "match"],
});

/**
 * "Monitores" del detalle de un proyecto: los recursos a vigilar (cpu, memoria, disco, gpu, red, carpetas, patrones
 * de log). Tabla con su monitor_uuid (copiable), ambito, umbrales y ultimo valor por agente; y un formulario (el mismo
 * para alta y edicion) cuyo bloque de config cambia con el tipo. El formulario es html fijo: el refresco solo repinta
 * la tabla, asi que lo que se este escribiendo no se pierde. Borrar pide un segundo clic en la propia fila.
 * Una sola instancia: guarda las referencias al DOM, el proyecto y los monitores pintados.
 */
export class MonitorsView {

    static #instance = null;

    #dateTimer;
    #dispatcher;
    #htmler;
    #elements = {};
    #projectUuid = "";
    #monitorsByUuid = new Map();
    #pendingDeleteMonitorUuid = "";

    constructor() {
        this.#dateTimer = DateTimer.getInstance();
        this.#dispatcher = Dispatcher.getInstance();
        this.#htmler = Htmler.getInstance();
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
        this.#elements.formMonitor.addEventListener("submit", (event) => this.#onSubmit(event));
        this.#elements.buttonMonitorCancel.addEventListener("click", () => this.resetForm());
        this.#elements.divMonitorsTable.addEventListener("click", (event) => this.#onTableClick(event));
    }

    /**
     * El proyecto que se esta viendo. Al cambiar de proyecto, tabla y formulario se vacian; en un refresco, no.
     * @param {string} projectUuid
     */
    renderProject(projectUuid) {
        if (projectUuid === this.#projectUuid) return;
        this.#projectUuid = projectUuid;
        this.#pendingDeleteMonitorUuid = "";
        this.#monitorsByUuid = new Map();
        this.#elements.divMonitorsTable.innerHTML = "";
        this.renderAgentOptions([]);
        this.resetForm();
        this.renderErrors([]);
    }

    /**
     * A que agentes se puede asignar un monitor (ademas de "todos"). Respeta lo que hubiera elegido.
     * @param {{agent_uuid: string, agent_name: string}[]} agents los del proyecto
     */
    renderAgentOptions(agents) {
        const h = this.#htmler;
        const selectedAgentUuid = this.#elements.selectMonitorAgent.value;
        this.#elements.selectMonitorAgent.innerHTML = `<option value="">${ALL_AGENTS_LABEL}</option>`
            + agents.map((agent) => `<option value="${h.getEscaped(agent.agent_uuid)}">${h.getEscaped(agent.agent_name)}</option>`).join("");
        this.#elements.selectMonitorAgent.value = agents.some((agent) => agent.agent_uuid === selectedAgentUuid) ? selectedAgentUuid : "";
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
    }

    /**
     * @param {string[]} errors de la lista (no se pudo cargar, o fallo activar / borrar)
     */
    renderErrors(errors) {
        this.#renderText(this.#elements.pMonitorsErrors, errors.join(" · "));
    }

    /**
     * @param {string[]} errors del formulario
     */
    renderFormErrors(errors) {
        this.#renderText(this.#elements.pMonitorFormErrors, errors.join(" · "));
    }

    /**
     * Guardado: formulario a "nuevo monitor" y aviso.
     * @param {string} message
     */
    renderSaved(message) {
        this.resetForm();
        this.#renderText(this.#elements.pMonitorFormMessage, message);
    }

    /**
     * @param {boolean} isBusy
     */
    renderBusy(isBusy) {
        this.#elements.buttonMonitorSubmit.disabled = isBusy;
    }

    resetForm() {
        this.#elements.formMonitor.reset();
        this.#elements.hiddenMonitorUuid.value = "";
        this.#elements.selectMonitorType.disabled = false;
        this.#elements.inputMonitorInterval.value = String(MonitorDefaultEnum.INTERVAL_SECONDS);
        this.#elements.checkboxMonitorEnabled.checked = true;
        this.#elements.h2MonitorForm.textContent = TITLE_CREATE;
        this.#elements.buttonMonitorSubmit.textContent = SUBMIT_CREATE;
        this.#elements.buttonMonitorCancel.classList.add(CSS_HIDDEN);
        this.renderFormErrors([]);
        this.#renderText(this.#elements.pMonitorFormMessage, "");
        this.#renderFormByType();
    }

    #loadElements() {
        this.#elements = {
            pMonitorsErrors: document.getElementById("p-monitors-errors"),
            divMonitorsTable: document.getElementById("div-monitors-table"),
            h2MonitorForm: document.getElementById("h2-monitor-form"),
            formMonitor: document.getElementById("form-monitor"),
            hiddenMonitorUuid: document.getElementById("hidden-monitor-uuid"),
            selectMonitorType: document.getElementById("select-monitor-type"),
            inputMonitorName: document.getElementById("input-monitor-name"),
            selectMonitorAgent: document.getElementById("select-monitor-agent"),
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
            buttonMonitorSubmit: document.getElementById("button-monitor-submit"),
            buttonMonitorCancel: document.getElementById("button-monitor-cancel"),
            pMonitorFormErrors: document.getElementById("p-monitor-form-errors"),
            pMonitorFormMessage: document.getElementById("p-monitor-form-message"),
        };
    }

    #loadSelectOptions() {
        this.#elements.selectMonitorType.innerHTML = MONITOR_TYPES
            .map((monitorType) => `<option value="${monitorType}">${this.#htmler.getEscaped(MONITOR_TYPE_LABELS[monitorType])}</option>`).join("");
        this.#elements.selectMonitorSeverity.innerHTML = `<option value="">la que deduzca el backend</option>`
            + MONITOR_SEVERITIES.map((severity) => `<option value="${severity}">${severity}</option>`).join("");
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
        this.#elements.inputMonitorPath.placeholder = monitorType === MonitorTypeEnum.DISK ? MonitorDefaultEnum.DISK_PATH : "/var/log/nginx";
    }

    #renderTable() {
        const monitors = [...this.#monitorsByUuid.values()];
        if (monitors.length === 0) {
            this.#elements.divMonitorsTable.innerHTML = `<p class="empty">Sin monitores. Crea el primero en el apartado 5: los agentes del proyecto lo aplican solos en menos de un minuto.</p>`;
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
            <td>${h.getEscaped(monitor.agent_uuid ? (monitor.agent_name || monitor.agent_uuid) : ALL_AGENTS_LABEL)}</td>
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
                if (monitor) this.#renderFormForEdit(monitor);
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

    #renderFormForEdit(monitor) {
        const primitives = this.#getFormPrimitivesOfMonitor(monitor);
        this.resetForm();
        this.#elements.hiddenMonitorUuid.value = primitives.monitorUuid;
        this.#elements.selectMonitorType.value = primitives.monitorType;
        // el tipo de un monitor no cambia: se borra y se crea otro
        this.#elements.selectMonitorType.disabled = true;
        this.#elements.inputMonitorName.value = primitives.monitorName;
        this.#elements.selectMonitorAgent.value = primitives.agentUuid;
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
        this.#elements.h2MonitorForm.textContent = `5 - Editar monitor ${primitives.monitorName}`;
        this.#elements.buttonMonitorSubmit.textContent = SUBMIT_UPDATE;
        this.#elements.buttonMonitorCancel.classList.remove(CSS_HIDDEN);
        this.#renderFormByType();
        this.#elements.formMonitor.scrollIntoView({ block: "nearest" });
        this.#elements.inputMonitorName.focus();
    }

    #onSubmit(event) {
        event.preventDefault();
        this.renderFormErrors([]);
        this.#renderText(this.#elements.pMonitorFormMessage, "");

        const monitorUuid = this.#elements.hiddenMonitorUuid.value;
        const primitives = {
            projectUuid: this.#projectUuid,
            monitorUuid,
            monitorType: this.#elements.selectMonitorType.value,
            monitorName: this.#elements.inputMonitorName.value,
            agentUuid: this.#elements.selectMonitorAgent.value,
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
        this.#dispatcher.dispatch(
            monitorUuid ? MonitorsEventEnum.UPDATE_MONITOR_REQUESTED : MonitorsEventEnum.CREATE_MONITOR_REQUESTED,
            primitives
        );
    }

    #renderText(element, text) {
        element.textContent = text;
        element.classList.toggle(CSS_HIDDEN, text === "");
    }
}
