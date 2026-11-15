import { HttpCodeEnum } from "/js/Modules/Shared/Domain/Enums/HttpCodeEnum.js";
import { ScreenEnum } from "/js/Modules/Shared/Domain/Enums/ScreenEnum.js";
import { SharedEventEnum } from "/js/Modules/Shared/Domain/Enums/SharedEventEnum.js";
import { Dispatcher } from "/js/Modules/Shared/Infrastructure/Components/Dispatcher.js";
import { SessionStorer } from "/js/Modules/Shared/Infrastructure/Components/SessionStorer.js";
import { Subscriber } from "/js/Modules/Shared/Infrastructure/Components/Subscriber.js";
import { NavigationView } from "/js/Modules/Shared/Infrastructure/Views/NavigationView.js";

import { SessionEventEnum } from "/js/Modules/Session/Domain/Enums/SessionEventEnum.js";

import { AgentsEventEnum } from "/js/Modules/Agents/Domain/Enums/AgentsEventEnum.js";
import { AgentsMessageEnum } from "/js/Modules/Agents/Domain/Enums/AgentsMessageEnum.js";
import { CreateAgentDto } from "/js/Modules/Agents/Application/CreateAgent/CreateAgentDto.js";
import { CreateAgentService } from "/js/Modules/Agents/Application/CreateAgent/CreateAgentService.js";
import { AgentFormView } from "/js/Modules/Agents/Infrastructure/Views/AgentFormView.js";
import { AgentsException } from "/js/Modules/Agents/Domain/Exceptions/AgentsException.js";

/**
 * Alta de un agente desde el detalle de un proyecto: la vista enseña su api_token una sola vez y la pantalla se
 * repinta (el agente ya sale en la lista, como never_seen). Tambien le dice a la vista en que proyecto esta, y al
 * salir de sesion le quita el token de la vista.
 */
export class CreateAgentController {

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
     * @returns {CreateAgentController}
     */
    static getInstance() {
        return new CreateAgentController();
    }

    init() {
        this.#subscriber.subscribe(AgentsEventEnum.CREATE_AGENT_REQUESTED, (request) => this.invoke(request));
        this.#subscriber.subscribe(SharedEventEnum.SCREEN_REQUESTED, (request) => {
            if (request.screen !== ScreenEnum.PROJECT_DETAIL) return;
            this.#agentFormView.renderProject(request.pathParam ?? "");
        });
        this.#subscriber.subscribe(SessionEventEnum.USER_LOGGED_OUT, () => this.#agentFormView.renderProject(""));
    }

    /**
     * @param {{projectUuid: string, agentName: string, heartbeatIntervalSeconds: string}} request
     */
    async invoke(request) {
        this.#agentFormView.renderBusy(true);
        try {
            const createAgentResultDto = await CreateAgentService.getInstance().invoke(
                CreateAgentDto.fromPrimitives({
                    dashboardToken: this.#sessionStorer.getDashboardToken(),
                    projectUuid: request.projectUuid,
                    agentName: request.agentName,
                    heartbeatIntervalSeconds: request.heartbeatIntervalSeconds,
                })
            );
            this.#agentFormView.renderCreated({
                agentName: createAgentResultDto.getAgentName(),
                apiToken: createAgentResultDto.getApiToken(),
                endpointUrl: createAgentResultDto.getEndpointUrl(),
            });
            this.#navigationView.renderCurrentScreen();
        }
        catch (error) {
            if (error instanceof AgentsException && error.getCode() === HttpCodeEnum.UNAUTHORIZED) {
                this.#dispatcher.dispatch(SessionEventEnum.SESSION_EXPIRED);
                return;
            }
            if (error instanceof AgentsException) {
                this.#agentFormView.renderErrors(error.getErrors());
                return;
            }
            console.error("CreateAgentController.invoke", error);
            this.#agentFormView.renderErrors([AgentsMessageEnum.UNEXPECTED_ERROR]);
        }
        finally {
            this.#agentFormView.renderBusy(false);
        }
    }
}
