import { HttpCodeEnum } from "/js/Modules/Shared/Domain/Enums/HttpCodeEnum.js";
import { ScreenEnum } from "/js/Modules/Shared/Domain/Enums/ScreenEnum.js";
import { SharedEventEnum } from "/js/Modules/Shared/Domain/Enums/SharedEventEnum.js";
import { Dispatcher } from "/js/Modules/Shared/Infrastructure/Components/Dispatcher.js";
import { SessionStorer } from "/js/Modules/Shared/Infrastructure/Components/SessionStorer.js";
import { Subscriber } from "/js/Modules/Shared/Infrastructure/Components/Subscriber.js";

import { SessionEventEnum } from "/js/Modules/Session/Domain/Enums/SessionEventEnum.js";

import { DownloadsMessageEnum } from "/js/Modules/Downloads/Domain/Enums/DownloadsMessageEnum.js";
import { GetDownloadsDto } from "/js/Modules/Downloads/Application/GetDownloads/GetDownloadsDto.js";
import { GetDownloadsService } from "/js/Modules/Downloads/Application/GetDownloads/GetDownloadsService.js";
import { DownloadsView } from "/js/Modules/Downloads/Infrastructure/Views/DownloadsView.js";
import { DownloadsException } from "/js/Modules/Downloads/Domain/Exceptions/DownloadsException.js";

/**
 * Descargas (#/downloads): binarios del agente e instrucciones de instalacion. No se repinta sola: los binarios solo
 * cambian con un `make dist`, y un repintado cerraria las instrucciones que se esten leyendo.
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
            this.#downloadsView.renderDownloads(getDownloadsResultDto.getEndpointUrl(), getDownloadsResultDto.getDownloads());
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
}
