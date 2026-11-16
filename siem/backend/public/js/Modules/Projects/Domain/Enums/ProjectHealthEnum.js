/**
 * Semaforo de un proyecto (lo calcula el backend: ProjectHealthResolver). LABELS: lo que se ve en el filtro.
 */
export const ProjectHealthEnum = Object.freeze({
    OK: "ok",
    WARNING: "warning",
    CRITICAL: "critical",
    EMPTY: "empty",
    VALUES: Object.freeze(["ok", "warning", "critical", "empty"]),
    LABELS: Object.freeze({
        ok: "ok",
        warning: "con avisos",
        critical: "crítico",
        empty: "sin datos",
    }),
});
