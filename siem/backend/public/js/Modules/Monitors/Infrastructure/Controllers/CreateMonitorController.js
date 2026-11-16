import { HttpCodeEnum } from "/js/Modules/Shared/Domain/Enums/HttpCodeEnum.js";
import { Dispatcher } from "/js/Modules/Shared/Infrastructure/Components/Dispatcher.js";
import { SessionStorer } from "/js/Modules/Shared/Infrastructure/Components/SessionStorer.js";
import { Subscriber } from "/js/Modules/Shared/Infrastructure/Components/Subscriber.js";
import { TextFormatter } from "/js/Modules/Shared/Infrastructure/Components/TextFormatter.js";
import { NavigationView } from "/js/Modules/Shared/Infrastructure/Views/NavigationView.js";

import { SessionEventEnum } from "/js/Modules/Session/Domain/Enums/SessionEventEnum.js";

import { AgentsEventEnum } from "/js/Modules/Agents/Domain/Enums/AgentsEventEnum.js";

import { MonitorsEventEnum } from "/js/Modules/Monitors/Domain/Enums/MonitorsEventEnum.js";
import { MonitorsMessageEnum } from "/js/Modules/Monitors/Domain/Enums/MonitorsMessageEnum.js";
import { MonitorsTextEnum } from "/js/Modules/Monitors/Domain/Enums/MonitorsTextEnum.js";
import { CreateMonitorDto } from "/js/Modules/Monitors/Application/CreateMonitor/CreateMonitorDto.js";
import { CreateMonitorService } from "/js/Modules/Monitors/Application/CreateMonitor/CreateMonitorService.js";
import { MonitorsView } from "/js/Modules/Monitors/Infrastructure/Views/MonitorsView.js";
import { MonitorsException } from "/js/Modules/Monitors/Domain/Exceptions/MonitorsException.js";

/**
 * Alta de un monitor desde el modal del detalle de un proyecto: creado, se repinta la pantalla (sale en la tabla).
 * El modal lo abre "Crear monitor" en la fila de un agente (Agents avisa con CREATE_MONITOR_REQUESTED: el monitor
 * sera de ese agente) o el boton de "todos los agentes" de la propia vista.
 */
export class CreateMonitorController {

    #dispatcher;
    #sessionStorer;
    #subscriber;
    #textFormatter;
    #navigationView;
    #monitorsView;

    constructor() {
        this.#dispatcher = Dispatcher.getInstance();
        this.#sessionStorer = SessionStorer.getInstance();
        this.#subscriber = Subscriber.getInstance();
        this.#textFormatter = TextFormatter.getInstance();
        this.#navigationView = NavigationView.getInstance();
        this.#monitorsView = MonitorsView.getInstance();
    }

    /**
     * @returns {CreateMonitorController}
     */
    static getInstance() {
        return new CreateMonitorController();
    }

    init() {
        this.#subscriber.subscribe(MonitorsEventEnum.CREATE_MONITOR_REQUESTED, (request) => this.invoke(request));
        this.#subscriber.subscribe(AgentsEventEnum.CREATE_MONITOR_REQUESTED, (request) => {
            this.#monitorsView.openCreateDialog(request.agentUuid ?? "", request.agentName ?? "");
        });
    }

    /**
     * @param {Object} request los campos del formulario (ver CreateMonitorDto) + projectUuid
     */
    async invoke(request) {
        this.#monitorsView.renderBusy(true);
        try {
            const createMonitorResultDto = await CreateMonitorService.getInstance().invoke(
                CreateMonitorDto.fromPrimitives({ ...request, dashboardToken: this.#sessionStorer.getDashboardToken() })
            );
            this.#monitorsView.renderSaved(
                this.#textFormatter.getFormatted(MonitorsTextEnum.MONITOR_CREATED, createMonitorResultDto.getMonitorName())
            );
            this.#navigationView.renderCurrentScreen();
        }
        catch (error) {
            if (error instanceof MonitorsException && error.getCode() === HttpCodeEnum.UNAUTHORIZED) {
                this.#dispatcher.dispatch(SessionEventEnum.SESSION_EXPIRED);
                return;
            }
            if (error instanceof MonitorsException) {
                this.#monitorsView.renderFormErrors(error.getErrors());
                return;
            }
            console.error("CreateMonitorController.invoke", error);
            this.#monitorsView.renderFormErrors([MonitorsMessageEnum.UNEXPECTED_ERROR]);
        }
        finally {
            this.#monitorsView.renderBusy(false);
        }
    }
}
