package main

func arrayContains(array []string, val string) bool {
	for _, v := range array {
		if val == v {
			return true
		}
	}

	return false
}

func arrayRemoveAll(arr []string, items ...string) []string {
	var newArr []string
	for _, a := range arr {
		if !arrayContains(items, a) {
			newArr = append(newArr, a)
		}
	}
	return newArr
}

func arrayRepeat(text string, count int) []string {
	a := make([]string, count)
	for i := 0; i < count; i++ {
		a[i] = text
	}
	return a
}

func arrayStringToInterface(s []string) []interface{} {
	a := make([]interface{}, len(s))
	for i, e := range s {
		a[i] = e
	}
	return a
}
