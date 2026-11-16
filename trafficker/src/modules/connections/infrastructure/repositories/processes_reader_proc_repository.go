package repositories

import (
    "os"
    "path/filepath"
    "strconv"
    "strings"

    "trafficker/src/modules/connections/domain/entities"
)

// ProcessesReaderProcRepository de quien es cada socket: recorre <proc>/<pid>/fd y cada enlace "socket:[inodo]" une
// ese inodo con el proceso. Sin root solo se pueden leer los fd de los procesos del propio usuario.
type ProcessesReaderProcRepository struct{}

func ProcessesReaderProcRepositoryGetInstance() *ProcessesReaderProcRepository {
    return &ProcessesReaderProcRepository{}
}

// GetProcessEntitiesByInode inodo -> proceso, y cuantos procesos no se han podido leer (sin permiso). Un socket
// compartido por varios procesos (un fork, los workers de nginx) se queda con el de pid mas bajo.
func (r *ProcessesReaderProcRepository) GetProcessEntitiesByInode(procPath string) (map[uint64]*entities.ProcessEntity, int, error) {
    procEntries, err := os.ReadDir(procPath)
    if err != nil {
        return nil, 0, err
    }

    processEntityByInode := map[uint64]*entities.ProcessEntity{}
    unreadableProcesses := 0
    for _, procEntry := range procEntries {
        pid, err := strconv.Atoi(procEntry.Name())
        if err != nil || !procEntry.IsDir() {
            continue
        }
        pidPath := filepath.Join(procPath, procEntry.Name())
        fdEntries, err := os.ReadDir(filepath.Join(pidPath, "fd"))
        if err != nil {
            // sin permiso, o el proceso ha terminado entre medias
            if os.IsPermission(err) {
                unreadableProcesses++
            }
            continue
        }

        var processEntity *entities.ProcessEntity
        isFdUnreadable := false
        for _, fdEntry := range fdEntries {
            inode, isSocket, err := r.getSocketInode(filepath.Join(pidPath, "fd", fdEntry.Name()))
            if os.IsPermission(err) {
                // se puede listar fd/ pero no leer sus enlaces: falta CAP_SYS_PTRACE (p. ej. root en un contenedor)
                isFdUnreadable = true
                break
            }
            if !isSocket {
                continue
            }
            if previous := processEntityByInode[inode]; previous != nil && previous.GetPid() < pid {
                continue
            }
            if processEntity == nil {
                processEntity = entities.NewProcessEntity(pid, r.getName(pidPath), r.getExe(pidPath))
            }
            processEntityByInode[inode] = processEntity
        }
        if isFdUnreadable {
            unreadableProcesses++
        }
    }
    return processEntityByInode, unreadableProcesses, nil
}

// getSocketInode "socket:[22512]" -> 22512. El error solo se devuelve si es de permisos (el resto: no es un socket,
// o el fd se ha cerrado entre medias).
func (r *ProcessesReaderProcRepository) getSocketInode(fdPath string) (uint64, bool, error) {
    link, err := os.Readlink(fdPath)
    if err != nil {
        if os.IsPermission(err) {
            return 0, false, err
        }
        return 0, false, nil
    }
    if !strings.HasPrefix(link, "socket:[") || !strings.HasSuffix(link, "]") {
        return 0, false, nil
    }
    inode, err := strconv.ParseUint(link[len("socket:["):len(link)-1], 10, 64)
    if err != nil {
        return 0, false, nil
    }
    return inode, true, nil
}

func (r *ProcessesReaderProcRepository) getName(pidPath string) string {
    rawName, err := os.ReadFile(filepath.Join(pidPath, "comm"))
    if err != nil {
        return ""
    }
    return strings.TrimSpace(string(rawName))
}

// getExe vacio si no se puede leer (proceso de otro usuario, o del kernel).
func (r *ProcessesReaderProcRepository) getExe(pidPath string) string {
    exe, err := os.Readlink(filepath.Join(pidPath, "exe"))
    if err != nil {
        return ""
    }
    return strings.TrimSuffix(exe, " (deleted)")
}
