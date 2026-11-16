export const AgentsEventEnum = Object.freeze({
    // {projectUuid, agentName, heartbeatIntervalSeconds}
    CREATE_AGENT_REQUESTED: "agents.create-agent-requested",
    // {agentUuid, agentName, heartbeatIntervalSeconds, isEnabled}
    UPDATE_AGENT_REQUESTED: "agents.update-agent-requested",
    // {agent}: "Editar" en la tabla; el agente tal cual viene del detalle del proyecto, para cargarlo en el formulario
    EDIT_AGENT_REQUESTED: "agents.edit-agent-requested",
    // {agentUuid, agentName}: ya confirmado en la fila ("Sí, borrar")
    DELETE_AGENT_REQUESTED: "agents.delete-agent-requested",
    // {filters: {agentText, agentStatus}}: cambio un filtro de la tabla (se filtra en el front, sin pedir nada)
    AGENTS_FILTERS_CHANGED: "agents.agents-filters-changed",
});
