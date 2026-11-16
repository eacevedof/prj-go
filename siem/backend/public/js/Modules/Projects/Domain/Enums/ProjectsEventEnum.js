export const ProjectsEventEnum = Object.freeze({
    // {projectName, projectDescription}
    CREATE_PROJECT_REQUESTED: "projects.create-project-requested",
    // {projectUuid, projectName, projectDescription}
    UPDATE_PROJECT_REQUESTED: "projects.update-project-requested",
    // {project}: "Editar" en la tabla; el proyecto tal cual viene del dashboard, para cargarlo en el formulario
    EDIT_PROJECT_REQUESTED: "projects.edit-project-requested",
    // {projectUuid, projectName}: "Borrar" en la tabla o en la cabecera del detalle; abre el modal que pide el nombre
    DELETE_PROJECT_DIALOG_REQUESTED: "projects.delete-project-dialog-requested",
    // {projectUuid}: desde el modal, ya confirmado escribiendo el nombre; borrado, vuelve a la portada con el aviso
    // de los siem-client que quedan por desinstalar
    DELETE_PROJECT_REQUESTED: "projects.delete-project-requested",
    // {filters: {projectText, projectHealth}}: cambio un filtro de la tabla (se filtra en el front, sin pedir nada)
    PROJECTS_FILTERS_CHANGED: "projects.projects-filters-changed",
});
