package entities

import (
    "siem-backend/src/modules/agents/domain/enums"
    "siem-backend/src/modules/shared/application"
)

const (
    AGENT_UUID_PREFIX                  = "agent-"
    DEFAULT_HEARTBEAT_INTERVAL_SECONDS = 60
)

// AgentPrimitives representacion publica: ni ids internos ni el hash del token.
type AgentPrimitives struct {
    AgentUuid                string `json:"agent_uuid"`
    AgentName                string `json:"agent_name"`
    HeartbeatIntervalSeconds int    `json:"heartbeat_interval_seconds"`
    IsEnabled                bool   `json:"is_enabled"`
    HostName                 string `json:"host_name"`
    HostOs                   string `json:"host_os"`
    HostIp                   string `json:"host_ip"`
    ClientVersion            string `json:"client_version"`
    LastSeenAt               string `json:"last_seen_at"`
    CreatedAt                string `json:"created_at"`
    UpdatedAt                string `json:"updated_at"`
}

// AgentWithApiTokenPrimitives el agente + su api_token en claro: solo en el alta.
type AgentWithApiTokenPrimitives struct {
    AgentPrimitives
    ApiToken string `json:"api_token"`
}

type AgentEntity struct {
    id                       int
    projectId                int
    agentUuid                string
    agentName                string
    apiTokenHash             string
    heartbeatIntervalSeconds int
    isEnabled                bool
    hostName                 string
    hostOs                   string
    hostIp                   string
    clientVersion            string
    lastSeenAt               string
    createdAt                string
    updatedAt                string
}

// AgentEntityFromPrimitives desde una fila de la tabla agents (columnas en snake_case). Las lecturas no traen
// api_token_hash: solo lo informa el alta.
func AgentEntityFromPrimitives(primitives map[string]any) *AgentEntity {
    return &AgentEntity{
        id:                       application.GetInt(primitives, enums.AgentsKeyEnum.ID),
        projectId:                application.GetInt(primitives, enums.AgentsKeyEnum.PROJECT_ID),
        agentUuid:                application.GetTrimmedString(primitives, enums.AgentsKeyEnum.AGENT_UUID),
        agentName:                application.GetTrimmedString(primitives, enums.AgentsKeyEnum.AGENT_NAME),
        apiTokenHash:             application.GetTrimmedString(primitives, enums.AgentsKeyEnum.API_TOKEN_HASH),
        heartbeatIntervalSeconds: application.GetIntOrDefault(primitives, enums.AgentsKeyEnum.HEARTBEAT_INTERVAL_SECONDS, DEFAULT_HEARTBEAT_INTERVAL_SECONDS),
        isEnabled:                application.GetIntOrDefault(primitives, enums.AgentsKeyEnum.IS_ENABLED, 1) == 1,
        hostName:                 application.GetTrimmedString(primitives, enums.AgentsKeyEnum.HOST_NAME),
        hostOs:                   application.GetTrimmedString(primitives, enums.AgentsKeyEnum.HOST_OS),
        hostIp:                   application.GetTrimmedString(primitives, enums.AgentsKeyEnum.HOST_IP),
        clientVersion:            application.GetTrimmedString(primitives, enums.AgentsKeyEnum.CLIENT_VERSION),
        lastSeenAt:               application.GetTrimmedString(primitives, enums.AgentsKeyEnum.LAST_SEEN_AT),
        createdAt:                application.GetTrimmedString(primitives, enums.AgentsKeyEnum.CREATED_AT),
        updatedAt:                application.GetTrimmedString(primitives, enums.AgentsKeyEnum.UPDATED_AT),
    }
}

func (e *AgentEntity) ToPrimitives() AgentPrimitives {
    return AgentPrimitives{
        AgentUuid:                e.agentUuid,
        AgentName:                e.agentName,
        HeartbeatIntervalSeconds: e.heartbeatIntervalSeconds,
        IsEnabled:                e.isEnabled,
        HostName:                 e.hostName,
        HostOs:                   e.hostOs,
        HostIp:                   e.hostIp,
        ClientVersion:            e.clientVersion,
        LastSeenAt:               e.lastSeenAt,
        CreatedAt:                e.createdAt,
        UpdatedAt:                e.updatedAt,
    }
}

// ToPrimitivesOrEmpty una lectura que no encuentra nada da el agente vacio.
func ToPrimitivesOrEmpty(agentEntity *AgentEntity) AgentPrimitives {
    if agentEntity == nil {
        return AgentPrimitives{}
    }
    return agentEntity.ToPrimitives()
}

func (e *AgentEntity) GetId() int {
    return e.id
}

func (e *AgentEntity) GetProjectId() int {
    return e.projectId
}

func (e *AgentEntity) GetAgentUuid() string {
    return e.agentUuid
}

func (e *AgentEntity) GetAgentName() string {
    return e.agentName
}

func (e *AgentEntity) GetApiTokenHash() string {
    return e.apiTokenHash
}

func (e *AgentEntity) GetHeartbeatIntervalSeconds() int {
    return e.heartbeatIntervalSeconds
}

func (e *AgentEntity) IsEnabled() bool {
    return e.isEnabled
}

func (e *AgentEntity) GetHostName() string {
    return e.hostName
}

func (e *AgentEntity) GetHostOs() string {
    return e.hostOs
}

func (e *AgentEntity) GetHostIp() string {
    return e.hostIp
}

func (e *AgentEntity) GetClientVersion() string {
    return e.clientVersion
}

func (e *AgentEntity) GetLastSeenAt() string {
    return e.lastSeenAt
}

func (e *AgentEntity) GetCreatedAt() string {
    return e.createdAt
}

func (e *AgentEntity) GetUpdatedAt() string {
    return e.updatedAt
}
