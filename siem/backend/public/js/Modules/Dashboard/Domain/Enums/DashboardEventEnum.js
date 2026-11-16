export const DashboardEventEnum = Object.freeze({
    // {projects}: la portada ya esta pintada (objetos planos de la API, tambien en cada auto-refresco). La escucha
    // Projects para pintar su tabla con filtros
    DASHBOARD_RENDERED: "dashboard.dashboard-rendered",
    // {project, agents, serverNow}: el detalle de un proyecto ya esta pintado (objetos planos de la API). Lo escuchan los modulos
    // que cuelgan de esa pantalla: Projects (borrar), Agents (su tabla y el formulario), Monitors (a que agentes se puede asignar)
    PROJECT_DETAIL_RENDERED: "dashboard.project-detail-rendered",
});
