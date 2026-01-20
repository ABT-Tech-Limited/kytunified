package kyt

// RiskLevel represents the unified risk level across all KYT providers.
// This is the standardized output that clients will receive regardless of
// which underlying provider (Beosin, Chainalysis, etc.) is being used.
type RiskLevel string

const (
	// RiskLevelLow indicates minimal risk. The address/transaction appears clean
	// with no significant exposure to illicit activities.
	RiskLevelLow RiskLevel = "Low"

	// RiskLevelMedium indicates moderate risk requiring attention.
	// Some exposure to potentially risky activities detected.
	RiskLevelMedium RiskLevel = "Medium"

	// RiskLevelHigh indicates significant risk requiring action.
	// Substantial exposure to high-risk activities detected.
	RiskLevelHigh RiskLevel = "High"

	// RiskLevelCritical indicates severe risk requiring immediate action.
	// Direct exposure to sanctioned entities, hackers, or other critical threats.
	// Note: This maps from Beosin's "Severe" level.
	RiskLevelCritical RiskLevel = "Critical"

	// RiskLevelUnknown indicates the risk level could not be determined.
	// This may occur when the provider returns an unexpected value or
	// when the assessment fails.
	RiskLevelUnknown RiskLevel = "Unknown"
)

// Severity returns the numeric severity (0-3) for comparison operations.
// Higher numbers indicate higher risk.
// Returns -1 for unknown risk levels.
func (r RiskLevel) Severity() int {
	switch r {
	case RiskLevelLow:
		return 0
	case RiskLevelMedium:
		return 1
	case RiskLevelHigh:
		return 2
	case RiskLevelCritical:
		return 3
	default:
		return -1
	}
}

// IsHigherOrEqual returns true if this risk level is >= the other.
// Useful for threshold-based decisions.
func (r RiskLevel) IsHigherOrEqual(other RiskLevel) bool {
	return r.Severity() >= other.Severity()
}

// IsHigherThan returns true if this risk level is > the other.
func (r RiskLevel) IsHigherThan(other RiskLevel) bool {
	return r.Severity() > other.Severity()
}

// IsValid returns true if this is a valid, known risk level.
func (r RiskLevel) IsValid() bool {
	return r.Severity() >= 0
}

// String returns the string representation of the risk level.
func (r RiskLevel) String() string {
	return string(r)
}

// ParseRiskLevel parses a string into RiskLevel.
// It handles case-insensitive matching and maps provider-specific
// values (like "Severe") to our unified levels.
func ParseRiskLevel(s string) RiskLevel {
	switch s {
	case "Low", "low", "LOW":
		return RiskLevelLow
	case "Medium", "medium", "MEDIUM":
		return RiskLevelMedium
	case "High", "high", "HIGH":
		return RiskLevelHigh
	case "Critical", "critical", "CRITICAL", "Severe", "severe", "SEVERE":
		return RiskLevelCritical
	default:
		return RiskLevelUnknown
	}
}

// AllRiskLevels returns all valid risk levels in order of severity.
func AllRiskLevels() []RiskLevel {
	return []RiskLevel{
		RiskLevelLow,
		RiskLevelMedium,
		RiskLevelHigh,
		RiskLevelCritical,
	}
}
