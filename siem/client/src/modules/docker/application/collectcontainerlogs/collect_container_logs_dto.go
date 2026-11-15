package collectcontainerlogs

// CollectContainerLogsDto sin campos: que contenedores se leen sale de siem-client.json (logs_include / logs_exclude).
type CollectContainerLogsDto struct{}

func CollectContainerLogsDtoFromPrimitives(primitives map[string]any) *CollectContainerLogsDto {
    return &CollectContainerLogsDto{}
}
