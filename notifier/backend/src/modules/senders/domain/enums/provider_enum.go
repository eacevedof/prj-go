package enums

// ProviderDefinition lo que trae un proveedor por defecto. Un perfil de envio solo guarda lo que el usuario puede
// cambiar; el resto sale de aqui.
type ProviderDefinition struct {
    Channel      string
    SmtpHost     string
    SmtpPort     int
    SmtpSecurity string
    // los limites del proveedor para una cuenta normal; el perfil puede bajarlos (daily_limit, minute_limit)
    DailyLimit  int
    MinuteLimit int
    // false = existe en el modelo pero aun no tiene MessageSender (fase posterior del roadmap)
    IsImplemented bool
}

// ProviderEnum catalogo de proveedores (sin tabla: anadir uno es anadir un valor aqui y, si hace falta, su
// MessageSender). Mismos codigos que el CHECK de app_senders.provider_code.
var ProviderEnum = struct {
    GMAIL          string
    YAHOO          string
    OUTLOOK        string
    RESEND         string
    SMTP_CUSTOM    string
    TELEGRAM_BOT   string
    WHATSAPP_CLOUD string

    DEFINITIONS map[string]ProviderDefinition
}{
    GMAIL:          "gmail",
    YAHOO:          "yahoo",
    OUTLOOK:        "outlook",
    RESEND:         "resend",
    SMTP_CUSTOM:    "smtp_custom",
    TELEGRAM_BOT:   "telegram_bot",
    WHATSAPP_CLOUD: "whatsapp_cloud",

    DEFINITIONS: map[string]ProviderDefinition{
        // Gmail con contraseña de aplicacion (exige verificacion en dos pasos): 500 destinatarios/dia en una cuenta
        // personal
        "gmail": {Channel: "email", SmtpHost: "smtp.gmail.com", SmtpPort: 587, SmtpSecurity: "starttls", DailyLimit: 500, MinuteLimit: 20, IsImplemented: true},
        "yahoo": {Channel: "email", SmtpHost: "smtp.mail.yahoo.com", SmtpPort: 465, SmtpSecurity: "ssl", DailyLimit: 500, MinuteLimit: 20},
        "outlook": {Channel: "email", SmtpHost: "smtp-mail.outlook.com", SmtpPort: 587, SmtpSecurity: "starttls", DailyLimit: 300, MinuteLimit: 30},
        "resend": {Channel: "email"},
        // cualquier servidor SMTP: host, puerto y seguridad los pone el perfil (en local, Mailpit)
        "smtp_custom":    {Channel: "email", IsImplemented: true},
        "telegram_bot":   {Channel: "telegram"},
        "whatsapp_cloud": {Channel: "whatsapp"},
    },
}
