/**
 * Claves de un proyecto tal cual las da GET /api/v1/dashboard (snake_case).
 */
export const ProjectsKeyEnum = Object.freeze({
    PROJECT_ID: "project_id",
    PROJECT_UUID: "project_uuid",
    PROJECT_NAME: "project_name",
    PROJECT_DESCRIPTION: "project_description",
    PROJECT_HEALTH: "project_health",
    WARNINGS_24H: "warnings_24h",
    ERRORS_24H: "errors_24h",
    CRITICALS_24H: "criticals_24h",
    // DELETE /api/v1/projects/:project_uuid: los agentes que tenia el proyecto borrado
    AGENTS: "agents",
});
