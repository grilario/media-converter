package styles

func GetCursor(cursor int, current int) string {
	indicator := " "
	if cursor == current {
		indicator = ">"
	}

	return indicator
}
