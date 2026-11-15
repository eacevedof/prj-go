import { HttpCodeEnum } from "/js/Modules/Shared/Domain/Enums/HttpCodeEnum.js";
import { Dispatcher } from "/js/Modules/Shared/Infrastructure/Components/Dispatcher.js";
import { SessionStorer } from "/js/Modules/Shared/Infrastructure/Components/SessionStorer.js";
import { Subscriber } from "/js/Modules/Shared/Infrastructure/Components/Subscriber.js";
import { NavigationView } from "/js/Modules/Shared/Infrastructure/Views/NavigationView.js";

import { SessionEventEnum } from "/js/Modules/Session/Domain/Enums/SessionEventEnum.js";

import { MonitorsEventEnum } from "/js/Modules/Monitors/Domain/Enums/MonitorsEventEnum.js";
import { MonitorsMessageEnum } from "/js/Modules/Monitors/Domain/Enums/MonitorsMessageEnum.js";
import { DeleteMonitorDto } from "/js/Modules/Monitors/Application/DeleteMonitor/DeleteMonitorDto.js";
import { DeleteMonitorService } from "/js/Modules/Monitors/Application/DeleteMonitor/DeleteMonitorService.js";
import { MonitorsView } from "/js/Modules/Monitors/Infrastructure/Views/MonitorsView.js";
import { MonitorsException } from "/js/Modules/Monitors/Domain/Exceptions/MonitorsException.js";

/**
 * Baja de un monitor (ya confirmada en su fila): borrado, se repinta la pantalla.
 */
export class DeleteMonitorController {

    #dispatcher;
    #sessionStorer;
    #subscriber;
    #navigationView;
    #monitorsView;

    constructor() {
        this.#dispatcher = Dispatcher.getInstance();
        this.#sessionStorer = SessionStorer.getInstance();
        this.#subscriber = Subscriber.getInstance();
        this.#navigationView = NavigationView.getInstance();
        this.#monitorsView = MonitorsView.getInstance();
    }

    /**
     * @returns {DeleteMonitorController}
     */
    static getInstance() {
        return new DeleteMonitorController();
    }

    init() {
        this.#subscriber.subscribe(MonitorsEventEnum.DELETE_MONITOR_REQUESTED, (request) => this.invoke(request));
    }

    /**
     * @param {{monitorUuid: string}} request
     */
    async invoke(request) {
        try {
            await DeleteMonitorService.getInstance().invoke(
                DeleteMonitorDto.fromPrimitives({
                    dashboardToken: this.#sessionStorer.getDashboardToken(),
                    monitorUuid: request.monitorUuid,
                })
            );
            // si se estaba editando ese monitor, el formulario no puede quedarse apuntando a uno que ya no existe
            this.#monitorsView.resetForm();
            this.#navigationView.renderCurrentScreen();
        }
        catch (error) {
            if (error instanceof MonitorsException && error.getCode() === HttpCodeEnum.UNAUTHORIZED) {
                this.#dispatcher.dispatch(SessionEventEnum.SESSION_EXPIRED);
                return;
            }
            if (error instanceof MonitorsException) {
                this.#monitorsView.renderErrors(error.getErrors());
                return;
            }
            console.error("DeleteMonitorController.invoke", error);
            this.#monitorsView.renderErrors([MonitorsMessageEnum.UNEXPECTED_ERROR]);
        }
    }
}
