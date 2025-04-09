package auth

// Role constants
const (
	RoleAdmin    = "admin"
	RoleUser     = "user"
	RoleTrader   = "trader"
	RoleAnalyst  = "analyst"
	RoleReadOnly = "readonly"
)

// Permission constants
const (
	PermissionReadTrades       = "read:trades"
	PermissionCreateTrades     = "create:trades"
	PermissionManageTrades     = "manage:trades"
	PermissionReadStrategies   = "read:strategies"
	PermissionManageStrategies = "manage:strategies"
	PermissionReadSettings     = "read:settings"
	PermissionManageSettings   = "manage:settings"
	PermissionManageUsers      = "manage:users"
	PermissionViewAnalytics    = "view:analytics"
	PermissionManageAPI        = "manage:api"
)

// RolePermissions maps roles to their permissions
var RolePermissions = map[string][]string{
	RoleAdmin: {
		PermissionReadTrades,
		PermissionCreateTrades,
		PermissionManageTrades,
		PermissionReadStrategies,
		PermissionManageStrategies,
		PermissionReadSettings,
		PermissionManageSettings,
		PermissionManageUsers,
		PermissionViewAnalytics,
		PermissionManageAPI,
	},
	RoleTrader: {
		PermissionReadTrades,
		PermissionCreateTrades,
		PermissionManageTrades,
		PermissionReadStrategies,
		PermissionReadSettings,
		PermissionViewAnalytics,
	},
	RoleAnalyst: {
		PermissionReadTrades,
		PermissionReadStrategies,
		PermissionViewAnalytics,
	},
	RoleUser: {
		PermissionReadTrades,
		PermissionReadStrategies,
		PermissionReadSettings,
	},
	RoleReadOnly: {
		PermissionReadTrades,
		PermissionReadStrategies,
	},
}

// HasPermission checks if a role has a specific permission
func HasPermission(role string, permission string) bool {
	permissions, exists := RolePermissions[role]
	if !exists {
		return false
	}
	for _, p := range permissions {
		if p == permission {
			return true
		}
	}
	return false
}

// GetRolePermissions returns all permissions for a role
func GetRolePermissions(role string) []string {
	permissions, exists := RolePermissions[role]
	if !exists {
		return []string{}
	}
	return permissions
}

// ValidateRole checks if a role is valid
func ValidateRole(role string) bool {
	_, exists := RolePermissions[role]
	return exists
}
