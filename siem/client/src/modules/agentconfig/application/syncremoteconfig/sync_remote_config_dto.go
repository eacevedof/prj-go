package syncremoteconfig

// SyncRemoteConfigDto sin campos: la url y el token salen del .env.
type SyncRemoteConfigDto struct{}

func SyncRemoteConfigDtoFromPrimitives(primitives map[string]any) *SyncRemoteConfigDto {
    return &SyncRemoteConfigDto{}
}
