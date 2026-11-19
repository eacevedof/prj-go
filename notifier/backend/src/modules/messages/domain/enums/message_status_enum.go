package enums

// MessageStatusEnum estado de un mensaje en la cola. Mismos valores que el CHECK de app_messages.message_status.
var MessageStatusEnum = struct {
    PENDING   string
    SENDING   string
    SENT      string
    FAILED    string
    CANCELLED string
}{
    PENDING:   "pending",
    SENDING:   "sending",
    SENT:      "sent",
    FAILED:    "failed",
    CANCELLED: "cancelled",
}
