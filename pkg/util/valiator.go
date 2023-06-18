package util

const (
	USD = "USD"
	EUR = "EUR"
	CAD = "CAD"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case EUR, USD, CAD:
		return true
	default:
		return false
	}
}
