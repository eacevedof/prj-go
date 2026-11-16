package enums

// ProjectFieldsLimitEnum limites de los campos de un proyecto (ProjectFieldsValidator). Si cambian, cambiar tambien
// los mensajes de ProjectsMessageEnum y los limites del front.
var ProjectFieldsLimitEnum = struct {
    PROJECT_NAME_MIN_LENGTH        int
    PROJECT_NAME_MAX_LENGTH        int
    PROJECT_DESCRIPTION_MAX_LENGTH int
}{
    PROJECT_NAME_MIN_LENGTH:        3,
    PROJECT_NAME_MAX_LENGTH:        100,
    PROJECT_DESCRIPTION_MAX_LENGTH: 500,
}
