package entities

import (
    "notifier-backend/src/modules/shared/application"

    "notifier-backend/src/modules/messages/domain/enums"
)

// MessageAttachmentPrimitives lo que se ensena de un adjunto: sin su contenido.
type MessageAttachmentPrimitives struct {
    FileName       string `json:"file_name"`
    AttachmentType string `json:"attachment_type"`
    SourceUrl      string `json:"source_url"`
    MimeType       string `json:"mime_type"`
    SizeBytes      int    `json:"size_bytes"`
}

type MessageAttachmentEntity struct {
    id             int
    messageId      int
    attachmentType string
    fileName       string
    content        string
    sourceUrl      string
    mimeType       string
    sizeBytes      int
}

// MessageAttachmentEntityFromPrimitives desde una fila de app_message_attachments.
func MessageAttachmentEntityFromPrimitives(primitives map[string]any) *MessageAttachmentEntity {
    keys := enums.MessagesKeyEnum
    return &MessageAttachmentEntity{
        id:             application.GetInt(primitives, keys.ID),
        messageId:      application.GetInt(primitives, keys.MESSAGE_ID),
        attachmentType: application.GetTrimmedString(primitives, keys.ATTACHMENT_TYPE),
        fileName:       application.GetTrimmedString(primitives, keys.FILE_NAME),
        // el contenido no se recorta: los espacios del principio o del final son del adjunto
        content:   application.ToString(primitives[keys.CONTENT]),
        sourceUrl: application.GetTrimmedString(primitives, keys.SOURCE_URL),
        mimeType:  application.GetTrimmedString(primitives, keys.MIME_TYPE),
        sizeBytes: application.GetInt(primitives, keys.SIZE_BYTES),
    }
}

func (e *MessageAttachmentEntity) ToPrimitives() MessageAttachmentPrimitives {
    return MessageAttachmentPrimitives{
        FileName:       e.fileName,
        AttachmentType: e.attachmentType,
        SourceUrl:      e.sourceUrl,
        MimeType:       e.mimeType,
        SizeBytes:      e.sizeBytes,
    }
}

func (e *MessageAttachmentEntity) GetId() int {
    return e.id
}

func (e *MessageAttachmentEntity) GetMessageId() int {
    return e.messageId
}

func (e *MessageAttachmentEntity) GetAttachmentType() string {
    return e.attachmentType
}

func (e *MessageAttachmentEntity) IsUrl() bool {
    return e.attachmentType == enums.AttachmentTypeEnum.URL
}

func (e *MessageAttachmentEntity) GetFileName() string {
    return e.fileName
}

func (e *MessageAttachmentEntity) GetContent() string {
    return e.content
}

func (e *MessageAttachmentEntity) GetSourceUrl() string {
    return e.sourceUrl
}

func (e *MessageAttachmentEntity) GetMimeType() string {
    return e.mimeType
}

func (e *MessageAttachmentEntity) GetSizeBytes() int {
    return e.sizeBytes
}
