import { HttpCodeEnum } from "/js/Modules/Shared/Domain/Enums/HttpCodeEnum.js";
import { ScreenEnum } from "/js/Modules/Shared/Domain/Enums/ScreenEnum.js";
import { SharedEventEnum } from "/js/Modules/Shared/Domain/Enums/SharedEventEnum.js";
import { Dispatcher } from "/js/Modules/Shared/Infrastructure/Components/Dispatcher.js";
import { SessionStorer } from "/js/Modules/Shared/Infrastructure/Components/SessionStorer.js";
import { Subscriber } from "/js/Modules/Shared/Infrastructure/Components/Subscriber.js";

import { SessionEventEnum } from "/js/Modules/Session/Domain/Enums/SessionEventEnum.js";

import { MonitorsMessageEnum } from "/js/Modules/Monitors/Domain/Enums/MonitorsMessageEnum.js";
import { ListMonitorsDto } from "/js/Modules/Monitors/Application/ListMonitors/ListMonitorsDto.js";
import { ListMonitorsService } from "/js/Modules/Monitors/Application/ListMonitors/ListMonitorsService.js";
import { MonitorsView } from "/js/Modules/Monitors/Infrastructure/Views/MonitorsView.js";
import { MonitorsException } from "/js/Modules/Monitors/Domain/Exceptions/MonitorsException.js";

/**
 * Los monitores del proyecto que se esta viendo (#/projects/<project_uuid>): se piden cada vez que se pide esa
 * pantalla (tambien en cada refresco: el ultimo valor por agente cambia solo). El refresco lo programa el controller
 * del detalle, no este.
 */
export class ListMonitorsController {

    #dispatcher;
    #sessionStorer;
    #subscriber;
    #monitorsView;

    constructor() {
        this.#dispatcher = Dispatcher.getInstance();
        this.#sessionStorer = SessionStorer.getInstance();
        this.#subscriber = Subscriber.getInstance();
        this.#monitorsView = MonitorsView.getInstance();
    }

    /**
     * @returns {ListMonitorsController}
     */
    static getInstance() {
        return new ListMonitorsController();
    }

    init() {
        this.#subscriber.subscribe(SharedEventEnum.SCREEN_REQUESTED, (request) => {
            if (request.screen !== ScreenEnum.PROJECT_DETAIL) return;
            this.#monitorsView.renderProject(request.pathParam ?? "");
            this.invoke({ projectUuid: request.pathParam });
        });
        this.#subscriber.subscribe(SessionEventEnum.USER_LOGGED_OUT, () => this.#monitorsView.renderProject(""));
    }

    /**
     * @param {{projectUuid: string}} request
     */
    async invoke(request) {
        try {
            const listMonitorsResultDto = await ListMonitorsService.getInstance().invoke(
                ListMonitorsDto.fromPrimitives({
                    dashboardToken: this.#sessionStorer.getDashboardToken(),
                    projectUuid: request.projectUuid,
                })
            );
            this.#monitorsView.renderMonitors(listMonitorsResultDto.getMonitors());
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
            console.error("ListMonitorsController.invoke", error);
            this.#monitorsView.renderErrors([MonitorsMessageEnum.UNEXPECTED_ERROR]);
        }
    }
}
