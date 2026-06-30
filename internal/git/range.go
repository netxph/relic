package git

import (
	"fmt"
	"strings"
)

// ParsedRange holds the resolved from/to refs for a git range query.
type ParsedRange struct {
	From string
	To   string
}

// ParseRange accepts "<hash>" or "<hash>..<hash>" and returns a ParsedRange.
// Single hash is expanded to <hash>..HEAD.
// Minimum hash length is 7 characters.
func ParseRange(input string) (ParsedRange, error) {
	if input == "" {
		return ParsedRange{}, fmt.Errorf("--range is required")
	}

	if strings.Contains(input, "..") {
		parts := strings.SplitN(input, "..", 2)
		from, to := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
		if err := validateHash(from); err != nil {
			return ParsedRange{}, err
		}
		if err := validateHash(to); err != nil {
			return ParsedRange{}, err
		}
		return ParsedRange{From: from, To: to}, nil
	}

	hash := strings.TrimSpace(input)
	if err := validateHash(hash); err != nil {
		return ParsedRange{}, err
	}
	return ParsedRange{From: hash, To: "HEAD"}, nil
}

func validateHash(h string) error {
	if len(h) < 7 {
		return fmt.Errorf("hash %q is too short (minimum 7 characters)", h)
	}
	return nil
}
