import { AgentFieldsEnum } from "/js/Modules/Agents/Domain/Enums/AgentFieldsEnum.js";
import { AgentsMessageEnum } from "/js/Modules/Agents/Domain/Enums/AgentsMessageEnum.js";
import { CreateAgentResultDto } from "/js/Modules/Agents/Application/CreateAgent/CreateAgentResultDto.js";
import { AgentsWriterHttpRepository } from "/js/Modules/Agents/Infrastructure/Repositories/AgentsWriterHttpRepository.js";
import { AgentsException } from "/js/Modules/Agents/Domain/Exceptions/AgentsException.js";

/**
 * Alta de un agente (un host con siem-client) en un proyecto. La API devuelve su api_token una sola vez: de ella solo
 * se guarda el hash.
 */
export class CreateAgentService {

    #agentsWriterHttpRepository;

    constructor() {
        this.#agentsWriterHttpRepository = AgentsWriterHttpRepository.getInstance();
    }

    /**
     * @returns {CreateAgentService}
     */
    static getInstance() {
        return new CreateAgentService();
    }

    /**
     * @param {import("./CreateAgentDto.js").CreateAgentDto} createAgentDto
     * @returns {Promise<CreateAgentResultDto>}
     * @throws {AgentsException}
     */
    async invoke(createAgentDto) {
        const errors = [];
        if (!createAgentDto.getProjectUuid()) {
            errors.push(AgentsMessageEnum.PROJECT_UUID_IS_REQUIRED);
        }
        if (!AgentFieldsEnum.AGENT_NAME_PATTERN.test(createAgentDto.getAgentName())) {
            errors.push(AgentsMessageEnum.AGENT_NAME_IS_NOT_VALID);
        }
        const heartbeatIntervalSeconds = createAgentDto.getHeartbeatIntervalSeconds();
        if (!Number.isInteger(heartbeatIntervalSeconds) || heartbeatIntervalSeconds < AgentFieldsEnum.HEARTBEAT_MIN_SECONDS
            || heartbeatIntervalSeconds > AgentFieldsEnum.HEARTBEAT_MAX_SECONDS) {
            errors.push(AgentsMessageEnum.HEARTBEAT_INTERVAL_MUST_BE_10_TO_86400);
        }
        if (errors.length) {
            AgentsException.badRequestCustom(errors);
        }

        const apiResponse = await this.#agentsWriterHttpRepository.createAgentInProject(
            createAgentDto.getProjectUuid(),
            { agent_name: createAgentDto.getAgentName(), heartbeat_interval_seconds: heartbeatIntervalSeconds },
            createAgentDto.getDashboardToken()
        );
        if (!apiResponse.status) {
            AgentsException.fromApiResponse(apiResponse);
        }

        const agent = apiResponse.data ?? {};
        return CreateAgentResultDto.fromPrimitives({
            agentUuid: agent.agent_uuid,
            agentName: agent.agent_name,
            apiToken: agent.api_token,
            endpointUrl: agent.endpoint_url,
        });
    }
}
