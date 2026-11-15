export const MonitorsEventEnum = Object.freeze({
    // {projectUuid, ...campos del formulario}
    CREATE_MONITOR_REQUESTED: "monitors.create-monitor-requested",
    // {monitorUuid, ...campos del formulario}: editar, o activar/desactivar desde la tabla
    UPDATE_MONITOR_REQUESTED: "monitors.update-monitor-requested",
    // {monitorUuid}: ya confirmado en la tabla (segundo clic)
    DELETE_MONITOR_REQUESTED: "monitors.delete-monitor-requested",
});
