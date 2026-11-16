import { HttpCodeEnum } from "/js/Modules/Shared/Domain/Enums/HttpCodeEnum.js";
import { Dispatcher } from "/js/Modules/Shared/Infrastructure/Components/Dispatcher.js";
import { SessionStorer } from "/js/Modules/Shared/Infrastructure/Components/SessionStorer.js";
import { Subscriber } from "/js/Modules/Shared/Infrastructure/Components/Subscriber.js";
import { TextFormatter } from "/js/Modules/Shared/Infrastructure/Components/TextFormatter.js";
import { NavigationView } from "/js/Modules/Shared/Infrastructure/Views/NavigationView.js";

import { SessionEventEnum } from "/js/Modules/Session/Domain/Enums/SessionEventEnum.js";

import { MonitorsEventEnum } from "/js/Modules/Monitors/Domain/Enums/MonitorsEventEnum.js";
import { MonitorsMessageEnum } from "/js/Modules/Monitors/Domain/Enums/MonitorsMessageEnum.js";
import { MonitorsTextEnum } from "/js/Modules/Monitors/Domain/Enums/MonitorsTextEnum.js";
import { UpdateMonitorDto } from "/js/Modules/Monitors/Application/UpdateMonitor/UpdateMonitorDto.js";
import { UpdateMonitorService } from "/js/Modules/Monitors/Application/UpdateMonitor/UpdateMonitorService.js";
import { MonitorsView } from "/js/Modules/Monitors/Infrastructure/Views/MonitorsView.js";
import { MonitorsException } from "/js/Modules/Monitors/Domain/Exceptions/MonitorsException.js";

/**
 * Cambio de un monitor: desde el formulario (editar) o desde la tabla (activar / desactivar, isFromTable). Los errores
 * se pintan donde se pidio el cambio: en el formulario o encima de la tabla.
 */
export class UpdateMonitorController {

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
     * @returns {UpdateMonitorController}
     */
    static getInstance() {
        return new UpdateMonitorController();
    }

    init() {
        this.#subscriber.subscribe(MonitorsEventEnum.UPDATE_MONITOR_REQUESTED, (request) => this.invoke(request));
    }

    /**
     * @param {Object} request los campos del formulario (ver UpdateMonitorDto) + monitorUuid + isFromTable?
     */
    async invoke(request) {
        const isFromTable = request.isFromTable === true;
        if (!isFromTable) this.#monitorsView.renderBusy(true);
        try {
            const updateMonitorResultDto = await UpdateMonitorService.getInstance().invoke(
                UpdateMonitorDto.fromPrimitives({ ...request, dashboardToken: this.#sessionStorer.getDashboardToken() })
            );
            if (!isFromTable) {
                this.#monitorsView.renderSaved(
                    this.#textFormatter.getFormatted(MonitorsTextEnum.MONITOR_UPDATED, updateMonitorResultDto.getMonitorName())
                );
            }
            this.#navigationView.renderCurrentScreen();
        }
        catch (error) {
            if (error instanceof MonitorsException && error.getCode() === HttpCodeEnum.UNAUTHORIZED) {
                this.#dispatcher.dispatch(SessionEventEnum.SESSION_EXPIRED);
                return;
            }
            const errors = error instanceof MonitorsException ? error.getErrors() : [MonitorsMessageEnum.UNEXPECTED_ERROR];
            if (!(error instanceof MonitorsException)) {
                console.error("UpdateMonitorController.invoke", error);
            }
            if (isFromTable) {
                this.#monitorsView.renderErrors(errors);
                return;
            }
            this.#monitorsView.renderFormErrors(errors);
        }
        finally {
            if (!isFromTable) this.#monitorsView.renderBusy(false);
        }
    }
}
