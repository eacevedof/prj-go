export const ProjectsEventEnum = Object.freeze({
    // {projectName, projectDescription}
    CREATE_PROJECT_REQUESTED: "projects.create-project-requested",
    // {projectUuid, projectName, projectDescription}
    UPDATE_PROJECT_REQUESTED: "projects.update-project-requested",
    // {project}: "Editar" en la tabla; el proyecto tal cual viene del dashboard, para cargarlo en el formulario
    EDIT_PROJECT_REQUESTED: "projects.edit-project-requested",
    // {projectUuid}: desde el detalle, ya confirmado escribiendo el nombre; borrado, vuelve a la portada
    DELETE_PROJECT_REQUESTED: "projects.delete-project-requested",
    // {projectUuid, projectName}: desde la tabla de la portada, ya confirmado con "Sí, borrar"; se queda en la portada
    DELETE_LISTED_PROJECT_REQUESTED: "projects.delete-listed-project-requested",
    // {filters: {projectText, projectHealth}}: cambio un filtro de la tabla (se filtra en el front, sin pedir nada)
    PROJECTS_FILTERS_CHANGED: "projects.projects-filters-changed",
});
