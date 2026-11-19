package runmigrations

import (
    "os"
    "path/filepath"
    "slices"

    "notifier-backend/src/core/boot"
    "notifier-backend/src/core/config"

    "notifier-backend/src/modules/shared/application"

    "notifier-backend/src/modules/devops/domain/enums"
    "notifier-backend/src/modules/devops/infrastructure/repositories"
)

// RunMigrationsService ejecuta, en orden de nombre, los .sql de database/migrations que no esten en la tabla
// migrations. Todos los pendientes de una ejecucion comparten batch (max(batch) + 1).
// ./console run-migrations
type RunMigrationsService struct {
    runMigrationsReaderSqliteRepository *repositories.RunMigrationsReaderSqliteRepository
    runMigrationsWriterSqliteRepository *repositories.RunMigrationsWriterSqliteRepository

    runMigrationsDto   *RunMigrationsDto
    sqlMigrationsFiles []string
    sqlFilesToBeRan    []string
    messages           []string
}

func RunMigrationsServiceGetInstance() *RunMigrationsService {
    return &RunMigrationsService{
        runMigrationsReaderSqliteRepository: repositories.RunMigrationsReaderSqliteRepositoryGetInstance(),
        runMigrationsWriterSqliteRepository: repositories.RunMigrationsWriterSqliteRepositoryGetInstance(),
    }
}

// Invoke error: solo inesperados (un .sql que falla, la BD).
func (s *RunMigrationsService) Invoke(runMigrationsDto *RunMigrationsDto) (*RunMigrationsResultDto, error) {
    s.runMigrationsDto = runMigrationsDto

    if err := s.createMigrationsTableIfDoesNotExist(); err != nil {
        return nil, err
    }

    if err := s.loadMigrationsFiles(); err != nil {
        return nil, err
    }
    if len(s.sqlMigrationsFiles) == 0 {
        s.addMessage(enums.MigrationsMessageEnum.NO_MIGRATIONS_FILES_FOUND)
        return RunMigrationsResultDtoFromPrimitives(map[string]any{"messages": s.messages}), nil
    }

    if err := s.loadFilesToBeRan(); err != nil {
        return nil, err
    }
    if len(s.sqlFilesToBeRan) == 0 {
        s.addMessage(enums.MigrationsMessageEnum.NO_MIGRATIONS_FILES_TO_BE_RAN)
        return RunMigrationsResultDtoFromPrimitives(map[string]any{"messages": s.messages}), nil
    }

    if err := s.runMigrations(); err != nil {
        return nil, err
    }
    return RunMigrationsResultDtoFromPrimitives(map[string]any{"messages": s.messages}), nil
}

func (s *RunMigrationsService) createMigrationsTableIfDoesNotExist() error {
    doesMigrationsTableExist, err := s.runMigrationsReaderSqliteRepository.DoesMigrationsTableExist(enums.MigrationsEnum.MIGRATIONS_TABLE)
    if err != nil {
        return err
    }
    if doesMigrationsTableExist {
        s.addMessage(enums.MigrationsMessageEnum.MIGRATIONS_TABLE_ALREADY_EXISTS)
        return nil
    }
    s.addMessage(enums.MigrationsMessageEnum.CREATING_MIGRATIONS_TABLE)
    return s.runMigrationsWriterSqliteRepository.RunSqlFile(s.getPathSqlFile(enums.MigrationsEnum.CREATE_MIGRATIONS_TABLE_FILE))
}

func (s *RunMigrationsService) loadMigrationsFiles() error {
    dirEntries, err := os.ReadDir(boot.BootPathsEnum.PATH_MIGRATIONS)
    if err != nil {
        return err
    }

    var sqlFiles []string
    for _, dirEntry := range dirEntries {
        fileName := dirEntry.Name()
        if dirEntry.IsDir() || filepath.Ext(fileName) != enums.MigrationsEnum.SQL_FILE_EXTENSION || fileName == enums.MigrationsEnum.CREATE_MIGRATIONS_TABLE_FILE {
            continue
        }
        sqlFiles = append(sqlFiles, fileName)
    }
    slices.Sort(sqlFiles)
    s.sqlMigrationsFiles = sqlFiles
    return nil
}

func (s *RunMigrationsService) loadFilesToBeRan() error {
    ranMigrationsFiles, err := s.runMigrationsReaderSqliteRepository.GetRanMigrationsFiles()
    if err != nil {
        return err
    }
    var ranFilesNames []string
    for _, ranMigrationsFile := range ranMigrationsFiles {
        ranFilesNames = append(ranFilesNames, application.ToString(ranMigrationsFile[enums.MigrationsEnum.MIGRATION_FILE]))
    }

    for _, sqlFile := range s.sqlMigrationsFiles {
        if slices.Contains(ranFilesNames, sqlFile) {
            continue
        }
        s.sqlFilesToBeRan = append(s.sqlFilesToBeRan, sqlFile)
    }
    return nil
}

func (s *RunMigrationsService) runMigrations() error {
    maxBatch, err := s.runMigrationsReaderSqliteRepository.GetMaxBatch()
    if err != nil {
        return err
    }
    batch := maxBatch + 1
    s.addMessage(config.Tr(enums.MigrationsMessageEnum.RUNNING_MIGRATIONS_BATCH_X, batch))
    for _, sqlFile := range s.sqlFilesToBeRan {
        s.addMessage(config.Tr(enums.MigrationsMessageEnum.RUNNING_MIGRATION_X, sqlFile))
        if err := s.runMigrationsWriterSqliteRepository.RunMigrationFileAndSave(s.getPathSqlFile(sqlFile), sqlFile, batch); err != nil {
            return err
        }
        s.addMessage(config.Tr(enums.MigrationsMessageEnum.SAVED_MIGRATION_X, sqlFile))
    }
    return nil
}

func (s *RunMigrationsService) getPathSqlFile(sqlFile string) string {
    return filepath.Join(boot.BootPathsEnum.PATH_MIGRATIONS, sqlFile)
}

// addMessage traza para el ResultDto: el command la imprime y la deja en el log de debug.
func (s *RunMigrationsService) addMessage(message string) {
    s.messages = append(s.messages, message)
}
