package services

import (
    "time"

    "siem-backend/src/modules/agents/domain/enums"
)

const (
    dateTimeFormat = "2006-01-02 15:04:05"
    // un agente sigue online mientras no falle mas de 3 latidos seguidos (+ margen por la latencia del envio)
    missedHeartbeatsAllowed = 3
    heartbeatGraceSeconds   = 15
)

// AgentStatusResolver estado de un agente a partir de su ultimo latido. El estado no se guarda: caduca solo.
type AgentStatusResolver struct{}

func AgentStatusResolverGetInstance() *AgentStatusResolver {
    return &AgentStatusResolver{}
}

// GetAgentStatus lastSeenAt y nowDt en UTC "Y-m-d H:i:s". Una fecha ilegible cuenta como offline.
func (r *AgentStatusResolver) GetAgentStatus(isEnabled bool, lastSeenAt string, heartbeatIntervalSeconds int, nowDt string) string {
    if !isEnabled {
        return enums.AgentStatusEnum.DISABLED
    }
    if lastSeenAt == "" {
        return enums.AgentStatusEnum.NEVER_SEEN
    }

    lastSeenTime, err := time.Parse(dateTimeFormat, lastSeenAt)
    if err != nil {
        return enums.AgentStatusEnum.OFFLINE
    }
    nowTime, err := time.Parse(dateTimeFormat, nowDt)
    if err != nil {
        return enums.AgentStatusEnum.OFFLINE
    }

    maxSilenceSeconds := heartbeatIntervalSeconds*missedHeartbeatsAllowed + heartbeatGraceSeconds
    if nowTime.Sub(lastSeenTime) <= time.Duration(maxSilenceSeconds)*time.Second {
        return enums.AgentStatusEnum.ONLINE
    }
    return enums.AgentStatusEnum.OFFLINE
}
