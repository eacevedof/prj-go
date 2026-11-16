/**
 * Textos del camino feliz de la tabla y el formulario de agentes (los ve una persona: en español). {0} lo rellena
 * TextFormatter.
 */
export const AgentsTextEnum = Object.freeze({
    FORM_TITLE_CREATE: "Nuevo agente",
    FORM_TITLE_UPDATE: "Editar agente #{0}",
    SUBMIT_CREATE: "Crear agente",
    SUBMIT_UPDATE: "Guardar cambios",
    AGENT_UPDATED: "Agente \"{0}\" guardado.",
    AGENT_DELETED: "Agente \"{0}\" borrado, con sus fuentes y eventos: su token ya no vale.",
    LIST_EMPTY: "Sin agentes. Da de alta el primero con \"Nuevo agente\": te dará el token para instalarlo.",
    LIST_NO_MATCHES: "Ningún agente cumple los filtros.",
    COUNT_ALL: "({0})",
    COUNT_FILTERED: "({0} de {1})",
    STATUS_ALL: "Todos",
    HOST_UNKNOWN: "host sin informar",
    NEVER_SEEN: "nunca",
    NO_SOURCES: "Aún no ha enviado eventos.",
    SHOW_DETAIL: "Ver detalle",
    HIDE_DETAIL: "Ocultar detalle",
    EDIT: "Editar",
    DELETE: "Borrar",
    EVENTS: "eventos",
    CONFIRM_DELETE: "¿Borrar con sus fuentes y eventos? Su token deja de valer.",
    CONFIRM_DELETE_YES: "Sí, borrar",
    CONFIRM_DELETE_NO: "Cancelar",
});
