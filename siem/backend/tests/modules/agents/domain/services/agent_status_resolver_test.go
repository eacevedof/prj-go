package services

import (
    "testing"

    "siem-backend/src/modules/agents/domain/services"
    "siem-backend/tests"
)

func TestAgentStatus(t *testing.T) {
    agentStatusResolver := services.AgentStatusResolverGetInstance()
    const now = "2026-11-15 10:00:00"

    tests.AssertEquals(t, "disabled", agentStatusResolver.GetAgentStatus(false, "2026-11-15 09:59:59", 60, now))
    tests.AssertEquals(t, "never_seen", agentStatusResolver.GetAgentStatus(true, "", 60, now))
    // 3 latidos de 60 s + 15 s de margen = 195 s
    tests.AssertEquals(t, "online", agentStatusResolver.GetAgentStatus(true, "2026-11-15 09:56:45", 60, now))
    tests.AssertEquals(t, "offline", agentStatusResolver.GetAgentStatus(true, "2026-11-15 09:56:44", 60, now))
    tests.AssertEquals(t, "offline", agentStatusResolver.GetAgentStatus(true, "garbage", 60, now))
}
