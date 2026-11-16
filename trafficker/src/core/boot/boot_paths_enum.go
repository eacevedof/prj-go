package boot

import (
    "os"
    "path/filepath"
)

// BootPathsEnum rutas absolutas de la app. Dos modos:
//   - REPO (hay go.mod subiendo desde el cwd): raiz = esa carpeta; todo lo generado en storage/ y la config en config/.
//   - BINARIO SUELTO (no hay go.mod): raiz = carpeta del ejecutable; logs/, spool/ y state/ junto a el, para
//     desplegar el agente en un host copiando el binario + .env.
//
// PATH_SPOOL: eventos recogidos y aun no entregados al backend. PATH_STATE: la ultima foto de los sockets.
var BootPathsEnum = func() struct {
    PATH_ROOT        string
    PATH_SRC         string
    PATH_LOGS        string
    PATH_SPOOL       string
    PATH_STATE       string
    IS_PACKAGED      bool
} {
    pathRoot, isPackaged := getPathRoot()
    pathStorage := filepath.Join(pathRoot, "storage")
    if isPackaged {
        pathStorage = pathRoot
    }
    return struct {
        PATH_ROOT        string
        PATH_SRC         string
        PATH_LOGS        string
        PATH_SPOOL       string
        PATH_STATE       string
        IS_PACKAGED      bool
    }{
        PATH_ROOT:        pathRoot,
        PATH_SRC:         filepath.Join(pathRoot, "src"),
        PATH_LOGS:        filepath.Join(pathStorage, "logs"),
        PATH_SPOOL:       filepath.Join(pathStorage, "spool"),
        PATH_STATE:       filepath.Join(pathStorage, "state"),
        IS_PACKAGED:      isPackaged,
    }
}()

// getPathRoot (raiz, empaquetado): la primera carpeta con go.mod subiendo desde el cwd; si no hay, la del ejecutable.
func getPathRoot() (string, bool) {
    if cwd, err := os.Getwd(); err == nil {
        for dir := cwd; ; dir = filepath.Dir(dir) {
            if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
                return dir, false
            }
            if filepath.Dir(dir) == dir {
                break
            }
        }
    }
    if executable, err := os.Executable(); err == nil {
        return filepath.Dir(executable), true
    }
    return ".", true
}
