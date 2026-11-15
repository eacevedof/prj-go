package collectfilelines

// CollectFileLinesDto sin campos: que carpetas se vigilan sale de siem-client.json (se relee en cada pasada).
type CollectFileLinesDto struct{}

func CollectFileLinesDtoFromPrimitives(primitives map[string]any) *CollectFileLinesDto {
    return &CollectFileLinesDto{}
}
