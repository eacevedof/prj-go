package deleteproject

import "siem-backend/src/modules/shared/application"

// DeleteProjectPrimitives lo que devuelve la API al borrar.
type DeleteProjectPrimitives struct {
    ProjectUuid string `json:"project_uuid"`
}

type DeleteProjectResultDto struct {
    projectUuid string
}

func DeleteProjectResultDtoFromPrimitives(primitives map[string]any) *DeleteProjectResultDto {
    return &DeleteProjectResultDto{projectUuid: application.GetTrimmedString(primitives, "projectUuid")}
}

func (d *DeleteProjectResultDto) ToArray() DeleteProjectPrimitives {
    return DeleteProjectPrimitives{ProjectUuid: d.projectUuid}
}
