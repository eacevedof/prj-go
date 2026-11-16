import { MonitorsMessageEnum } from "/js/Modules/Monitors/Domain/Enums/MonitorsMessageEnum.js";
import { MonitorValidator } from "/js/Modules/Monitors/Domain/Services/MonitorValidator.js";
import { CreateMonitorResultDto } from "/js/Modules/Monitors/Application/CreateMonitor/CreateMonitorResultDto.js";
import { MonitorsWriterHttpRepository } from "/js/Modules/Monitors/Infrastructure/Repositories/MonitorsWriterHttpRepository.js";
import { MonitorsException } from "/js/Modules/Monitors/Domain/Exceptions/MonitorsException.js";

/**
 * Alta de un monitor (recurso a vigilar) en un proyecto. Los agentes del proyecto lo descargan solos en su siguiente
 * sincronizacion: no hay que tocar nada en los servidores.
 */
export class CreateMonitorService {

    #monitorValidator;
    #monitorsWriterHttpRepository;

    constructor() {
        this.#monitorValidator = MonitorValidator.getInstance();
        this.#monitorsWriterHttpRepository = MonitorsWriterHttpRepository.getInstance();
    }

    /**
     * @returns {CreateMonitorService}
     */
    static getInstance() {
        return new CreateMonitorService();
    }

    /**
     * @param {import("./CreateMonitorDto.js").CreateMonitorDto} createMonitorDto
     * @returns {Promise<CreateMonitorResultDto>}
     * @throws {MonitorsException}
     */
    async invoke(createMonitorDto) {
        const errors = this.#monitorValidator.getErrors({
            monitorType: createMonitorDto.getMonitorType(),
            monitorName: createMonitorDto.getMonitorName(),
            intervalSeconds: createMonitorDto.getIntervalSeconds(),
            warningThreshold: createMonitorDto.getWarningThreshold(),
            criticalThreshold: createMonitorDto.getCriticalThreshold(),
            config: createMonitorDto.getConfig(),
        });
        if (!createMonitorDto.getProjectUuid()) {
            errors.unshift(MonitorsMessageEnum.PROJECT_UUID_IS_REQUIRED);
        }
        if (errors.length) {
            MonitorsException.badRequestCustom(errors);
        }

        const apiResponse = await this.#monitorsWriterHttpRepository.createMonitorInProject(
            createMonitorDto.getProjectUuid(), createMonitorDto.getApiMonitor(), createMonitorDto.getDashboardToken()
        );
        if (!apiResponse.status) {
            MonitorsException.fromApiResponse(apiResponse);
        }
        return CreateMonitorResultDto.fromPrimitives({ monitor: apiResponse.data });
    }
}
