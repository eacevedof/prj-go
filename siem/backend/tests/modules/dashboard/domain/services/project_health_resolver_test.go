package services

import (
    "testing"

    "siem-backend/src/modules/dashboard/domain/services"
    "siem-backend/tests"
)

func TestProjectHealth(t *testing.T) {
    projectHealthResolver := services.ProjectHealthResolverGetInstance()

    tests.AssertEquals(t, "empty", projectHealthResolver.GetProjectHealth(0, 0, 5, 5))
    tests.AssertEquals(t, "critical", projectHealthResolver.GetProjectHealth(2, 0, 0, 0))
    tests.AssertEquals(t, "critical", projectHealthResolver.GetProjectHealth(2, 2, 0, 1))
    tests.AssertEquals(t, "warning", projectHealthResolver.GetProjectHealth(2, 1, 0, 0))
    tests.AssertEquals(t, "warning", projectHealthResolver.GetProjectHealth(2, 2, 3, 0))
    tests.AssertEquals(t, "ok", projectHealthResolver.GetProjectHealth(2, 2, 0, 0))
}
