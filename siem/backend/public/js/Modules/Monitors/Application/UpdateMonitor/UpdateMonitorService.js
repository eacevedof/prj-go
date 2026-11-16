import { MonitorsMessageEnum } from "/js/Modules/Monitors/Domain/Enums/MonitorsMessageEnum.js";
import { MonitorValidator } from "/js/Modules/Monitors/Domain/Services/MonitorValidator.js";
import { UpdateMonitorResultDto } from "/js/Modules/Monitors/Application/UpdateMonitor/UpdateMonitorResultDto.js";
import { MonitorsWriterHttpRepository } from "/js/Modules/Monitors/Infrastructure/Repositories/MonitorsWriterHttpRepository.js";
import { MonitorsException } from "/js/Modules/Monitors/Domain/Exceptions/MonitorsException.js";

/**
 * Cambia un monitor (nombre, ambito, intervalo, umbrales, config, activo). El PUT lleva el monitor entero.
 */
export class UpdateMonitorService {

    #monitorValidator;
    #monitorsWriterHttpRepository;

    constructor() {
        this.#monitorValidator = MonitorValidator.getInstance();
        this.#monitorsWriterHttpRepository = MonitorsWriterHttpRepository.getInstance();
    }

    /**
     * @returns {UpdateMonitorService}
     */
    static getInstance() {
        return new UpdateMonitorService();
    }

    /**
     * @param {import("./UpdateMonitorDto.js").UpdateMonitorDto} updateMonitorDto
     * @returns {Promise<UpdateMonitorResultDto>}
     * @throws {MonitorsException}
     */
    async invoke(updateMonitorDto) {
        const errors = this.#monitorValidator.getErrors({
            monitorType: updateMonitorDto.getMonitorType(),
            monitorName: updateMonitorDto.getMonitorName(),
            intervalSeconds: updateMonitorDto.getIntervalSeconds(),
            warningThreshold: updateMonitorDto.getWarningThreshold(),
            criticalThreshold: updateMonitorDto.getCriticalThreshold(),
            config: updateMonitorDto.getConfig(),
        });
        if (!updateMonitorDto.getMonitorUuid()) {
            errors.unshift(MonitorsMessageEnum.MONITOR_UUID_IS_REQUIRED);
        }
        if (errors.length) {
            MonitorsException.badRequestCustom(errors);
        }

        const apiResponse = await this.#monitorsWriterHttpRepository.updateMonitorByMonitorUuid(
            updateMonitorDto.getMonitorUuid(), updateMonitorDto.getApiMonitor(), updateMonitorDto.getDashboardToken()
        );
        if (!apiResponse.status) {
            MonitorsException.fromApiResponse(apiResponse);
        }
        return UpdateMonitorResultDto.fromPrimitives({ monitor: apiResponse.data });
    }
}
