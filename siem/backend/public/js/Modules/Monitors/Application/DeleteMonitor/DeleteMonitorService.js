import { MonitorsMessageEnum } from "/js/Modules/Monitors/Domain/Enums/MonitorsMessageEnum.js";
import { MonitorsException } from "/js/Modules/Monitors/Domain/Exceptions/MonitorsException.js";
import { DeleteMonitorResultDto } from "/js/Modules/Monitors/Application/DeleteMonitor/DeleteMonitorResultDto.js";
import { MonitorsWriterHttpRepository } from "/js/Modules/Monitors/Infrastructure/Repositories/MonitorsWriterHttpRepository.js";

/**
 * Baja de un monitor: los agentes dejan de vigilarlo en su siguiente sincronizacion. Sus eventos ya guardados siguen.
 */
export class DeleteMonitorService {

    #monitorsWriterHttpRepository;

    constructor() {
        this.#monitorsWriterHttpRepository = MonitorsWriterHttpRepository.getInstance();
    }

    /**
     * @returns {DeleteMonitorService}
     */
    static getInstance() {
        return new DeleteMonitorService();
    }

    /**
     * @param {import("./DeleteMonitorDto.js").DeleteMonitorDto} deleteMonitorDto
     * @returns {Promise<DeleteMonitorResultDto>}
     * @throws {MonitorsException}
     */
    async invoke(deleteMonitorDto) {
        if (!deleteMonitorDto.getMonitorUuid()) {
            MonitorsException.badRequestCustom([MonitorsMessageEnum.MONITOR_UUID_IS_REQUIRED]);
        }

        const apiResponse = await this.#monitorsWriterHttpRepository.deleteMonitorByMonitorUuid(
            deleteMonitorDto.getMonitorUuid(), deleteMonitorDto.getDashboardToken()
        );
        if (!apiResponse.status) {
            MonitorsException.fromApiResponse(apiResponse);
        }
        return DeleteMonitorResultDto.fromPrimitives({ monitorUuid: deleteMonitorDto.getMonitorUuid() });
    }
}
