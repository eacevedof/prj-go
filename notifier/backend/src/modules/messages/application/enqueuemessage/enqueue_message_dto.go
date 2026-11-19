package enqueuemessage

import "notifier-backend/src/modules/shared/application"

// EnqueueMessageDto un mensaje de POST /api/v1/messages. to, cc, bcc y attachments llegan serializados en JSON tal cual
// los mando el cliente (pueden ser string o lista): los normaliza el servicio. Lo auth* sale del token de envio.
type EnqueueMessageDto struct {
    application.AbstractInputDto
    authApiTokenId    int
    authSenderId      int
    authChannel       string
    authDefaultToJson string
    toJson            string
    ccJson            string
    bccJson           string
    replyTo           string
    subject           string
    body              string
    bodyFormat        string
    attachmentsJson   string
    idempotencyKey    string
}

func EnqueueMessageDtoFromPrimitives(primitives map[string]any) *EnqueueMessageDto {
    return &EnqueueMessageDto{
        AbstractInputDto:  application.NewAbstractInputDto(primitives),
        authApiTokenId:    application.GetInt(primitives, "authApiTokenId"),
        authSenderId:      application.GetInt(primitives, "authSenderId"),
        authChannel:       application.GetTrimmedString(primitives, "authChannel"),
        authDefaultToJson: application.GetTrimmedString(primitives, "authDefaultTo"),
        toJson:            application.GetTrimmedString(primitives, "toJson"),
        ccJson:            application.GetTrimmedString(primitives, "ccJson"),
        bccJson:           application.GetTrimmedString(primitives, "bccJson"),
        replyTo:           application.GetTrimmedString(primitives, "replyTo"),
        subject:           application.GetTrimmedString(primitives, "subject"),
        // el cuerpo no se recorta: sus espacios y saltos de linea son del mensaje
        body:            application.ToString(primitives["body"]),
        bodyFormat:      application.GetTrimmedString(primitives, "bodyFormat"),
        attachmentsJson: application.GetTrimmedString(primitives, "attachmentsJson"),
        idempotencyKey:  application.GetTrimmedString(primitives, "idempotencyKey"),
    }
}

func (d *EnqueueMessageDto) GetAuthApiTokenId() int {
    return d.authApiTokenId
}

func (d *EnqueueMessageDto) GetAuthSenderId() int {
    return d.authSenderId
}

func (d *EnqueueMessageDto) GetAuthChannel() string {
    return d.authChannel
}

func (d *EnqueueMessageDto) GetAuthDefaultToJson() string {
    return d.authDefaultToJson
}

func (d *EnqueueMessageDto) GetToJson() string {
    return d.toJson
}

func (d *EnqueueMessageDto) GetCcJson() string {
    return d.ccJson
}

func (d *EnqueueMessageDto) GetBccJson() string {
    return d.bccJson
}

func (d *EnqueueMessageDto) GetReplyTo() string {
    return d.replyTo
}

func (d *EnqueueMessageDto) GetSubject() string {
    return d.subject
}

func (d *EnqueueMessageDto) GetBody() string {
    return d.body
}

func (d *EnqueueMessageDto) GetBodyFormat() string {
    return d.bodyFormat
}

func (d *EnqueueMessageDto) GetAttachmentsJson() string {
    return d.attachmentsJson
}

func (d *EnqueueMessageDto) GetIdempotencyKey() string {
    return d.idempotencyKey
}
