package enums

// MigrationsEnum la tabla de control de migraciones y los ficheros .sql de database/migrations.
var MigrationsEnum = struct {
    MIGRATIONS_TABLE             string
    CREATE_MIGRATIONS_TABLE_FILE string
    SQL_FILE_EXTENSION           string
    MIGRATION_FILE               string
}{
    MIGRATIONS_TABLE:             "migrations",
    CREATE_MIGRATIONS_TABLE_FILE: "0000-create-migrations-table.sql",
    SQL_FILE_EXTENSION:           ".sql",
    MIGRATION_FILE:               "migration_file",
}
