package services

import (
    "encoding/json"
    "testing"
)

// toJsonPrimitives pasa por JSON, como la cache en disco (los bool siguen siendo bool, los numeros float64).
func toJsonPrimitives(t *testing.T, primitives map[string]any) map[string]any {
    t.Helper()
    rawJson, err := json.Marshal(primitives)
    if err != nil {
        t.Fatal(err)
    }
    restored := map[string]any{}
    if err := json.Unmarshal(rawJson, &restored); err != nil {
        t.Fatal(err)
    }
    return restored
}
