package enums

// UserPolicyEnum politicas "modulo:read|write" y el mapa rol -> politicas (sin tabla de permisos, como mypromos).
// Un owner tiene senders/apitokens/messages, pero solo sobre lo suyo: ese filtro lo pone cada caso de uso.
var UserPolicyEnum = struct {
    USERS_READ      string
    USERS_WRITE     string
    SENDERS_READ    string
    SENDERS_WRITE   string
    APITOKENS_READ  string
    APITOKENS_WRITE string
    MESSAGES_READ   string
    MESSAGES_WRITE  string
    QUEUES_READ     string

    // todo lo de cualquier usuario, no solo lo propio
    ALL_RESOURCES string

    POLICIES_BY_ROLE map[string][]string
}{
    USERS_READ:      "users:read",
    USERS_WRITE:     "users:write",
    SENDERS_READ:    "senders:read",
    SENDERS_WRITE:   "senders:write",
    APITOKENS_READ:  "apitokens:read",
    APITOKENS_WRITE: "apitokens:write",
    MESSAGES_READ:   "messages:read",
    MESSAGES_WRITE:  "messages:write",
    QUEUES_READ:     "queues:read",

    ALL_RESOURCES: "all:resources",

    POLICIES_BY_ROLE: map[string][]string{
        "root": {
            "users:read", "users:write", "senders:read", "senders:write", "apitokens:read", "apitokens:write",
            "messages:read", "messages:write", "queues:read", "all:resources",
        },
        "admin": {
            "users:read", "users:write", "senders:read", "senders:write", "apitokens:read", "apitokens:write",
            "messages:read", "messages:write", "queues:read", "all:resources",
        },
        "owner": {
            "senders:read", "senders:write", "apitokens:read", "apitokens:write", "messages:read", "messages:write",
        },
    },
}
