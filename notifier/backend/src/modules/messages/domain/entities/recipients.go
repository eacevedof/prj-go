package entities

import "encoding/json"

// Recipients los destinatarios de un mensaje, tal cual se guardan en app_messages.recipients (JSON). Bcc nunca va en
// las cabeceras del email, solo en el sobre SMTP.
type Recipients struct {
    To      []string `json:"to"`
    Cc      []string `json:"cc"`
    Bcc     []string `json:"bcc"`
    ReplyTo string   `json:"reply_to"`
}

// RecipientsFromJson la columna recipients. Un JSON roto da destinatarios vacios (el worker lo dara por fallido).
func RecipientsFromJson(recipientsJson string) Recipients {
    recipients := Recipients{}
    _ = json.Unmarshal([]byte(recipientsJson), &recipients)
    return recipients.getWithoutNilLists()
}

func (r Recipients) ToJson() string {
    recipientsJson, err := json.Marshal(r.getWithoutNilLists())
    if err != nil {
        return "{}"
    }
    return string(recipientsJson)
}

// GetEnvelopeRecipients a quien se entrega de verdad: to + cc + bcc.
func (r Recipients) GetEnvelopeRecipients() []string {
    envelopeRecipients := make([]string, 0, len(r.To)+len(r.Cc)+len(r.Bcc))
    envelopeRecipients = append(envelopeRecipients, r.To...)
    envelopeRecipients = append(envelopeRecipients, r.Cc...)
    return append(envelopeRecipients, r.Bcc...)
}

func (r Recipients) getWithoutNilLists() Recipients {
    if r.To == nil {
        r.To = []string{}
    }
    if r.Cc == nil {
        r.Cc = []string{}
    }
    if r.Bcc == nil {
        r.Bcc = []string{}
    }
    return r
}
