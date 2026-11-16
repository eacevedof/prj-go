import { HttpCodeEnum } from "/js/Modules/Shared/Domain/Enums/HttpCodeEnum.js";
import { Dispatcher } from "/js/Modules/Shared/Infrastructure/Components/Dispatcher.js";
import { SessionStorer } from "/js/Modules/Shared/Infrastructure/Components/SessionStorer.js";
import { Subscriber } from "/js/Modules/Shared/Infrastructure/Components/Subscriber.js";
import { TextFormatter } from "/js/Modules/Shared/Infrastructure/Components/TextFormatter.js";
import { NavigationView } from "/js/Modules/Shared/Infrastructure/Views/NavigationView.js";

import { SessionEventEnum } from "/js/Modules/Session/Domain/Enums/SessionEventEnum.js";

import { AgentsEventEnum } from "/js/Modules/Agents/Domain/Enums/AgentsEventEnum.js";
import { AgentsMessageEnum } from "/js/Modules/Agents/Domain/Enums/AgentsMessageEnum.js";
import { AgentsTextEnum } from "/js/Modules/Agents/Domain/Enums/AgentsTextEnum.js";
import { DeleteAgentDto } from "/js/Modules/Agents/Application/DeleteAgent/DeleteAgentDto.js";
import { DeleteAgentService } from "/js/Modules/Agents/Application/DeleteAgent/DeleteAgentService.js";
import { AgentFormView } from "/js/Modules/Agents/Infrastructure/Views/AgentFormView.js";
import { AgentsTableView } from "/js/Modules/Agents/Infrastructure/Views/AgentsTableView.js";
import { AgentsException } from "/js/Modules/Agents/Domain/Exceptions/AgentsException.js";

/**
 * Baja de un agente desde la tabla del detalle de su proyecto (ya confirmada en la fila): borrado, aviso y detalle
 * repintado. Si se estaba editando ese agente, el formulario vuelve a alta.
 */
export class DeleteAgentController {

    #dispatcher;
    #sessionStorer;
    #subscriber;
    #textFormatter;
    #navigationView;
    #agentFormView;
    #agentsTableView;

    constructor() {
        this.#dispatcher = Dispatcher.getInstance();
        this.#sessionStorer = SessionStorer.getInstance();
        this.#subscriber = Subscriber.getInstance();
        this.#textFormatter = TextFormatter.getInstance();
        this.#navigationView = NavigationView.getInstance();
        this.#agentFormView = AgentFormView.getInstance();
        this.#agentsTableView = AgentsTableView.getInstance();
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
     * @param {{agentUuid: string, agentName: string}} request
     */
    async invoke(request) {
        try {
            const deleteAgentResultDto = await DeleteAgentService.getInstance().invoke(
                DeleteAgentDto.fromPrimitives({
                    dashboardToken: this.#sessionStorer.getDashboardToken(),
                    agentUuid: request.agentUuid,
                })
            );
            this.#agentFormView.resetIfEditing(deleteAgentResultDto.getAgentUuid());
            this.#agentsTableView.renderMessage(
                this.#textFormatter.getFormatted(AgentsTextEnum.AGENT_DELETED, request.agentName)
            );
        }
        catch (error) {
            if (error instanceof AgentsException && error.getCode() === HttpCodeEnum.UNAUTHORIZED) {
                this.#dispatcher.dispatch(SessionEventEnum.SESSION_EXPIRED);
                return;
            }
            if (error instanceof AgentsException) {
                this.#agentsTableView.renderErrors(error.getErrors());
            }
            else {
                console.error("DeleteAgentController.invoke", error);
                this.#agentsTableView.renderErrors([AgentsMessageEnum.UNEXPECTED_ERROR]);
            }
        }
        this.#navigationView.renderCurrentScreen();
    }
}
