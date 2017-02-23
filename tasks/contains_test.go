package tasks_test

func contains(a []string, v string) bool {

	for _, av := range a {
		if av == v {
			return true
		}
	}

	return false
}
