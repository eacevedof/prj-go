package enums

// RequestMetadataKeyEnum claves del map de metadata que HttpRequester.GetMetadata() rellena y AbstractInputDto lee.
var RequestMetadataKeyEnum = struct {
    USER_IP_ADDRESS      string
    USER_BROWSER         string
    USER_BROWSER_VERSION string
    USER_OS              string
    USER_TIME_ZONE       string
    USER_LANGUAGE_CODE   string
}{
    USER_IP_ADDRESS:      "userIpAddress",
    USER_BROWSER:         "userBrowser",
    USER_BROWSER_VERSION: "userBrowserVersion",
    USER_OS:              "userOs",
    USER_TIME_ZONE:       "userTimeZone",
    USER_LANGUAGE_CODE:   "userLanguageCode",
}
