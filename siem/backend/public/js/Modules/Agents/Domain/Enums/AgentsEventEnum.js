export const AgentsEventEnum = Object.freeze({
    // {projectUuid, agentName, heartbeatIntervalSeconds}
    CREATE_AGENT_REQUESTED: "agents.create-agent-requested",
    // {agentUuid}: ya confirmado en la vista (segundo clic)
    DELETE_AGENT_REQUESTED: "agents.delete-agent-requested",
});
