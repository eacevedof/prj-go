package enums

// SmtpSecurityEnum como se protege la conexion SMTP: starttls (587), ssl (465, TLS desde el primer byte) o none (solo
// smtp_custom: Mailpit en local).
var SmtpSecurityEnum = struct {
    STARTTLS string
    SSL      string
    NONE     string
    VALUES   []string
}{
    STARTTLS: "starttls",
    SSL:      "ssl",
    NONE:     "none",
    VALUES:   []string{"starttls", "ssl", "none"},
}
