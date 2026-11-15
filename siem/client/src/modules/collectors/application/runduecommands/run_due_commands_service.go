package runduecommands

import (
    "time"

    agentconfigentities "siem-client/src/modules/agentconfig/domain/entities"
    agentconfigtraits "siem-client/src/modules/agentconfig/infrastructure/traits"
    "siem-client/src/modules/collectors/domain/entities"
    "siem-client/src/modules/collectors/domain/enums"
    "siem-client/src/modules/collectors/domain/services"
    "siem-client/src/modules/collectors/infrastructure/repositories"
    "siem-client/src/modules/shared/application"
    "siem-client/src/modules/shared/infrastructure/components/uuid"
    "siem-client/src/modules/shared/infrastructure/repositories/spool"
    "siem-client/src/modules/shared/infrastructure/repositories/state"
)

const dateTimeFormat = "2006-01-02 15:04:05"

// RunDueCommandsService lanza los comandos programados a los que les toca (han pasado interval_seconds desde la
// ultima vez, o nunca se lanzaron) y deja su resultado en el spool como eventos command_result. Uno detras de otro:
// timeout_seconds acota lo que puede tardar cada uno. Cuando se lanzo cada uno queda en state/commands-last-run.json.
type RunDueCommandsService struct {
    agentconfigtraits.AgentConfigLoaderTrait

    uuidGenerator                       *uuid.UuidGenerator
    collectedEventBuilder               *services.CollectedEventBuilder
    commandResultsReaderShellRepository *repositories.CommandResultsReaderShellRepository
    stateReaderJsonRepository           *state.StateReaderJsonRepository
    stateWriterJsonRepository           *state.StateWriterJsonRepository
    spoolWriterDiskRepository           *spool.SpoolWriterDiskRepository

    runDueCommandsDto *RunDueCommandsDto
}

func RunDueCommandsServiceGetInstance() *RunDueCommandsService {
    return &RunDueCommandsService{
        AgentConfigLoaderTrait:              agentconfigtraits.NewAgentConfigLoaderTrait(),
        uuidGenerator:                       uuid.GetInstance(),
        collectedEventBuilder:               services.CollectedEventBuilderGetInstance(),
        commandResultsReaderShellRepository: repositories.CommandResultsReaderShellRepositoryGetInstance(),
        stateReaderJsonRepository:           state.StateReaderJsonRepositoryGetInstance(),
        stateWriterJsonRepository:           state.StateWriterJsonRepositoryGetInstance(),
        spoolWriterDiskRepository:           spool.SpoolWriterDiskRepositoryGetInstance(),
    }
}

// Invoke error: *AgentConfigException (400, 404) o uno inesperado (disco). Un comando que falla o no arranca no es
// un error del servicio: es un evento con exit_code distinto de 0.
func (s *RunDueCommandsService) Invoke(runDueCommandsDto *RunDueCommandsDto) (*RunDueCommandsResultDto, error) {
    s.runDueCommandsDto = runDueCommandsDto

    agentConfigEntity, err := s.GetValidAgentConfigEntityOrFail()
    if err != nil {
        return nil, err
    }
    lastRunByCommandName, err := s.stateReaderJsonRepository.GetState(s.GetStatePath(), enums.CollectorsStateFileEnum.COMMANDS_LAST_RUN)
    if err != nil {
        return nil, err
    }

    // solo los comandos que siguen en la configuracion: uno que se quita deja de figurar en el estado
    nextLastRunByCommandName := map[string]any{}
    events := []map[string]any{}
    failedCommands := 0
    for _, scheduledCommandEntity := range agentConfigEntity.GetScheduledCommandEntities() {
        commandName := scheduledCommandEntity.GetName()
        lastRunUnix := int64(application.ToInt(lastRunByCommandName[commandName]))
        nextLastRunByCommandName[commandName] = lastRunUnix

        now := time.Now().UTC()
        if !s.isCommandDue(scheduledCommandEntity, lastRunUnix, now) {
            continue
        }
        event := s.getCommandResultEvent(scheduledCommandEntity, now)
        if application.GetInt(event, enums.CollectorsKeyEnum.EXIT_CODE) != 0 {
            failedCommands++
        }
        events = append(events, event)
        nextLastRunByCommandName[commandName] = now.Unix()
    }

    spooledBatches := 0
    if len(events) > 0 {
        if err := s.spoolWriterDiskRepository.SaveSpoolBatch(s.GetSpoolPath(), s.uuidGenerator.GetUniqueIdWithPrefix("b-"), events); err != nil {
            return nil, err
        }
        spooledBatches = 1
    }
    if err := s.stateWriterJsonRepository.SaveState(s.GetStatePath(), enums.CollectorsStateFileEnum.COMMANDS_LAST_RUN, nextLastRunByCommandName); err != nil {
        return nil, err
    }

    return RunDueCommandsResultDtoFromPrimitives(map[string]any{
        "ranCommands":    len(events),
        "failedCommands": failedCommands,
        "spooledBatches": spooledBatches,
        "results":        events,
    }), nil
}

func (s *RunDueCommandsService) isCommandDue(scheduledCommandEntity *agentconfigentities.ScheduledCommandEntity, lastRunUnix int64, now time.Time) bool {
    if s.runDueCommandsDto.IsForced() || lastRunUnix == 0 {
        return true
    }
    return now.Unix()-lastRunUnix >= int64(scheduledCommandEntity.GetIntervalSeconds())
}

// getCommandResultEvent un comando que no llega a arrancar se anota como exit_code -1 con el motivo como salida.
func (s *RunDueCommandsService) getCommandResultEvent(scheduledCommandEntity *agentconfigentities.ScheduledCommandEntity, now time.Time) map[string]any {
    commandResultEntity, err := s.commandResultsReaderShellRepository.GetCommandResult(
        scheduledCommandEntity.GetCommandLine(), scheduledCommandEntity.GetTimeoutSeconds(),
    )
    if err != nil {
        commandResultEntity = entities.CommandResultEntityFromPrimitives(map[string]any{
            "output":   "[cannot start] " + err.Error(),
            "exitCode": -1,
        })
    }
    return s.collectedEventBuilder.GetCommandResultEvent(
        scheduledCommandEntity.GetName(), commandResultEntity, scheduledCommandEntity.GetTimeoutSeconds(), now.Format(dateTimeFormat),
    )
}
