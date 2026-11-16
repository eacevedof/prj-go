const ONLINE = "online";
const OFFLINE = "offline";
const NEVER_SEEN = "never_seen";
const DISABLED = "disabled";

/**
 * Estado de un agente (lo calcula el backend: AgentStatusResolver). LABELS: lo que se ve en el filtro.
 */
export const AgentStatusEnum = Object.freeze({
    ONLINE,
    OFFLINE,
    NEVER_SEEN,
    DISABLED,
    VALUES: Object.freeze([ONLINE, OFFLINE, NEVER_SEEN, DISABLED]),
    LABELS: Object.freeze({
        [ONLINE]: "online",
        [OFFLINE]: "offline",
        [NEVER_SEEN]: "nunca visto",
        [DISABLED]: "desactivado",
    }),
});
