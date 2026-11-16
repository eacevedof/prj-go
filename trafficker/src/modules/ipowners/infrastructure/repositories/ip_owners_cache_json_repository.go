package repositories

import (
    "trafficker/src/modules/ipowners/domain/entities"
    "trafficker/src/modules/ipowners/domain/enums"
    "trafficker/src/modules/shared/infrastructure/repositories/state"
)

// IpOwnersCacheJsonRepository la cache de whois en <state>/ip-owners.json: {"ranges": [ {start_ip, end_ip, ...} ]}.
type IpOwnersCacheJsonRepository struct {
    stateReaderJsonRepository *state.StateReaderJsonRepository
    stateWriterJsonRepository *state.StateWriterJsonRepository
}

func IpOwnersCacheJsonRepositoryGetInstance() *IpOwnersCacheJsonRepository {
    return &IpOwnersCacheJsonRepository{
        stateReaderJsonRepository: state.StateReaderJsonRepositoryGetInstance(),
        stateWriterJsonRepository: state.StateWriterJsonRepositoryGetInstance(),
    }
}

// GetIpOwnerEntities vacio si aun no hay cache; una entrada que no se entiende se ignora.
func (r *IpOwnersCacheJsonRepository) GetIpOwnerEntities(statePath string) ([]*entities.IpOwnerEntity, error) {
    cache, err := r.stateReaderJsonRepository.GetState(statePath, enums.IpOwnersCacheEnum.FILE_NAME)
    if err != nil {
        return nil, err
    }
    ipOwnerEntities := []*entities.IpOwnerEntity{}
    rawRanges, _ := cache[enums.IpOwnersCacheEnum.RANGES].([]any)
    for _, rawRange := range rawRanges {
        rangePrimitives, _ := rawRange.(map[string]any)
        if ipOwnerEntity := entities.IpOwnerEntityFromCachePrimitives(rangePrimitives); ipOwnerEntity != nil {
            ipOwnerEntities = append(ipOwnerEntities, ipOwnerEntity)
        }
    }
    return ipOwnerEntities, nil
}

func (r *IpOwnersCacheJsonRepository) SaveIpOwnerEntities(statePath string, ipOwnerEntities []*entities.IpOwnerEntity) error {
    ranges := make([]map[string]any, 0, len(ipOwnerEntities))
    for _, ipOwnerEntity := range ipOwnerEntities {
        ranges = append(ranges, ipOwnerEntity.ToCachePrimitives())
    }
    return r.stateWriterJsonRepository.SaveState(statePath, enums.IpOwnersCacheEnum.FILE_NAME, map[string]any{
        enums.IpOwnersCacheEnum.RANGES: ranges,
    })
}
