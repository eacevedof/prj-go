import { DashboardKeyEnum } from "/js/Modules/Dashboard/Domain/Enums/DashboardKeyEnum.js";
import { DashboardException } from "/js/Modules/Dashboard/Domain/Exceptions/DashboardException.js";
import { GetDashboardResultDto } from "/js/Modules/Dashboard/Application/GetDashboard/GetDashboardResultDto.js";
import { DashboardReaderHttpRepository } from "/js/Modules/Dashboard/Infrastructure/Repositories/DashboardReaderHttpRepository.js";

/**
 * La portada: estado del servidor y un resumen por proyecto. Los estados (salud, agentes online) los calcula la API:
 * aqui solo se pide y se da forma al resultado.
 */
export class GetDashboardService {

    #dashboardReaderHttpRepository;

    constructor() {
        this.#dashboardReaderHttpRepository = DashboardReaderHttpRepository.getInstance();
    }

    /**
     * @returns {GetDashboardService}
     */
    static getInstance() {
        return new GetDashboardService();
    }

    /**
     * @param {import("./GetDashboardDto.js").GetDashboardDto} getDashboardDto
     * @returns {Promise<GetDashboardResultDto>}
     * @throws {DashboardException}
     */
    async invoke(getDashboardDto) {
        const apiResponse = await this.#dashboardReaderHttpRepository.getDashboard(getDashboardDto.getDashboardToken());
        if (!apiResponse.status) {
            DashboardException.fromApiResponse(apiResponse);
        }

        const dashboard = apiResponse.data ?? {};
        return GetDashboardResultDto.fromPrimitives({
            server: dashboard[DashboardKeyEnum.SERVER],
            projects: dashboard[DashboardKeyEnum.PROJECTS],
        });
    }
}
