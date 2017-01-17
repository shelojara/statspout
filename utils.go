package statspout

func contains(slice []string, name string) bool {
	for _, n := range slice {
		if n == name {
			return true
		}
	}

	return false
}
