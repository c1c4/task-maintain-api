package constants

func GetSupportedPermissionsByUserType() map[string][]string {
	return map[string][]string{
		"Technician": {"create", "update", "get_one", "list_own_tasks"},
		"Manager": {"list", "delete", "notified"},
	}
}