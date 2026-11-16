/**
 * Estado de un agente (lo calcula el backend: AgentStatusResolver). LABELS: lo que se ve en el filtro.
 */
export const AgentStatusEnum = Object.freeze({
    ONLINE: "online",
    OFFLINE: "offline",
    NEVER_SEEN: "never_seen",
    DISABLED: "disabled",
    VALUES: Object.freeze(["online", "offline", "never_seen", "disabled"]),
    LABELS: Object.freeze({
        online: "online",
        offline: "offline",
        never_seen: "nunca visto",
        disabled: "desactivado",
    }),
});
