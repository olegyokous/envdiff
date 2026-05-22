package diff

// Status describes the comparison outcome for a single key.
type Status string

const (
	// StatusMatch indicates all environments define the key with the same value.
	StatusMatch Status = "match"
	// StatusMissing indicates the key is absent in one or more environments.
	StatusMissing Status = "missing"
	// StatusMismatch indicates the key is present everywhere but values differ.
	StatusMismatch Status = "mismatch"
)

// Result holds the comparison outcome for a single key across environments.
type Result struct {
	// Key is the environment variable name.
	Key string `json:"key"`
	// Status is the overall comparison outcome.
	Status Status `json:"status"`
	// Values maps each environment name to its observed value.
	// A missing key is represented by an empty string with StatusMissing.
	Values map[string]string `json:"values"`
}

// IsMissing reports whether the result represents a missing key.
func (r Result) IsMissing() bool { return r.Status == StatusMissing }

// IsMismatch reports whether the result represents a value mismatch.
func (r Result) IsMismatch() bool { return r.Status == StatusMismatch }

// IsMatch reports whether the result represents a full match.
func (r Result) IsMatch() bool { return r.Status == StatusMatch }
