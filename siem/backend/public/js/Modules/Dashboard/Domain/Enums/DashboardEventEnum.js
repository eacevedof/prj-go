export const DashboardEventEnum = Object.freeze({
    // {project, agents}: el detalle de un proyecto ya esta pintado (objetos planos de la API). Lo escuchan los modulos
    // que cuelgan de esa pantalla: Projects (borrar), Agents (alta), Monitors (a que agentes se puede asignar)
    PROJECT_DETAIL_RENDERED: "dashboard.project-detail-rendered",
});
