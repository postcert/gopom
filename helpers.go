package main

func btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}

func itob(u int) bool {
	return u == 1
}

func getIndexWithFallback(index int, fallback int, list []int) int {
	if index < len(list) {
		return list[index]
	}
	return fallback
}
