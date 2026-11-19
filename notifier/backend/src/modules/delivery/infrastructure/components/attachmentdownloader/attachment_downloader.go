package attachmentdownloader

import (
    "context"
    "errors"
    "fmt"
    "io"
    "mime"
    "net"
    "net/http"
    "path"
    "syscall"
    "time"

    "notifier-backend/src/core/config"

    messagesservices "notifier-backend/src/modules/messages/domain/services"

    "notifier-backend/src/modules/delivery/domain/enums"
    "notifier-backend/src/modules/delivery/domain/exceptions"
)

// cgnat 100.64.0.0/10: no es IsPrivate() pero tampoco es internet
var cgnatNetwork = &net.IPNet{IP: net.IPv4(100, 64, 0, 0), Mask: net.CIDRMask(10, 32)}

var errPrivateAddress = errors.New("private address")

// AttachmentDownloader descarga un adjunto por url para meterlo en el email. Defensas contra SSRF: la url (y cada
// redireccion) pasa por AttachmentUrlChecker (https + host de la lista blanca), y al conectar se rechaza cualquier ip
// que no sea publica (loopback, privada, link-local, cgnat...), resuelva el DNS lo que resuelva. Con tope de tamano y
// de tiempo.
type AttachmentDownloader struct {
    attachmentUrlChecker *messagesservices.AttachmentUrlChecker
    maxBytes             int
    httpClient           *http.Client
}

func AttachmentDownloaderGetInstance(attachmentUrlChecker *messagesservices.AttachmentUrlChecker, maxBytes int, timeoutSeconds int) *AttachmentDownloader {
    attachmentDownloader := &AttachmentDownloader{attachmentUrlChecker: attachmentUrlChecker, maxBytes: maxBytes}
    dialer := &net.Dialer{Timeout: 10 * time.Second, Control: attachmentDownloader.failIfPrivateAddress}
    attachmentDownloader.httpClient = &http.Client{
        Timeout: time.Duration(timeoutSeconds) * time.Second,
        Transport: &http.Transport{
            // sin proxy del entorno: la comprobacion de ip tiene que ser contra el destino real
            Proxy:               nil,
            DialContext:         dialer.DialContext,
            TLSHandshakeTimeout: 10 * time.Second,
        },
        CheckRedirect: attachmentDownloader.checkRedirect,
    }
    return attachmentDownloader
}

// GetContent el contenido y su tipo (el de la cabecera o, si no viene, el de la extension). position es el numero del
// adjunto en el mensaje, para el error.
func (d *AttachmentDownloader) GetContent(sourceUrl string, fileName string, position int) ([]byte, string, error) {
    if !d.attachmentUrlChecker.IsAllowed(sourceUrl) {
        return nil, "", exceptions.PermanentCustom(config.Tr(enums.DeliveryMessageEnum.ATTACHMENT_X_URL_NOT_ALLOWED, position))
    }

    request, err := http.NewRequestWithContext(context.Background(), http.MethodGet, sourceUrl, nil)
    if err != nil {
        return nil, "", exceptions.PermanentCustom(config.Tr(enums.DeliveryMessageEnum.ATTACHMENT_X_URL_NOT_ALLOWED, position))
    }
    request.Header.Set("User-Agent", "notifier/1.0")
    response, err := d.httpClient.Do(request)
    if err != nil {
        return nil, "", d.getDownloadException(err, position)
    }
    defer response.Body.Close()

    if response.StatusCode != http.StatusOK {
        message := config.Tr(enums.DeliveryMessageEnum.ATTACHMENT_X_HTTP_STATUS_Y, position, response.StatusCode)
        // un 4xx no se arregla reintentando; un 5xx puede que si
        if response.StatusCode >= 400 && response.StatusCode < 500 {
            return nil, "", exceptions.PermanentCustom(message)
        }
        return nil, "", exceptions.TemporaryCustom(message)
    }

    content, err := io.ReadAll(io.LimitReader(response.Body, int64(d.maxBytes)+1))
    if err != nil {
        return nil, "", d.getDownloadException(err, position)
    }
    if len(content) > d.maxBytes {
        return nil, "", exceptions.PermanentCustom(config.Tr(enums.DeliveryMessageEnum.ATTACHMENT_X_TOO_BIG, position, d.maxBytes))
    }
    return content, d.getMimeType(response.Header.Get("Content-Type"), fileName), nil
}

func (d *AttachmentDownloader) checkRedirect(request *http.Request, previousRequests []*http.Request) error {
    if len(previousRequests) > enums.DeliveryLimitEnum.ATTACHMENT_MAX_REDIRECTS || !d.attachmentUrlChecker.IsAllowed(request.URL.String()) {
        return http.ErrUseLastResponse
    }
    return nil
}

// failIfPrivateAddress el Control del dialer: se ejecuta con la ip ya resuelta, justo antes de conectar.
func (d *AttachmentDownloader) failIfPrivateAddress(network string, address string, rawConn syscall.RawConn) error {
    host, _, err := net.SplitHostPort(address)
    if err != nil {
        return err
    }
    ip := net.ParseIP(host)
    if ip == nil || !IsPublicIp(ip) {
        return errPrivateAddress
    }
    return nil
}

// IsPublicIp false para loopback, privadas, link-local, multicast, sin especificar y cgnat.
func IsPublicIp(ip net.IP) bool {
    return !(ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() ||
        ip.IsInterfaceLocalMulticast() || ip.IsMulticast() || ip.IsUnspecified() || cgnatNetwork.Contains(ip))
}

func (d *AttachmentDownloader) getDownloadException(err error, position int) error {
    if errors.Is(err, errPrivateAddress) {
        return exceptions.PermanentCustom(config.Tr(enums.DeliveryMessageEnum.ATTACHMENT_X_HOST_IS_PRIVATE, position))
    }
    return exceptions.TemporaryCustom(config.Tr(enums.DeliveryMessageEnum.ATTACHMENT_X_DOWNLOAD_FAILED, position, fmt.Sprint(err)))
}

// getMimeType "application/pdf; charset=..." -> "application/pdf"; sin cabecera, por la extension del nombre.
func (d *AttachmentDownloader) getMimeType(contentType string, fileName string) string {
    if mediaType, _, err := mime.ParseMediaType(contentType); err == nil && mediaType != "" {
        return mediaType
    }
    if mimeType := mime.TypeByExtension(path.Ext(fileName)); mimeType != "" {
        mediaType, _, _ := mime.ParseMediaType(mimeType)
        return mediaType
    }
    return "application/octet-stream"
}
