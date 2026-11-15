package collectcontainerstates

// CollectContainerStatesDto sin campos: que se vigila sale de siem-client.json (se relee en cada pasada).
type CollectContainerStatesDto struct{}

func CollectContainerStatesDtoFromPrimitives(primitives map[string]any) *CollectContainerStatesDto {
    return &CollectContainerStatesDto{}
}
