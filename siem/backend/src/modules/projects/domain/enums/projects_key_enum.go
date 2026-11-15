package enums

// ProjectsKeyEnum claves de la API y de consola de projects, mas las columnas con las que se monta ProjectEntity.
var ProjectsKeyEnum = struct {
    PROJECT_UUID        string
    PROJECT_NAME        string
    PROJECT_DESCRIPTION string
    CREATED_AT          string
    UPDATED_AT          string

    // columna interna de BD: nunca sale por la API
    ID string
}{
    PROJECT_UUID:        "project_uuid",
    PROJECT_NAME:        "project_name",
    PROJECT_DESCRIPTION: "project_description",
    CREATED_AT:          "created_at",
    UPDATED_AT:          "updated_at",

    ID: "id",
}
