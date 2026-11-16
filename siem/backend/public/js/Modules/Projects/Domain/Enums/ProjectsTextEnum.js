/**
 * Textos del camino feliz de la pantalla de proyectos (los ve una persona: en español). {0} lo rellena la vista.
 */
export const ProjectsTextEnum = Object.freeze({
    FORM_TITLE_CREATE: "Nuevo proyecto",
    FORM_TITLE_UPDATE: "Editar proyecto #{0}",
    SUBMIT_CREATE: "Crear proyecto",
    SUBMIT_UPDATE: "Guardar cambios",
    PROJECT_CREATED: "Proyecto \"{0}\" creado. Entra en él para dar de alta sus agentes y monitores.",
    PROJECT_UPDATED: "Proyecto \"{0}\" guardado.",
    LIST_EMPTY: "Sin proyectos todavía. Crea el primero con \"Nuevo proyecto\".",
    LIST_NO_MATCHES: "Ningún proyecto cumple los filtros.",
    COUNT_ALL: "({0})",
    COUNT_FILTERED: "({0} de {1})",
    HEALTH_ALL: "Todos",
    NO_DESCRIPTION: "sin descripción",
    EDIT: "Editar",
    DELETE: "Borrar",
    // aviso tras borrar (ProjectDeletedView): el borrado no desinstala los siem-client de sus servidores
    DELETED_TITLE: "Proyecto \"{0}\" borrado, con sus agentes, monitores y eventos.",
    DELETED_WITHOUT_AGENTS: "No tenía agentes: no queda nada instalado.",
    DELETED_UNINSTALL_ONE: "Recuerda desinstalar su cliente: sigue instalado en su servidor (el panel ya lo rechaza).",
    DELETED_UNINSTALL_MANY: "Recuerda desinstalar sus {0} clientes: siguen instalados en sus servidores (el panel ya los rechaza).",
    DELETED_AGENT_ON_HOST: "{0} en {1}",
    DELETED_HOST_UNKNOWN: "un servidor que nunca llegó a enviar nada",
    DELETED_PATH_DEFAULT: "ruta por defecto del instalador: el cliente no la informó (anterior a 0.4.0)",
    DELETED_NEVER_SEEN: "No llegó a conectar: si se instaló, estará en {0} (Linux) o {1} (Windows).",
    DELETED_COMMAND_UNIX: "bash",
    DELETED_COMMAND_WINDOWS: "PowerShell como administrador",
    DELETED_CLOSE: "Entendido",
});
