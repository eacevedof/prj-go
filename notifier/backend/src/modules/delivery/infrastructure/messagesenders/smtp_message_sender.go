package messagesenders

import (
    "crypto/tls"
    "net"
    "net/smtp"
    "strconv"
    "time"

    "notifier-backend/src/core/config"

    sendersenums "notifier-backend/src/modules/senders/domain/enums"

    "notifier-backend/src/modules/delivery/domain/entities"
    "notifier-backend/src/modules/delivery/domain/enums"
    "notifier-backend/src/modules/delivery/domain/exceptions"
    "notifier-backend/src/modules/delivery/domain/services"
)

// SmtpMessageSender entrega un email por SMTP: Gmail, Yahoo, Outlook.com o cualquier servidor (smtp_custom, Mailpit
// en local). starttls (587) sube a TLS tras el EHLO y no sigue si el servidor no lo ofrece; ssl (465) es TLS desde el
// primer byte; none solo sin credenciales (net/smtp no manda PLAIN en claro salvo a localhost).
// Toda la conversacion tiene un tope de SMTP_TIMEOUT_SECONDS.
type SmtpMessageSender struct {
    mimeMessageBuilder *services.MimeMessageBuilder
}

func SmtpMessageSenderGetInstance() *SmtpMessageSender {
    return &SmtpMessageSender{mimeMessageBuilder: services.MimeMessageBuilderGetInstance()}
}

func (s *SmtpMessageSender) Send(outgoingEmail entities.OutgoingEmail, smtpSettings entities.SmtpSettings) (string, error) {
    message, err := s.mimeMessageBuilder.GetMessage(outgoingEmail)
    if err != nil {
        return "", exceptions.PermanentCustom(err.Error())
    }

    smtpClient, err := s.getConnectedClient(smtpSettings)
    if err != nil {
        return "", err
    }
    defer smtpClient.Close()

    if err := s.authenticateIfNeeded(smtpClient, smtpSettings); err != nil {
        return "", err
    }
    if err := smtpClient.Mail(outgoingEmail.FromEmail); err != nil {
        return "", err
    }
    for _, recipient := range outgoingEmail.GetEnvelopeRecipients() {
        if err := smtpClient.Rcpt(recipient); err != nil {
            return "", err
        }
    }
    dataWriter, err := smtpClient.Data()
    if err != nil {
        return "", err
    }
    if _, err := dataWriter.Write(message); err != nil {
        return "", err
    }
    // el Close del DATA es cuando el servidor acepta (250) o rechaza el mensaje
    if err := dataWriter.Close(); err != nil {
        return "", err
    }
    _ = smtpClient.Quit()
    return config.Tr(enums.DeliveryMessageEnum.SENT_WITH_MESSAGE_ID_X, outgoingEmail.MessageId), nil
}

func (s *SmtpMessageSender) getConnectedClient(smtpSettings entities.SmtpSettings) (*smtp.Client, error) {
    address := net.JoinHostPort(smtpSettings.Host, strconv.Itoa(smtpSettings.Port))
    timeout := time.Duration(enums.DeliveryLimitEnum.SMTP_TIMEOUT_SECONDS) * time.Second
    tlsConfig := &tls.Config{ServerName: smtpSettings.Host, MinVersion: tls.VersionTLS12}

    var connection net.Conn
    var err error
    if smtpSettings.Security == sendersenums.SmtpSecurityEnum.SSL {
        connection, err = tls.DialWithDialer(&net.Dialer{Timeout: timeout}, "tcp", address, tlsConfig)
    } else {
        connection, err = net.DialTimeout("tcp", address, timeout)
    }
    if err != nil {
        return nil, err
    }
    _ = connection.SetDeadline(time.Now().Add(timeout))

    smtpClient, err := smtp.NewClient(connection, smtpSettings.Host)
    if err != nil {
        _ = connection.Close()
        return nil, err
    }
    if smtpSettings.Security != sendersenums.SmtpSecurityEnum.STARTTLS {
        return smtpClient, nil
    }
    if isStartTlsOffered, _ := smtpClient.Extension("STARTTLS"); !isStartTlsOffered {
        _ = smtpClient.Close()
        return nil, exceptions.PermanentCustom(config.Tr(enums.DeliveryMessageEnum.SMTP_SERVER_HAS_NO_STARTTLS))
    }
    if err := smtpClient.StartTLS(tlsConfig); err != nil {
        _ = smtpClient.Close()
        return nil, err
    }
    return smtpClient, nil
}

func (s *SmtpMessageSender) authenticateIfNeeded(smtpClient *smtp.Client, smtpSettings entities.SmtpSettings) error {
    if smtpSettings.Password == "" {
        return nil
    }
    if smtpSettings.Security == sendersenums.SmtpSecurityEnum.NONE {
        return exceptions.PermanentCustom(config.Tr(enums.DeliveryMessageEnum.SMTP_AUTH_NEEDS_TLS))
    }
    return smtpClient.Auth(smtp.PlainAuth("", smtpSettings.Username, smtpSettings.Password, smtpSettings.Host))
}
