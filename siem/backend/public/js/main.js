import { SessionStorer } from "/js/Modules/Shared/Infrastructure/Components/SessionStorer.js";
import { Dispatcher } from "/js/Modules/Shared/Infrastructure/Components/Dispatcher.js";
import { RefreshScreenController } from "/js/Modules/Shared/Infrastructure/Controllers/RefreshScreenController.js";
import { CopyView } from "/js/Modules/Shared/Infrastructure/Views/CopyView.js";
import { NavigationView } from "/js/Modules/Shared/Infrastructure/Views/NavigationView.js";

import { SessionEventEnum } from "/js/Modules/Session/Domain/Enums/SessionEventEnum.js";
import { LogInController } from "/js/Modules/Session/Infrastructure/Controllers/LogInController.js";
import { LogOutController } from "/js/Modules/Session/Infrastructure/Controllers/LogOutController.js";
import { StartSessionController } from "/js/Modules/Session/Infrastructure/Controllers/StartSessionController.js";
import { AuthView } from "/js/Modules/Session/Infrastructure/Views/AuthView.js";
import { SessionView } from "/js/Modules/Session/Infrastructure/Views/SessionView.js";

import { GetDashboardController } from "/js/Modules/Dashboard/Infrastructure/Controllers/GetDashboardController.js";
import { GetProjectDetailController } from "/js/Modules/Dashboard/Infrastructure/Controllers/GetProjectDetailController.js";
import { DashboardView } from "/js/Modules/Dashboard/Infrastructure/Views/DashboardView.js";
import { ProjectDetailView } from "/js/Modules/Dashboard/Infrastructure/Views/ProjectDetailView.js";

import { CreateProjectController } from "/js/Modules/Projects/Infrastructure/Controllers/CreateProjectController.js";
import { DeleteProjectController } from "/js/Modules/Projects/Infrastructure/Controllers/DeleteProjectController.js";
import { ProjectDangerView } from "/js/Modules/Projects/Infrastructure/Views/ProjectDangerView.js";
import { ProjectFormView } from "/js/Modules/Projects/Infrastructure/Views/ProjectFormView.js";

import { CreateAgentController } from "/js/Modules/Agents/Infrastructure/Controllers/CreateAgentController.js";
import { DeleteAgentController } from "/js/Modules/Agents/Infrastructure/Controllers/DeleteAgentController.js";
import { AgentFormView } from "/js/Modules/Agents/Infrastructure/Views/AgentFormView.js";

import { CreateMonitorController } from "/js/Modules/Monitors/Infrastructure/Controllers/CreateMonitorController.js";
import { DeleteMonitorController } from "/js/Modules/Monitors/Infrastructure/Controllers/DeleteMonitorController.js";
import { ListMonitorsController } from "/js/Modules/Monitors/Infrastructure/Controllers/ListMonitorsController.js";
import { UpdateMonitorController } from "/js/Modules/Monitors/Infrastructure/Controllers/UpdateMonitorController.js";
import { MonitorsView } from "/js/Modules/Monitors/Infrastructure/Views/MonitorsView.js";

import { SearchEventsController } from "/js/Modules/Events/Infrastructure/Controllers/SearchEventsController.js";
import { EventsView } from "/js/Modules/Events/Infrastructure/Views/EventsView.js";

import { GetDownloadsController } from "/js/Modules/Downloads/Infrastructure/Controllers/GetDownloadsController.js";
import { DownloadsView } from "/js/Modules/Downloads/Infrastructure/Views/DownloadsView.js";

/**
 * Arranque de la SPA (patron task-js, sin build ni dependencias): vistas (DOM) -> controllers (suscritos a sus
 * eventos) -> sesion. El token guardado en la pestaña se revalida contra la API; sin token, pantalla de entrada.
 * Es el punto de entrada: captura lo que no tenga a nadie mas que lo capture.
 */
const start = async function () {
    document.getElementById("p-app-loading")?.classList.add("hidden");
    [
        NavigationView, CopyView, AuthView, SessionView,
        DashboardView, ProjectDetailView, ProjectFormView, ProjectDangerView, AgentFormView, MonitorsView,
        EventsView, DownloadsView,
    ].forEach((view) => view.getInstance().init());
    [
        RefreshScreenController,
        LogInController, StartSessionController, LogOutController,
        GetDashboardController, GetProjectDetailController,
        CreateProjectController, DeleteProjectController,
        CreateAgentController, DeleteAgentController,
        ListMonitorsController, CreateMonitorController, UpdateMonitorController, DeleteMonitorController,
        SearchEventsController,
        GetDownloadsController,
    ].forEach((controller) => controller.getInstance().init());

    const dashboardToken = SessionStorer.getInstance().getDashboardToken();
    if (!dashboardToken) {
        AuthView.getInstance().show();
        return;
    }
    Dispatcher.getInstance().dispatch(SessionEventEnum.LOG_IN_REQUESTED, { dashboardToken });
};

start().catch((error) => {
    console.error("main.start", error);
});
