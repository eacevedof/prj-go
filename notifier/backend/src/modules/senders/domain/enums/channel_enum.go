package enums

// ChannelEnum canales de salida: cada uno es una cola (app_messages.channel) con su worker. Mismos valores que los
// CHECK de app_senders.channel y app_messages.channel.
var ChannelEnum = struct {
    EMAIL    string
    WHATSAPP string
    TELEGRAM string
    VALUES   []string
}{
    EMAIL:    "email",
    WHATSAPP: "whatsapp",
    TELEGRAM: "telegram",
    VALUES:   []string{"email", "whatsapp", "telegram"},
}
