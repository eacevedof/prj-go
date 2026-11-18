package repositories

import (
    "bufio"
    "errors"
    "io/fs"
    "os"
    "path/filepath"
    "runtime"

    "siem-client/src/modules/connections/domain/entities"
    "siem-client/src/modules/connections/domain/enums"
    "siem-client/src/modules/connections/domain/services"
)

// PROC_PATH donde estan los sockets (net/) y los procesos (<pid>/fd) del host.
const PROC_PATH = "/proc"

// SocketsReaderProcRepository la foto de todos los sockets TCP y UDP (IPv4 e IPv6) del espacio de red del proceso,
// leida de <proc>/net/{tcp,tcp6,udp,udp6}. Solo stdlib: no hace falta ss ni netstat. Solo Linux: en el resto el
// modulo compila pero no muestrea (isProcNetSupported, por build tags).
type SocketsReaderProcRepository struct {
    procNetLineParser *services.ProcNetLineParser
}

func SocketsReaderProcRepositoryGetInstance() *SocketsReaderProcRepository {
    return &SocketsReaderProcRepository{procNetLineParser: services.ProcNetLineParserGetInstance()}
}

// IsSupported false fuera de Linux (Windows y macOS no tienen /proc/net).
func (r *SocketsReaderProcRepository) IsSupported() bool {
    return isProcNetSupported
}

func (r *SocketsReaderProcRepository) GetOsName() string {
    return runtime.GOOS
}

// GetAllSocketEntities error si no existe <proc>/net/tcp (no es Linux o /proc no esta montado). Sin IPv6 en el
// kernel faltan tcp6/udp6: se sigue sin ellos.
func (r *SocketsReaderProcRepository) GetAllSocketEntities(procPath string) ([]*entities.SocketEntity, error) {
    socketEntities := []*entities.SocketEntity{}
    procNetFiles := []string{
        enums.ProcNetFileEnum.TCP, enums.ProcNetFileEnum.TCP6, enums.ProcNetFileEnum.UDP, enums.ProcNetFileEnum.UDP6,
    }
    for _, procNetFile := range procNetFiles {
        fileSocketEntities, err := r.getFileSocketEntities(procPath, procNetFile)
        if errors.Is(err, fs.ErrNotExist) && procNetFile != enums.ProcNetFileEnum.TCP {
            continue
        }
        if err != nil {
            return nil, err
        }
        socketEntities = append(socketEntities, fileSocketEntities...)
    }
    return socketEntities, nil
}

func (r *SocketsReaderProcRepository) getFileSocketEntities(procPath string, procNetFile string) ([]*entities.SocketEntity, error) {
    file, err := os.Open(filepath.Join(procPath, "net", procNetFile))
    if err != nil {
        return nil, err
    }
    defer file.Close()

    socketEntities := []*entities.SocketEntity{}
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        if socketEntity := r.procNetLineParser.GetSocketEntityOrNil(procNetFile, scanner.Text()); socketEntity != nil {
            socketEntities = append(socketEntities, socketEntity)
        }
    }
    return socketEntities, scanner.Err()
}
