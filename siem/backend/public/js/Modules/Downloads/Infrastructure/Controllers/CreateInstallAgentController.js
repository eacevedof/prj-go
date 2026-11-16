import { HttpCodeEnum } from "/js/Modules/Shared/Domain/Enums/HttpCodeEnum.js";
import { Dispatcher } from "/js/Modules/Shared/Infrastructure/Components/Dispatcher.js";
import { SessionStorer } from "/js/Modules/Shared/Infrastructure/Components/SessionStorer.js";
import { Subscriber } from "/js/Modules/Shared/Infrastructure/Components/Subscriber.js";

import { SessionEventEnum } from "/js/Modules/Session/Domain/Enums/SessionEventEnum.js";

import { AgentsMessageEnum } from "/js/Modules/Agents/Domain/Enums/AgentsMessageEnum.js";
import { CreateAgentDto } from "/js/Modules/Agents/Application/CreateAgent/CreateAgentDto.js";
import { CreateAgentService } from "/js/Modules/Agents/Application/CreateAgent/CreateAgentService.js";
import { AgentsException } from "/js/Modules/Agents/Domain/Exceptions/AgentsException.js";

import { DownloadsEventEnum } from "/js/Modules/Downloads/Domain/Enums/DownloadsEventEnum.js";
import { DownloadsView } from "/js/Modules/Downloads/Infrastructure/Views/DownloadsView.js";

/**
 * Paso 2 del asistente de instalacion (Descargas): alta del agente ahi mismo, con el mismo caso de uso que el detalle
 * del proyecto (CreateAgent). Su api_token, que la API da una sola vez, entra directamente en los comandos del paso 3.
 */
export class CreateInstallAgentController {

    #dispatcher;
    #sessionStorer;
    #subscriber;
    #downloadsView;

    constructor() {
        this.#dispatcher = Dispatcher.getInstance();
        this.#sessionStorer = SessionStorer.getInstance();
        this.#subscriber = Subscriber.getInstance();
        this.#downloadsView = DownloadsView.getInstance();
    }

    /**
     * @returns {CreateInstallAgentController}
     */
    static getInstance() {
        return new CreateInstallAgentController();
    }

    init() {
        this.#subscriber.subscribe(DownloadsEventEnum.CREATE_AGENT_REQUESTED, (request) => this.invoke(request));
    }

    /**
     * @param {{projectUuid: string, agentName: string}} request
     */
    async invoke(request) {
        this.#downloadsView.renderAgentBusy(true);
        try {
            const createAgentResultDto = await CreateAgentService.getInstance().invoke(
                CreateAgentDto.fromPrimitives({
                    dashboardToken: this.#sessionStorer.getDashboardToken(),
                    projectUuid: request.projectUuid,
                    agentName: request.agentName,
                })
            );
            this.#downloadsView.renderCreatedAgent({
                agentName: createAgentResultDto.getAgentName(),
                apiToken: createAgentResultDto.getApiToken(),
                endpointUrl: createAgentResultDto.getEndpointUrl(),
            });
        }
        catch (error) {
            if (error instanceof AgentsException && error.getCode() === HttpCodeEnum.UNAUTHORIZED) {
                this.#dispatcher.dispatch(SessionEventEnum.SESSION_EXPIRED);
                return;
            }
            if (error instanceof AgentsException) {
                this.#downloadsView.renderAgentErrors(error.getErrors());
                return;
            }
            console.error("CreateInstallAgentController.invoke", error);
            this.#downloadsView.renderAgentErrors([AgentsMessageEnum.UNEXPECTED_ERROR]);
        }
        finally {
            this.#downloadsView.renderAgentBusy(false);
        }
    }
}
