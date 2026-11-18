package entities

import (
    "siem-client/src/modules/shared/application"

    "siem-client/src/modules/connections/domain/enums"
)

// ProcessEntity el proceso que tiene abierto un socket: pid, nombre corto (/proc/<pid>/comm) y ejecutable
// (/proc/<pid>/exe; vacio si no se puede leer).
type ProcessEntity struct {
    pid  int
    name string
    exe  string
}

// ProcessEntityFromPrimitives claves de ConnectionKeyEnum: pid, process (el nombre) y exe.
func ProcessEntityFromPrimitives(primitives map[string]any) *ProcessEntity {
    return &ProcessEntity{
        pid:  application.GetInt(primitives, enums.ConnectionKeyEnum.PID),
        name: application.GetTrimmedString(primitives, enums.ConnectionKeyEnum.PROCESS),
        exe:  application.GetTrimmedString(primitives, enums.ConnectionKeyEnum.EXE),
    }
}

// UnknownProcessEntity el socket existe pero no se sabe de quien es (sin permiso para leer los fd de su proceso).
func UnknownProcessEntity() *ProcessEntity {
    return &ProcessEntity{pid: 0, name: enums.UNKNOWN_PROCESS, exe: ""}
}

func (e *ProcessEntity) GetPid() int {
    return e.pid
}

func (e *ProcessEntity) GetName() string {
    return e.name
}

func (e *ProcessEntity) GetExe() string {
    return e.exe
}

func (e *ProcessEntity) IsUnknown() bool {
    return e.pid == 0
}
