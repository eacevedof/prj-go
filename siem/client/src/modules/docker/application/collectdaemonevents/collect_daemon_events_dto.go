package collectdaemonevents

// CollectDaemonEventsDto sin campos: que se vigila sale de siem-client.json.
type CollectDaemonEventsDto struct{}

func CollectDaemonEventsDtoFromPrimitives(primitives map[string]any) *CollectDaemonEventsDto {
    return &CollectDaemonEventsDto{}
}
