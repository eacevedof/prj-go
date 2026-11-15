import { HttpCodeEnum } from "/js/Modules/Shared/Domain/Enums/HttpCodeEnum.js";
import { Dispatcher } from "/js/Modules/Shared/Infrastructure/Components/Dispatcher.js";
import { SessionStorer } from "/js/Modules/Shared/Infrastructure/Components/SessionStorer.js";
import { Subscriber } from "/js/Modules/Shared/Infrastructure/Components/Subscriber.js";
import { NavigationView } from "/js/Modules/Shared/Infrastructure/Views/NavigationView.js";

import { SessionEventEnum } from "/js/Modules/Session/Domain/Enums/SessionEventEnum.js";

import { AgentsEventEnum } from "/js/Modules/Agents/Domain/Enums/AgentsEventEnum.js";
import { AgentsMessageEnum } from "/js/Modules/Agents/Domain/Enums/AgentsMessageEnum.js";
import { DeleteAgentDto } from "/js/Modules/Agents/Application/DeleteAgent/DeleteAgentDto.js";
import { DeleteAgentService } from "/js/Modules/Agents/Application/DeleteAgent/DeleteAgentService.js";
import { AgentFormView } from "/js/Modules/Agents/Infrastructure/Views/AgentFormView.js";
import { AgentsException } from "/js/Modules/Agents/Domain/Exceptions/AgentsException.js";

/**
 * Baja de un agente (ya confirmada en su tarjeta): borrado, se repinta el detalle del proyecto.
 */
export class DeleteAgentController {

    #dispatcher;
    #sessionStorer;
    #subscriber;
    #navigationView;
    #agentFormView;

    constructor() {
        this.#dispatcher = Dispatcher.getInstance();
        this.#sessionStorer = SessionStorer.getInstance();
        this.#subscriber = Subscriber.getInstance();
        this.#navigationView = NavigationView.getInstance();
        this.#agentFormView = AgentFormView.getInstance();
    }

    /**
     * @returns {DeleteAgentController}
     */
    static getInstance() {
        return new DeleteAgentController();
    }

    init() {
        this.#subscriber.subscribe(AgentsEventEnum.DELETE_AGENT_REQUESTED, (request) => this.invoke(request));
    }

    /**
     * @param {{agentUuid: string}} request
     */
    async invoke(request) {
        try {
            await DeleteAgentService.getInstance().invoke(
                DeleteAgentDto.fromPrimitives({
                    dashboardToken: this.#sessionStorer.getDashboardToken(),
                    agentUuid: request.agentUuid,
                })
            );
            this.#agentFormView.renderListErrors([]);
            this.#navigationView.renderCurrentScreen();
        }
        catch (error) {
            if (error instanceof AgentsException && error.getCode() === HttpCodeEnum.UNAUTHORIZED) {
                this.#dispatcher.dispatch(SessionEventEnum.SESSION_EXPIRED);
                return;
            }
            if (error instanceof AgentsException) {
                this.#agentFormView.renderListErrors(error.getErrors());
                return;
            }
            console.error("DeleteAgentController.invoke", error);
            this.#agentFormView.renderListErrors([AgentsMessageEnum.UNEXPECTED_ERROR]);
        }
    }
}
