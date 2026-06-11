package catalog

import "fmt"

// PriceForSeatType returns the tier price for a layout seat type.
func PriceForSeatType(tiers PriceTiers, seatType string) (int64, error) {
	switch seatType {
	case SeatTypeStandard:
		return tiers.Standard, nil
	case SeatTypeVIP:
		return tiers.VIP, nil
	case SeatTypeWheelchair:
		return tiers.Wheelchair, nil
	case SeatTypeBlocked:
		return 0, fmt.Errorf("blocked seat cannot be priced")
	default:
		return 0, fmt.Errorf("unknown seat type: %s", seatType)
	}
}

// TotalForSeats sums tier prices for the given seat IDs against a screen layout.
func TotalForSeats(tiers PriceTiers, layout []LayoutSeat, seatIDs []string) (int64, error) {
	byID := make(map[string]LayoutSeat, len(layout))
	for _, s := range layout {
		byID[s.SeatID] = s
	}

	var total int64
	for _, id := range seatIDs {
		seat, ok := byID[id]
		if !ok {
			return 0, fmt.Errorf("unknown seat: %s", id)
		}
		price, err := PriceForSeatType(tiers, seat.Type)
		if err != nil {
			return 0, fmt.Errorf("seat %s: %w", id, err)
		}
		total += price
	}
	return total, nil
}
