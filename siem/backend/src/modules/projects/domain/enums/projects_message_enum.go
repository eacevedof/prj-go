package enums

// ProjectsMessageEnum mensajes de error de projects ({0}, {1}... los rellena config.Tr).
var ProjectsMessageEnum = struct {
    PROJECT_NAME_X_ALREADY_EXISTS          string
    PROJECT_X_NOT_FOUND                    string
    PROJECT_NAME_LENGTH_MUST_BE_3_TO_100   string
    PROJECT_NAME_HAS_INVALID_CHARS         string
    PROJECT_DESCRIPTION_MAX_LENGTH_IS_500  string
}{
    PROJECT_NAME_X_ALREADY_EXISTS:         "Project name {0} already exists",
    PROJECT_X_NOT_FOUND:                   "Project {0} not found",
    PROJECT_NAME_LENGTH_MUST_BE_3_TO_100:  "project_name: length must be between 3 and 100",
    PROJECT_NAME_HAS_INVALID_CHARS:        "project_name: only letters, digits, spaces and . _ - are allowed",
    PROJECT_DESCRIPTION_MAX_LENGTH_IS_500: "project_description: max length is 500",
}
