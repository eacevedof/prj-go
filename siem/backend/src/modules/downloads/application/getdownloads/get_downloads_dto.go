package getdownloads

import "siem-backend/src/modules/shared/application"

// GetDownloadsDto GET /api/v1/downloads: sin campos propios (la carpeta y la url base salen de la configuracion).
type GetDownloadsDto struct {
    application.AbstractInputDto
}

func GetDownloadsDtoFromPrimitives(primitives map[string]any) *GetDownloadsDto {
    return &GetDownloadsDto{AbstractInputDto: application.NewAbstractInputDto(primitives)}
}
