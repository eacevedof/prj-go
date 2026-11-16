package enums

// MigrationsMessageEnum traza de RunMigrationsService ({0}... los rellena config.Tr).
var MigrationsMessageEnum = struct {
    NO_MIGRATIONS_FILES_FOUND       string
    NO_MIGRATIONS_FILES_TO_BE_RAN   string
    MIGRATIONS_TABLE_ALREADY_EXISTS string
    CREATING_MIGRATIONS_TABLE       string
    RUNNING_MIGRATIONS_BATCH_X      string
    RUNNING_MIGRATION_X             string
    SAVED_MIGRATION_X               string
}{
    NO_MIGRATIONS_FILES_FOUND:       "No migrations files found",
    NO_MIGRATIONS_FILES_TO_BE_RAN:   "No migrations files to be ran",
    MIGRATIONS_TABLE_ALREADY_EXISTS: "migrations table wont be created, already exists",
    CREATING_MIGRATIONS_TABLE:       "creating migrations table",
    RUNNING_MIGRATIONS_BATCH_X:      "running migrations batch: {0}",
    RUNNING_MIGRATION_X:             "running migration: {0}",
    SAVED_MIGRATION_X:               "saved migration: {0}",
}
