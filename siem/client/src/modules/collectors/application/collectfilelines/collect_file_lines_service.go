package collectfilelines

import (
    "slices"
    "time"

    agentconfigentities "siem-client/src/modules/agentconfig/domain/entities"
    agentconfigtraits "siem-client/src/modules/agentconfig/infrastructure/traits"
    "siem-client/src/modules/collectors/domain/enums"
    "siem-client/src/modules/collectors/domain/services"
    "siem-client/src/modules/collectors/infrastructure/repositories"
    "siem-client/src/modules/shared/application"
    "siem-client/src/modules/shared/infrastructure/components/uuid"
    "siem-client/src/modules/shared/infrastructure/repositories/spool"
    "siem-client/src/modules/shared/infrastructure/repositories/state"
)

const (
    // por fichero y pasada: un log enorme (o start_at=beginning) se va leyendo en varias pasadas sin cargarlo entero
    MAX_BYTES_PER_FILE_PER_ROUND = 1 << 20
    dateTimeFormat               = "2006-01-02 15:04:05"
)

// CollectFileLinesService una pasada por las carpetas vigiladas: las lineas nuevas de cada fichero pasan al spool
// como eventos log_line, y se guarda hasta donde se leyo cada uno (state/file-offsets.json).
// Primero el spool y despues el estado: si el proceso muere entre medias, en el arranque se reenvian lineas (el
// backend ya las tiene: duplicado) pero nunca se pierden.
type CollectFileLinesService struct {
    agentconfigtraits.AgentConfigLoaderTrait

    uuidGenerator                    *uuid.UuidGenerator
    collectedEventBuilder            *services.CollectedEventBuilder
    eventsBatcher                    *services.EventsBatcher
    watchedFilesReaderDiskRepository *repositories.WatchedFilesReaderDiskRepository
    logLinesReaderDiskRepository     *repositories.LogLinesReaderDiskRepository
    stateReaderJsonRepository        *state.StateReaderJsonRepository
    stateWriterJsonRepository        *state.StateWriterJsonRepository
    spoolWriterDiskRepository        *spool.SpoolWriterDiskRepository

    collectFileLinesDto *CollectFileLinesDto
    offsetsByFilePath   map[string]int64
    initializedWatchers []string
    nextOffsets         map[string]int64
    events              []map[string]any
    filesRead           int
    missingFolders      []string
    unreadableFiles     []string
}

func CollectFileLinesServiceGetInstance() *CollectFileLinesService {
    return &CollectFileLinesService{
        AgentConfigLoaderTrait:           agentconfigtraits.NewAgentConfigLoaderTrait(),
        uuidGenerator:                    uuid.GetInstance(),
        collectedEventBuilder:            services.CollectedEventBuilderGetInstance(),
        eventsBatcher:                    services.EventsBatcherGetInstance(),
        watchedFilesReaderDiskRepository: repositories.WatchedFilesReaderDiskRepositoryGetInstance(),
        logLinesReaderDiskRepository:     repositories.LogLinesReaderDiskRepositoryGetInstance(),
        stateReaderJsonRepository:        state.StateReaderJsonRepositoryGetInstance(),
        stateWriterJsonRepository:        state.StateWriterJsonRepositoryGetInstance(),
        spoolWriterDiskRepository:        spool.SpoolWriterDiskRepositoryGetInstance(),
    }
}

// Invoke error: *AgentConfigException (400, 404) o uno inesperado (disco). Un fichero que no se puede leer no para
// la pasada: va en unreadable_files.
func (s *CollectFileLinesService) Invoke(collectFileLinesDto *CollectFileLinesDto) (*CollectFileLinesResultDto, error) {
    s.collectFileLinesDto = collectFileLinesDto

    agentConfigEntity, err := s.GetValidAgentConfigEntityOrFail()
    if err != nil {
        return nil, err
    }
    if err := s.loadState(); err != nil {
        return nil, err
    }

    occurredAt := time.Now().UTC().Format(dateTimeFormat)
    for _, watcherEntity := range agentConfigEntity.GetWatcherEntities() {
        if err := s.collectWatcher(watcherEntity, occurredAt); err != nil {
            return nil, err
        }
    }

    spooledBatches := 0
    for _, batchEvents := range s.eventsBatcher.GetBatches(s.events, agentConfigEntity.GetMaxEventsPerBatch()) {
        if err := s.spoolWriterDiskRepository.SaveSpoolBatch(s.GetSpoolPath(), s.uuidGenerator.GetUniqueIdWithPrefix("b-"), batchEvents); err != nil {
            return nil, err
        }
        spooledBatches++
    }
    if err := s.saveState(); err != nil {
        return nil, err
    }

    return CollectFileLinesResultDtoFromPrimitives(map[string]any{
        "filesRead":       s.filesRead,
        "linesCollected":  len(s.events),
        "spooledBatches":  spooledBatches,
        "missingFolders":  s.missingFolders,
        "unreadableFiles": s.unreadableFiles,
    }), nil
}

func (s *CollectFileLinesService) loadState() error {
    fileOffsetsState, err := s.stateReaderJsonRepository.GetState(s.GetStatePath(), enums.CollectorsStateFileEnum.FILE_OFFSETS)
    if err != nil {
        return err
    }

    s.offsetsByFilePath = map[string]int64{}
    rawOffsets, _ := fileOffsetsState[enums.CollectorsKeyEnum.OFFSETS].(map[string]any)
    for filePath, rawOffset := range rawOffsets {
        s.offsetsByFilePath[filePath] = int64(application.ToInt(rawOffset))
    }
    s.initializedWatchers = []string{}
    rawInitializedWatchers, _ := fileOffsetsState[enums.CollectorsKeyEnum.INITIALIZED_WATCHERS].([]any)
    for _, rawWatcherName := range rawInitializedWatchers {
        s.initializedWatchers = append(s.initializedWatchers, application.ToString(rawWatcherName))
    }

    s.nextOffsets = map[string]int64{}
    s.events = []map[string]any{}
    s.missingFolders = []string{}
    s.unreadableFiles = []string{}
    return nil
}

// collectWatcher un fichero que ya no esta deja de figurar en el estado (si vuelve, se lee desde el principio).
func (s *CollectFileLinesService) collectWatcher(watcherEntity *agentconfigentities.WatcherEntity, occurredAt string) error {
    if !s.watchedFilesReaderDiskRepository.DoesFolderExist(watcherEntity.GetPath()) {
        s.missingFolders = append(s.missingFolders, watcherEntity.GetPath())
        return nil
    }
    filePaths, err := s.watchedFilesReaderDiskRepository.GetFilePathsByPattern(
        watcherEntity.GetPath(), watcherEntity.GetPattern(), watcherEntity.IsRecursive(),
    )
    if err != nil {
        return err
    }

    isFirstScan := !slices.Contains(s.initializedWatchers, watcherEntity.GetName())
    for _, filePath := range filePaths {
        offset := s.getStartOffset(filePath, isFirstScan && watcherEntity.IsStartAtEnd())
        lines, nextOffset, err := s.logLinesReaderDiskRepository.GetLinesFromOffset(filePath, offset, MAX_BYTES_PER_FILE_PER_ROUND)
        s.nextOffsets[filePath] = nextOffset
        if err != nil {
            s.unreadableFiles = append(s.unreadableFiles, filePath)
            continue
        }
        if len(lines) > 0 {
            s.filesRead++
        }
        for _, line := range lines {
            // match (un log_pattern del panel o un watcher local con match): solo las lineas que casan
            if !watcherEntity.IsLineMatching(line) {
                continue
            }
            s.events = append(s.events, s.collectedEventBuilder.GetWatcherLogLineEvent(watcherEntity, filePath, line, occurredAt))
        }
    }

    if isFirstScan {
        s.initializedWatchers = append(s.initializedWatchers, watcherEntity.GetName())
    }
    return nil
}

// getStartOffset hasta donde se leyo; un fichero nuevo, desde 0, salvo en el primer escaneo con start_at=end (no se
// reenvia el historico que ya habia al instalar el agente).
func (s *CollectFileLinesService) getStartOffset(filePath string, isStartAtEnd bool) int64 {
    if offset, isKnown := s.offsetsByFilePath[filePath]; isKnown {
        return offset
    }
    if !isStartAtEnd {
        return 0
    }
    fileSize, err := s.watchedFilesReaderDiskRepository.GetFileSize(filePath)
    if err != nil {
        return 0
    }
    return fileSize
}

func (s *CollectFileLinesService) saveState() error {
    return s.stateWriterJsonRepository.SaveState(s.GetStatePath(), enums.CollectorsStateFileEnum.FILE_OFFSETS, map[string]any{
        enums.CollectorsKeyEnum.OFFSETS:              s.nextOffsets,
        enums.CollectorsKeyEnum.INITIALIZED_WATCHERS: s.initializedWatchers,
    })
}
