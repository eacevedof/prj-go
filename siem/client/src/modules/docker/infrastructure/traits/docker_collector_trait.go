package traits

import (
    agentconfigentities "siem-client/src/modules/agentconfig/domain/entities"
    agentconfigtraits "siem-client/src/modules/agentconfig/infrastructure/traits"
    collectorsservices "siem-client/src/modules/collectors/domain/services"
    "siem-client/src/core/config"
    "siem-client/src/modules/docker/domain/enums"
    "siem-client/src/modules/docker/domain/exceptions"
    "siem-client/src/modules/shared/infrastructure/components/uuid"
    "siem-client/src/modules/shared/infrastructure/repositories/spool"
    "siem-client/src/modules/shared/infrastructure/repositories/state"
)

// DockerCollectorTrait lo comun de los cuatro colectores de Docker (estados, eventos, recursos, logs): la config
// (se relee en cada pasada), el estado entre pasadas, el spool y el error de "Docker no responde".
type DockerCollectorTrait struct {
    agentconfigtraits.AgentConfigLoaderTrait

    uuidGenerator             *uuid.UuidGenerator
    eventsBatcher             *collectorsservices.EventsBatcher
    StateReaderJsonRepository *state.StateReaderJsonRepository
    StateWriterJsonRepository *state.StateWriterJsonRepository
    spoolWriterDiskRepository *spool.SpoolWriterDiskRepository
}

func NewDockerCollectorTrait() DockerCollectorTrait {
    return DockerCollectorTrait{
        AgentConfigLoaderTrait:    agentconfigtraits.NewAgentConfigLoaderTrait(),
        uuidGenerator:             uuid.GetInstance(),
        eventsBatcher:             collectorsservices.EventsBatcherGetInstance(),
        StateReaderJsonRepository: state.StateReaderJsonRepositoryGetInstance(),
        StateWriterJsonRepository: state.StateWriterJsonRepositoryGetInstance(),
        spoolWriterDiskRepository: spool.SpoolWriterDiskRepositoryGetInstance(),
    }
}

// GetValidConfigOrFail la config entera (el limite de eventos por lote) y su bloque docker.
func (t DockerCollectorTrait) GetValidConfigOrFail() (*agentconfigentities.AgentConfigEntity, *agentconfigentities.DockerConfigEntity, error) {
    agentConfigEntity, err := t.GetValidAgentConfigEntityOrFail()
    if err != nil {
        return nil, nil, err
    }
    return agentConfigEntity, agentConfigEntity.GetDockerConfigEntity(), nil
}

// GetDockerUnavailableException el socket no esta montado, el daemon no responde o contesta algo inesperado.
func (t DockerCollectorTrait) GetDockerUnavailableException(socketPath string, err error) error {
    return exceptions.ServiceUnavailableCustom(config.Tr(enums.DockerMessageEnum.DOCKER_UNAVAILABLE_X_X, socketPath, err.Error()))
}

// SaveEventsInSpool en lotes de como mucho maxEventsPerBatch; devuelve cuantos lotes.
func (t DockerCollectorTrait) SaveEventsInSpool(events []map[string]any, maxEventsPerBatch int) (int, error) {
    spooledBatches := 0
    for _, batchEvents := range t.eventsBatcher.GetBatches(events, maxEventsPerBatch) {
        if err := t.spoolWriterDiskRepository.SaveSpoolBatch(t.GetSpoolPath(), t.uuidGenerator.GetUniqueIdWithPrefix("b-"), batchEvents); err != nil {
            return spooledBatches, err
        }
        spooledBatches++
    }
    return spooledBatches, nil
}
