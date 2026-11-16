import { HttpCodeEnum } from "/js/Modules/Shared/Domain/Enums/HttpCodeEnum.js";
import { ScreenEnum } from "/js/Modules/Shared/Domain/Enums/ScreenEnum.js";
import { SharedEventEnum } from "/js/Modules/Shared/Domain/Enums/SharedEventEnum.js";
import { Dispatcher } from "/js/Modules/Shared/Infrastructure/Components/Dispatcher.js";
import { SessionStorer } from "/js/Modules/Shared/Infrastructure/Components/SessionStorer.js";
import { Subscriber } from "/js/Modules/Shared/Infrastructure/Components/Subscriber.js";

import { SessionEventEnum } from "/js/Modules/Session/Domain/Enums/SessionEventEnum.js";

import { GetDashboardDto } from "/js/Modules/Dashboard/Application/GetDashboard/GetDashboardDto.js";
import { GetDashboardService } from "/js/Modules/Dashboard/Application/GetDashboard/GetDashboardService.js";

import { DownloadsMessageEnum } from "/js/Modules/Downloads/Domain/Enums/DownloadsMessageEnum.js";
import { GetDownloadsDto } from "/js/Modules/Downloads/Application/GetDownloads/GetDownloadsDto.js";
import { GetDownloadsService } from "/js/Modules/Downloads/Application/GetDownloads/GetDownloadsService.js";
import { DownloadsView } from "/js/Modules/Downloads/Infrastructure/Views/DownloadsView.js";
import { DownloadsException } from "/js/Modules/Downloads/Domain/Exceptions/DownloadsException.js";

/**
 * Descargas (#/downloads): asistente de instalacion (binarios, instaladores y proyectos). No se repinta sola: solo
 * cambia con un `make dist`, y un repintado borraria lo que se este escribiendo en el asistente.
 */
export class GetDownloadsController {

    #dispatcher;
    #sessionStorer;
    #subscriber;
    #downloadsView;

    constructor() {
        this.#dispatcher = Dispatcher.getInstance();
        this.#sessionStorer = SessionStorer.getInstance();
        this.#subscriber = Subscriber.getInstance();
        this.#downloadsView = DownloadsView.getInstance();
    }

    /**
     * @returns {GetDownloadsController}
     */
    static getInstance() {
        return new GetDownloadsController();
    }

    init() {
        this.#subscriber.subscribe(SharedEventEnum.SCREEN_REQUESTED, (request) => {
            if (request.screen !== ScreenEnum.DOWNLOADS) return;
            this.invoke();
        });
        // el token del agente que se haya pegado para montar los comandos no sobrevive a la sesion
        this.#subscriber.subscribe(SessionEventEnum.USER_LOGGED_OUT, () => this.#downloadsView.reset());
    }

    async invoke() {
        try {
            const getDownloadsResultDto = await GetDownloadsService.getInstance().invoke(
                GetDownloadsDto.fromPrimitives({
                    dashboardToken: this.#sessionStorer.getDashboardToken(),
                })
            );
            this.#downloadsView.renderDownloads(
                getDownloadsResultDto.getEndpointUrl(), getDownloadsResultDto.getDownloads(), getDownloadsResultDto.getInstallers()
            );
            await this.#renderProjects();
        }
        catch (error) {
            if (error instanceof DownloadsException && error.getCode() === HttpCodeEnum.UNAUTHORIZED) {
                this.#dispatcher.dispatch(SessionEventEnum.SESSION_EXPIRED);
                return;
            }
            if (error instanceof DownloadsException) {
                this.#downloadsView.renderErrors(error.getErrors());
            }
            else {
                console.error("GetDownloadsController.invoke", error);
                this.#downloadsView.renderErrors([DownloadsMessageEnum.UNEXPECTED_ERROR]);
            }
        }
        this.#dispatcher.dispatch(SharedEventEnum.SCREEN_RENDERED, { isAutoRefreshed: false });
    }

    /**
     * Los proyectos del paso 2 (crear el agente en uno de ellos). Si no se pueden leer, el desplegable queda vacio y el
     * resto del asistente sigue sirviendo (con un token pegado).
     */
    async #renderProjects() {
        try {
            const getDashboardResultDto = await GetDashboardService.getInstance().invoke(
                GetDashboardDto.fromPrimitives({ dashboardToken: this.#sessionStorer.getDashboardToken() })
            );
            this.#downloadsView.renderProjects(getDashboardResultDto.getProjects());
        }
        catch (error) {
            console.error("GetDownloadsController.renderProjects", error);
            this.#downloadsView.renderProjects([]);
        }
    }
}
