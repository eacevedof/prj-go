package services

import (
    "net/url"
    "slices"
    "strings"
)

// AttachmentUrlChecker un adjunto por url solo se descarga si es https y su host esta en la lista blanca
// (APP_ATTACHMENT_ALLOWED_HOSTS): asi la API no sirve para que alguien la use de proxy contra la red interna (SSRF).
// Que la ip resuelta no sea privada lo comprueba ademas el descargador, al conectar.
type AttachmentUrlChecker struct {
    allowedHosts []string
}

func AttachmentUrlCheckerFromAllowedHosts(allowedHosts []string) *AttachmentUrlChecker {
    return &AttachmentUrlChecker{allowedHosts: allowedHosts}
}

func (c *AttachmentUrlChecker) IsAllowed(rawUrl string) bool {
    parsedUrl, err := url.Parse(rawUrl)
    if err != nil || parsedUrl.Scheme != "https" || parsedUrl.User != nil {
        return false
    }
    // solo el puerto de https: un host permitido no abre la puerta a sus otros servicios
    if port := parsedUrl.Port(); port != "" && port != "443" {
        return false
    }
    return slices.Contains(c.allowedHosts, strings.ToLower(parsedUrl.Hostname()))
}
