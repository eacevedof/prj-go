/**
 * Claves de la API de dashboard (JSON de salida de /api/v1/dashboard y /api/v1/projects/:project_uuid).
 */
export const DashboardKeyEnum = Object.freeze({
    SERVER: "server",
    PROJECTS: "projects",
    NOW: "now",

    PROJECT: "project",
    AGENTS: "agents",

    PROJECT_UUID: "project_uuid",
    PROJECT_NAME: "project_name",
    PROJECT_DESCRIPTION: "project_description",
});
