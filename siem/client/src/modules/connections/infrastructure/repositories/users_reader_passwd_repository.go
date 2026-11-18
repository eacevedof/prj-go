package repositories

import (
    "bufio"
    "errors"
    "io/fs"
    "os"
    "strconv"
    "strings"

    "siem-client/src/modules/connections/domain/enums"
)

// UsersReaderPasswdRepository uid -> nombre de usuario, para el duenyo de cada socket. Los usuarios de LDAP/SSSD no
// estan en /etc/passwd: se muestran por su uid.
type UsersReaderPasswdRepository struct{}

func UsersReaderPasswdRepositoryGetInstance() *UsersReaderPasswdRepository {
    return &UsersReaderPasswdRepository{}
}

// GetUserNameByUid vacio si el fichero no existe (los sockets se muestran con el uid); otro error se devuelve.
func (r *UsersReaderPasswdRepository) GetUserNameByUid(passwdPath string) (map[int]string, error) {
    userNameByUid := map[int]string{}
    file, err := os.Open(passwdPath)
    if errors.Is(err, fs.ErrNotExist) {
        return userNameByUid, nil
    }
    if err != nil {
        return nil, err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        columns := strings.Split(scanner.Text(), enums.PasswdColumnEnum.SEPARATOR)
        if len(columns) < enums.PasswdColumnEnum.MIN_COLUMNS {
            continue
        }
        uid, err := strconv.Atoi(columns[enums.PasswdColumnEnum.UID])
        if err != nil {
            continue
        }
        if _, isSet := userNameByUid[uid]; !isSet {
            userNameByUid[uid] = columns[enums.PasswdColumnEnum.USER_NAME]
        }
    }
    return userNameByUid, scanner.Err()
}
