package collectfilelines

import "siem-client/src/modules/shared/application"

// CollectFileLinesPrimitives lo que imprime files-collect (claves en este orden).
type CollectFileLinesPrimitives struct {
    FilesRead       int      `json:"files_read"`
    LinesCollected  int      `json:"lines_collected"`
    SpooledBatches  int      `json:"spooled_batches"`
    MissingFolders  []string `json:"missing_folders"`
    UnreadableFiles []string `json:"unreadable_files"`
}

type CollectFileLinesResultDto struct {
    filesRead       int
    linesCollected  int
    spooledBatches  int
    missingFolders  []string
    unreadableFiles []string
}

func CollectFileLinesResultDtoFromPrimitives(primitives map[string]any) *CollectFileLinesResultDto {
    missingFolders, _ := primitives["missingFolders"].([]string)
    unreadableFiles, _ := primitives["unreadableFiles"].([]string)
    return &CollectFileLinesResultDto{
        filesRead:       application.GetInt(primitives, "filesRead"),
        linesCollected:  application.GetInt(primitives, "linesCollected"),
        spooledBatches:  application.GetInt(primitives, "spooledBatches"),
        missingFolders:  append([]string{}, missingFolders...),
        unreadableFiles: append([]string{}, unreadableFiles...),
    }
}

func (d *CollectFileLinesResultDto) GetLinesCollected() int {
    return d.linesCollected
}

func (d *CollectFileLinesResultDto) GetFilesRead() int {
    return d.filesRead
}

func (d *CollectFileLinesResultDto) GetMissingFolders() []string {
    return d.missingFolders
}

func (d *CollectFileLinesResultDto) GetUnreadableFiles() []string {
    return d.unreadableFiles
}

func (d *CollectFileLinesResultDto) ToArray() CollectFileLinesPrimitives {
    return CollectFileLinesPrimitives{
        FilesRead:       d.filesRead,
        LinesCollected:  d.linesCollected,
        SpooledBatches:  d.spooledBatches,
        MissingFolders:  d.missingFolders,
        UnreadableFiles: d.unreadableFiles,
    }
}
