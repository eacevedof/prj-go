package services

import (
    "slices"

    "notifier-backend/src/modules/users/domain/enums"
)

// UserPolicyChecker que politicas tiene un rol (UserPolicyEnum.POLICIES_BY_ROLE). Un rol desconocido no tiene ninguna.
type UserPolicyChecker struct{}

func UserPolicyCheckerGetInstance() *UserPolicyChecker {
    return &UserPolicyChecker{}
}

func (c *UserPolicyChecker) HasPolicy(userRole string, policy string) bool {
    return slices.Contains(enums.UserPolicyEnum.POLICIES_BY_ROLE[userRole], policy)
}

// CanSeeAllResources root y admin ven y tocan lo de todos; un owner, solo lo suyo.
func (c *UserPolicyChecker) CanSeeAllResources(userRole string) bool {
    return c.HasPolicy(userRole, enums.UserPolicyEnum.ALL_RESOURCES)
}
