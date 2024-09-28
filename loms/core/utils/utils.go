package utils

func Contains(statuses []string, status string) bool {
	for _, s := range statuses {
		if s == status {
			return true
		}
	}
	return false
}
