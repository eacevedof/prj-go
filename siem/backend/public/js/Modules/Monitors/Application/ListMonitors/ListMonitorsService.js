import { MonitorsMessageEnum } from "/js/Modules/Monitors/Domain/Enums/MonitorsMessageEnum.js";
import { MonitorsException } from "/js/Modules/Monitors/Domain/Exceptions/MonitorsException.js";
import { ListMonitorsResultDto } from "/js/Modules/Monitors/Application/ListMonitors/ListMonitorsResultDto.js";
import { MonitorsReaderHttpRepository } from "/js/Modules/Monitors/Infrastructure/Repositories/MonitorsReaderHttpRepository.js";

/**
 * Los monitores (recursos a vigilar) de un proyecto con su ultimo valor por agente.
 */
export class ListMonitorsService {

    #monitorsReaderHttpRepository;

    constructor() {
        this.#monitorsReaderHttpRepository = MonitorsReaderHttpRepository.getInstance();
    }

    /**
     * @returns {ListMonitorsService}
     */
    static getInstance() {
        return new ListMonitorsService();
    }

    /**
     * @param {import("./ListMonitorsDto.js").ListMonitorsDto} listMonitorsDto
     * @returns {Promise<ListMonitorsResultDto>}
     * @throws {MonitorsException}
     */
    async invoke(listMonitorsDto) {
        if (!listMonitorsDto.getProjectUuid()) {
            MonitorsException.badRequestCustom([MonitorsMessageEnum.PROJECT_UUID_IS_REQUIRED]);
        }

        const apiResponse = await this.#monitorsReaderHttpRepository.getMonitorsByProjectUuid(
            listMonitorsDto.getProjectUuid(), listMonitorsDto.getDashboardToken()
        );
        if (!apiResponse.status) {
            MonitorsException.fromApiResponse(apiResponse);
        }
        return ListMonitorsResultDto.fromPrimitives({ monitors: apiResponse.data?.monitors });
    }
}
