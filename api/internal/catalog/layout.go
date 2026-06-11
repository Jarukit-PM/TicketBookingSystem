package catalog

// SellableSeatCount returns layout seats that can be sold (excludes blocked).
func SellableSeatCount(layout ScreenLayout) int {
	count := 0
	for _, seat := range layout.Seats {
		if seat.Type != SeatTypeBlocked {
			count++
		}
	}
	return count
}
