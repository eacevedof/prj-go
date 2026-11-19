package enums

// UserRoleEnum roles del panel (los del CHECK de app_users.user_role). Que puede hacer cada uno: UserPolicyEnum.
var UserRoleEnum = struct {
    ROOT   string
    ADMIN  string
    OWNER  string
    VALUES []string
}{
    ROOT:   "root",
    ADMIN:  "admin",
    OWNER:  "owner",
    VALUES: []string{"root", "admin", "owner"},
}
