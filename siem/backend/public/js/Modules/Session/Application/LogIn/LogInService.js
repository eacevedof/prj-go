import { SessionMessageEnum } from "/js/Modules/Session/Domain/Enums/SessionMessageEnum.js";
import { LogInResultDto } from "/js/Modules/Session/Application/LogIn/LogInResultDto.js";
import { SessionReaderHttpRepository } from "/js/Modules/Session/Infrastructure/Repositories/SessionReaderHttpRepository.js";
import { SessionException } from "/js/Modules/Session/Domain/Exceptions/SessionException.js";

/**
 * Entrar con el token del dashboard: no hay usuarios ni login en la API, asi que el token se comprueba con una
 * lectura autenticada; si la API lo acepta, vale.
 */
export class LogInService {

    #sessionReaderHttpRepository;

    constructor() {
        this.#sessionReaderHttpRepository = SessionReaderHttpRepository.getInstance();
    }

    /**
     * @returns {LogInService}
     */
    static getInstance() {
        return new LogInService();
    }

    /**
     * @param {import("./LogInDto.js").LogInDto} logInDto
     * @returns {Promise<LogInResultDto>}
     * @throws {SessionException}
     */
    async invoke(logInDto) {
        if (!logInDto.getDashboardToken()) {
            SessionException.badRequestCustom([SessionMessageEnum.DASHBOARD_TOKEN_IS_REQUIRED]);
        }

        const apiResponse = await this.#sessionReaderHttpRepository.getDashboardTokenCheck(logInDto.getDashboardToken());
        if (!apiResponse.status) {
            SessionException.fromApiResponse(apiResponse);
        }

        return LogInResultDto.fromPrimitives({
            dashboardToken: logInDto.getDashboardToken(),
        });
    }
}
