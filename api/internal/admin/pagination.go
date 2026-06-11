package admin

const (
	defaultPageLimit = 50
	maxPageLimit     = 100
)

// ClampPageLimit normalizes a requested page size.
func ClampPageLimit(limit int) int {
	if limit <= 0 {
		return defaultPageLimit
	}
	if limit > maxPageLimit {
		return maxPageLimit
	}
	return limit
}

// SkipFor returns the MongoDB skip offset for a 1-based page.
func SkipFor(page, limit int) int {
	if page <= 1 {
		return 0
	}
	return (page - 1) * ClampPageLimit(limit)
}
