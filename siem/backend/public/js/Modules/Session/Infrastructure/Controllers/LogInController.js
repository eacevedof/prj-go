import { HttpCodeEnum } from "/js/Modules/Shared/Domain/Enums/HttpCodeEnum.js";
import { Dispatcher } from "/js/Modules/Shared/Infrastructure/Components/Dispatcher.js";
import { Subscriber } from "/js/Modules/Shared/Infrastructure/Components/Subscriber.js";

import { SessionEventEnum } from "/js/Modules/Session/Domain/Enums/SessionEventEnum.js";
import { SessionMessageEnum } from "/js/Modules/Session/Domain/Enums/SessionMessageEnum.js";
import { LogInDto } from "/js/Modules/Session/Application/LogIn/LogInDto.js";
import { LogInService } from "/js/Modules/Session/Application/LogIn/LogInService.js";
import { AuthView } from "/js/Modules/Session/Infrastructure/Views/AuthView.js";
import { SessionException } from "/js/Modules/Session/Domain/Exceptions/SessionException.js";

/**
 * Entrar con el token del dashboard (formulario o el guardado al recargar): si la API lo acepta se entra
 * (SESSION_STARTED); un 401 es sesion caducada (borra el guardado y vuelve a la entrada con aviso).
 */
export class LogInController {

    #dispatcher;
    #subscriber;
    #authView;

    constructor() {
        this.#dispatcher = Dispatcher.getInstance();
        this.#subscriber = Subscriber.getInstance();
        this.#authView = AuthView.getInstance();
    }

    /**
     * @returns {LogInController}
     */
    static getInstance() {
        return new LogInController();
    }

    init() {
        this.#subscriber.subscribe(SessionEventEnum.LOG_IN_REQUESTED, (request) => this.invoke(request));
    }

    /**
     * @param {{dashboardToken: string}} request
     */
    async invoke(request) {
        try {
            const logInResultDto = await LogInService.getInstance().invoke(
                LogInDto.fromPrimitives({
                    dashboardToken: request.dashboardToken,
                })
            );
            this.#dispatcher.dispatch(SessionEventEnum.SESSION_STARTED, {
                dashboardToken: logInResultDto.getDashboardToken(),
            });
        }
        catch (error) {
            if (error instanceof SessionException && error.getCode() === HttpCodeEnum.UNAUTHORIZED) {
                this.#dispatcher.dispatch(SessionEventEnum.SESSION_EXPIRED);
                return;
            }
            if (error instanceof SessionException) {
                this.#authView.renderErrors(error.getErrors());
                return;
            }
            console.error("LogInController.invoke", error);
            this.#authView.show();
            this.#authView.renderErrors([SessionMessageEnum.UNEXPECTED_ERROR]);
        }
    }
}
