import { DashboardKeyEnum } from "/js/Modules/Dashboard/Domain/Enums/DashboardKeyEnum.js";
import { DashboardMessageEnum } from "/js/Modules/Dashboard/Domain/Enums/DashboardMessageEnum.js";
import { DashboardException } from "/js/Modules/Dashboard/Domain/Exceptions/DashboardException.js";
import { GetProjectDetailResultDto } from "/js/Modules/Dashboard/Application/GetProjectDetail/GetProjectDetailResultDto.js";
import { DashboardReaderHttpRepository } from "/js/Modules/Dashboard/Infrastructure/Repositories/DashboardReaderHttpRepository.js";

/**
 * Un proyecto con sus agentes y, por agente, cada carpeta y comando con su ultimo resultado. El detalle no trae el
 * now del servidor: se pide a la vez el dashboard (los "hace X" no pueden fiarse del reloj del navegador).
 */
export class GetProjectDetailService {

    #dashboardReaderHttpRepository;

    constructor() {
        this.#dashboardReaderHttpRepository = DashboardReaderHttpRepository.getInstance();
    }

    /**
     * @returns {GetProjectDetailService}
     */
    static getInstance() {
        return new GetProjectDetailService();
    }

    /**
     * @param {import("./GetProjectDetailDto.js").GetProjectDetailDto} getProjectDetailDto
     * @returns {Promise<GetProjectDetailResultDto>}
     * @throws {DashboardException}
     */
    async invoke(getProjectDetailDto) {
        if (!getProjectDetailDto.getProjectUuid()) {
            DashboardException.badRequestCustom([DashboardMessageEnum.PROJECT_UUID_IS_REQUIRED]);
        }

        const dashboardToken = getProjectDetailDto.getDashboardToken();
        const [projectDetailResponse, dashboardResponse] = await Promise.all([
            this.#dashboardReaderHttpRepository.getProjectDetailByProjectUuid(getProjectDetailDto.getProjectUuid(), dashboardToken),
            this.#dashboardReaderHttpRepository.getDashboard(dashboardToken),
        ]);
        if (!projectDetailResponse.status) {
            DashboardException.fromApiResponse(projectDetailResponse);
        }
        if (!dashboardResponse.status) {
            DashboardException.fromApiResponse(dashboardResponse);
        }

        const projectDetail = projectDetailResponse.data ?? {};
        const server = dashboardResponse.data?.[DashboardKeyEnum.SERVER] ?? {};
        return GetProjectDetailResultDto.fromPrimitives({
            project: projectDetail[DashboardKeyEnum.PROJECT],
            agents: projectDetail[DashboardKeyEnum.AGENTS],
            serverNow: server[DashboardKeyEnum.NOW],
        });
    }
}
