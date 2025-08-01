package utils

func ContainsString(list []string, val string) bool {
	for _, v := range list {
		if v == val {
			return true
		}
	}
	return false
}
