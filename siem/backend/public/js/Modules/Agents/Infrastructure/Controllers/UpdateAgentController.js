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
import { UpdateAgentDto } from "/js/Modules/Agents/Application/UpdateAgent/UpdateAgentDto.js";
import { UpdateAgentService } from "/js/Modules/Agents/Application/UpdateAgent/UpdateAgentService.js";
import { AgentFormView } from "/js/Modules/Agents/Infrastructure/Views/AgentFormView.js";
import { AgentsTableView } from "/js/Modules/Agents/Infrastructure/Views/AgentsTableView.js";
import { AgentsException } from "/js/Modules/Agents/Domain/Exceptions/AgentsException.js";

/**
 * Edicion de un agente desde el detalle de su proyecto: "Editar" en la tabla lo carga en el formulario; guardado, el
 * formulario vuelve a alta y se repinta el detalle (la tabla ya sale con el cambio).
 */
export class UpdateAgentController {

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
     * @returns {UpdateAgentController}
     */
    static getInstance() {
        return new UpdateAgentController();
    }

    init() {
        this.#subscriber.subscribe(AgentsEventEnum.EDIT_AGENT_REQUESTED, (request) => this.#agentFormView.renderAgent(request.agent ?? {}));
        this.#subscriber.subscribe(AgentsEventEnum.UPDATE_AGENT_REQUESTED, (request) => this.invoke(request));
    }

    /**
     * @param {{agentUuid: string, agentName: string, heartbeatIntervalSeconds: string, isEnabled: boolean}} request
     */
    async invoke(request) {
        try {
            const updateAgentResultDto = await UpdateAgentService.getInstance().invoke(
                UpdateAgentDto.fromPrimitives({
                    dashboardToken: this.#sessionStorer.getDashboardToken(),
                    agentUuid: request.agentUuid,
                    agentName: request.agentName,
                    heartbeatIntervalSeconds: request.heartbeatIntervalSeconds,
                    isEnabled: request.isEnabled,
                })
            );
            this.#agentFormView.reset();
            this.#agentsTableView.renderMessage(
                this.#textFormatter.getFormatted(AgentsTextEnum.AGENT_UPDATED, updateAgentResultDto.getAgentName())
            );
            this.#navigationView.renderCurrentScreen();
        }
        catch (error) {
            this.#agentFormView.renderBusy(false);
            if (error instanceof AgentsException && error.getCode() === HttpCodeEnum.UNAUTHORIZED) {
                this.#dispatcher.dispatch(SessionEventEnum.SESSION_EXPIRED);
                return;
            }
            if (error instanceof AgentsException) {
                this.#agentFormView.renderErrors(error.getErrors());
                return;
            }
            console.error("UpdateAgentController.invoke", error);
            this.#agentFormView.renderErrors([AgentsMessageEnum.UNEXPECTED_ERROR]);
        }
    }
}
