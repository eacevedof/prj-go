package services

import (
    "regexp"
    "unicode/utf8"

    "siem-backend/src/modules/projects/domain/enums"
)

var projectNamePattern = regexp.MustCompile(`^[\p{L}\p{N} ._-]+$`)

// ProjectFieldsValidator reglas de los campos de un proyecto. Devuelve la pila de errores (vacia = todo bien):
// quien la usa decide con que excepcion escalar.
type ProjectFieldsValidator struct{}

func ProjectFieldsValidatorGetInstance() *ProjectFieldsValidator {
    return &ProjectFieldsValidator{}
}

func (v *ProjectFieldsValidator) GetErrors(projectName string, projectDescription string) []string {
    errors := []string{}
    projectNameLength := utf8.RuneCountInString(projectName)
    if projectNameLength < enums.ProjectFieldsLimitEnum.PROJECT_NAME_MIN_LENGTH || projectNameLength > enums.ProjectFieldsLimitEnum.PROJECT_NAME_MAX_LENGTH {
        errors = append(errors, enums.ProjectsMessageEnum.PROJECT_NAME_LENGTH_MUST_BE_3_TO_100)
    } else if !projectNamePattern.MatchString(projectName) {
        errors = append(errors, enums.ProjectsMessageEnum.PROJECT_NAME_HAS_INVALID_CHARS)
    }
    if utf8.RuneCountInString(projectDescription) > enums.ProjectFieldsLimitEnum.PROJECT_DESCRIPTION_MAX_LENGTH {
        errors = append(errors, enums.ProjectsMessageEnum.PROJECT_DESCRIPTION_MAX_LENGTH_IS_500)
    }
    return errors
}
