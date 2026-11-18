package repositories

import (
    "bufio"
    "os"
    "strconv"
    "strings"
)

// PASSWD_PATH los usuarios locales. Los de LDAP/SSSD no estan: se muestran por su uid.
const PASSWD_PATH = "/etc/passwd"

// UsersReaderPasswdRepository uid -> nombre de usuario, para el duenyo de cada socket.
type UsersReaderPasswdRepository struct{}

func UsersReaderPasswdRepositoryGetInstance() *UsersReaderPasswdRepository {
    return &UsersReaderPasswdRepository{}
}

// GetUserNameByUid vacio si no se puede leer el fichero: los sockets se muestran con el uid.
func (r *UsersReaderPasswdRepository) GetUserNameByUid(passwdPath string) map[int]string {
    userNameByUid := map[int]string{}
    file, err := os.Open(passwdPath)
    if err != nil {
        return userNameByUid
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        // nombre:x:uid:gid:...
        fields := strings.Split(scanner.Text(), ":")
        if len(fields) < 3 {
            continue
        }
        uid, err := strconv.Atoi(fields[2])
        if err != nil {
            continue
        }
        if _, isSet := userNameByUid[uid]; !isSet {
            userNameByUid[uid] = fields[0]
        }
    }
    return userNameByUid
}
