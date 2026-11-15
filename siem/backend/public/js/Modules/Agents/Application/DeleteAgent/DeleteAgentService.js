import { AgentsMessageEnum } from "/js/Modules/Agents/Domain/Enums/AgentsMessageEnum.js";
import { AgentsException } from "/js/Modules/Agents/Domain/Exceptions/AgentsException.js";
import { DeleteAgentResultDto } from "/js/Modules/Agents/Application/DeleteAgent/DeleteAgentResultDto.js";
import { AgentsWriterHttpRepository } from "/js/Modules/Agents/Infrastructure/Repositories/AgentsWriterHttpRepository.js";

/**
 * Baja de un agente: la API borra tambien sus fuentes y eventos, y su token deja de valer.
 */
export class DeleteAgentService {

    #agentsWriterHttpRepository;

    constructor() {
        this.#agentsWriterHttpRepository = AgentsWriterHttpRepository.getInstance();
    }

    /**
     * @returns {DeleteAgentService}
     */
    static getInstance() {
        return new DeleteAgentService();
    }

    /**
     * @param {import("./DeleteAgentDto.js").DeleteAgentDto} deleteAgentDto
     * @returns {Promise<DeleteAgentResultDto>}
     * @throws {AgentsException}
     */
    async invoke(deleteAgentDto) {
        if (!deleteAgentDto.getAgentUuid()) {
            AgentsException.badRequestCustom([AgentsMessageEnum.AGENT_UUID_IS_REQUIRED]);
        }

        const apiResponse = await this.#agentsWriterHttpRepository.deleteAgentByAgentUuid(
            deleteAgentDto.getAgentUuid(), deleteAgentDto.getDashboardToken()
        );
        if (!apiResponse.status) {
            AgentsException.fromApiResponse(apiResponse);
        }
        return DeleteAgentResultDto.fromPrimitives({ agentUuid: deleteAgentDto.getAgentUuid() });
    }
}
