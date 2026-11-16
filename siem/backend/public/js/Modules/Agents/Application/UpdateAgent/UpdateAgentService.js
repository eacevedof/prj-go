import { AgentFieldsEnum } from "/js/Modules/Agents/Domain/Enums/AgentFieldsEnum.js";
import { AgentsMessageEnum } from "/js/Modules/Agents/Domain/Enums/AgentsMessageEnum.js";
import { AgentsException } from "/js/Modules/Agents/Domain/Exceptions/AgentsException.js";
import { UpdateAgentResultDto } from "/js/Modules/Agents/Application/UpdateAgent/UpdateAgentResultDto.js";
import { AgentsWriterHttpRepository } from "/js/Modules/Agents/Infrastructure/Repositories/AgentsWriterHttpRepository.js";

/**
 * Cambia nombre, latido y si esta activo un agente. Desactivado, su token deja de valer hasta que se reactive. Aqui
 * solo lo evidente: el nombre duplicado en el proyecto (409) lo decide la API.
 */
export class UpdateAgentService {

    #agentsWriterHttpRepository;

    constructor() {
        this.#agentsWriterHttpRepository = AgentsWriterHttpRepository.getInstance();
    }

    /**
     * @returns {UpdateAgentService}
     */
    static getInstance() {
        return new UpdateAgentService();
    }

    /**
     * @param {import("./UpdateAgentDto.js").UpdateAgentDto} updateAgentDto
     * @returns {Promise<UpdateAgentResultDto>}
     * @throws {AgentsException}
     */
    async invoke(updateAgentDto) {
        const errors = [];
        if (!updateAgentDto.getAgentUuid()) {
            errors.push(AgentsMessageEnum.AGENT_UUID_IS_REQUIRED);
        }
        if (!AgentFieldsEnum.AGENT_NAME_PATTERN.test(updateAgentDto.getAgentName())) {
            errors.push(AgentsMessageEnum.AGENT_NAME_IS_NOT_VALID);
        }
        const heartbeat = updateAgentDto.getHeartbeatIntervalSeconds();
        if (!Number.isInteger(heartbeat) || heartbeat < AgentFieldsEnum.HEARTBEAT_MIN_SECONDS || heartbeat > AgentFieldsEnum.HEARTBEAT_MAX_SECONDS) {
            errors.push(AgentsMessageEnum.HEARTBEAT_INTERVAL_MUST_BE_10_TO_86400);
        }
        if (errors.length) {
            AgentsException.badRequestCustom(errors);
        }

        const apiResponse = await this.#agentsWriterHttpRepository.updateAgentByAgentUuid(
            updateAgentDto.getAgentUuid(),
            {
                agent_name: updateAgentDto.getAgentName(),
                heartbeat_interval_seconds: heartbeat,
                is_enabled: updateAgentDto.isEnabled(),
            },
            updateAgentDto.getDashboardToken()
        );
        if (!apiResponse.status) {
            AgentsException.fromApiResponse(apiResponse);
        }

        const agent = apiResponse.data ?? {};
        return UpdateAgentResultDto.fromPrimitives({
            agentUuid: agent.agent_uuid,
            agentName: agent.agent_name,
        });
    }
}
