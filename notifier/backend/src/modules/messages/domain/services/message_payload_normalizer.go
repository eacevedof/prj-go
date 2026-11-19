package services

import (
    "encoding/json"
    "net/url"
    "path"
    "strings"

    "notifier-backend/src/modules/shared/application"
    sharedservices "notifier-backend/src/modules/shared/domain/services"

    "notifier-backend/src/modules/messages/domain/enums"
)

const httpsPrefix = "https://"

// AttachmentInput un adjunto tal cual llega, ya normalizado: o contenido (texto) o url.
type AttachmentInput struct {
    FileName   string
    Content    string
    Url        string
    HasContent bool
    HasUrl     bool
}

// MessagePayloadNormalizer lo que el cliente puede mandar de varias formas se deja en una sola (contrato, design §3):
//   - to, cc, bcc: un string ("a@x.com" o "a@x.com, b@x.com") o una lista de strings -> lista sin repetidos.
//   - attachments: lista de {file_name, content|url} o de strings sueltos ("https://..." es url; si no, texto que va
//     como attachment.txt).
//
// Recibe cada campo serializado en JSON (los dto solo llevan primitivos). ok=false = el tipo no es el esperado.
type MessagePayloadNormalizer struct {
    emailAddressValidator *sharedservices.EmailAddressValidator
}

func MessagePayloadNormalizerGetInstance() *MessagePayloadNormalizer {
    return &MessagePayloadNormalizer{emailAddressValidator: sharedservices.EmailAddressValidatorGetInstance()}
}

func (n *MessagePayloadNormalizer) GetEmailList(rawJson string) ([]string, bool) {
    var rawValue any
    if rawJson == "" || json.Unmarshal([]byte(rawJson), &rawValue) != nil || rawValue == nil {
        return []string{}, true
    }

    rawEmails := []string{}
    switch typedValue := rawValue.(type) {
    case string:
        rawEmails = strings.Split(typedValue, ",")
    case []any:
        for _, rawItem := range typedValue {
            emailAddress, isString := rawItem.(string)
            if !isString {
                return []string{}, false
            }
            rawEmails = append(rawEmails, emailAddress)
        }
    default:
        return []string{}, false
    }

    emails := []string{}
    seen := map[string]bool{}
    for _, emailAddress := range rawEmails {
        emailAddress = n.emailAddressValidator.GetNormalized(emailAddress)
        if emailAddress == "" || seen[emailAddress] {
            continue
        }
        seen[emailAddress] = true
        emails = append(emails, emailAddress)
    }
    return emails, true
}

func (n *MessagePayloadNormalizer) GetAttachments(rawJson string) ([]AttachmentInput, bool) {
    var rawValue any
    if rawJson == "" || json.Unmarshal([]byte(rawJson), &rawValue) != nil || rawValue == nil {
        return []AttachmentInput{}, true
    }
    rawItems, isList := rawValue.([]any)
    if !isList {
        return []AttachmentInput{}, false
    }

    attachments := make([]AttachmentInput, 0, len(rawItems))
    for _, rawItem := range rawItems {
        switch typedItem := rawItem.(type) {
        case string:
            attachments = append(attachments, n.getAttachmentFromString(typedItem))
        case map[string]any:
            attachments = append(attachments, n.getAttachmentFromObject(typedItem))
        default:
            return []AttachmentInput{}, false
        }
    }
    return attachments, true
}

func (n *MessagePayloadNormalizer) getAttachmentFromString(rawAttachment string) AttachmentInput {
    trimmedAttachment := strings.TrimSpace(rawAttachment)
    if strings.HasPrefix(strings.ToLower(trimmedAttachment), httpsPrefix) {
        return AttachmentInput{FileName: n.getFileNameFromUrl(trimmedAttachment), Url: trimmedAttachment, HasUrl: true}
    }
    return AttachmentInput{FileName: enums.MessagesLimitEnum.DEFAULT_TEXT_FILE_NAME, Content: rawAttachment, HasContent: true}
}

func (n *MessagePayloadNormalizer) getAttachmentFromObject(rawAttachment map[string]any) AttachmentInput {
    keys := enums.MessagesKeyEnum
    _, hasContent := rawAttachment[keys.CONTENT]
    sourceUrl := application.GetTrimmedString(rawAttachment, keys.URL)
    attachment := AttachmentInput{
        FileName:   application.GetTrimmedString(rawAttachment, keys.FILE_NAME),
        Content:    application.ToString(rawAttachment[keys.CONTENT]),
        Url:        sourceUrl,
        HasContent: hasContent && rawAttachment[keys.CONTENT] != nil,
        HasUrl:     sourceUrl != "",
    }
    if attachment.FileName == "" && attachment.HasUrl {
        attachment.FileName = n.getFileNameFromUrl(sourceUrl)
    }
    if attachment.FileName == "" {
        attachment.FileName = enums.MessagesLimitEnum.DEFAULT_TEXT_FILE_NAME
    }
    return attachment
}

// getFileNameFromUrl https://cdn.x.com/docs/factura.pdf?v=2 -> factura.pdf ("attachment" si no hay nombre).
func (n *MessagePayloadNormalizer) getFileNameFromUrl(rawUrl string) string {
    parsedUrl, err := url.Parse(rawUrl)
    if err != nil {
        return "attachment"
    }
    fileName := path.Base(parsedUrl.Path)
    if fileName == "" || fileName == "." || fileName == "/" {
        return "attachment"
    }
    return fileName
}
