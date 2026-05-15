package diff

// KeyStatus represents the comparison status of a key across environments.
type KeyStatus string

const (
	StatusMissing   KeyStatus = "missing"
	StatusMismatch  KeyStatus = "mismatch"
	StatusMatch     KeyStatus = "match"
)

// Result holds the comparison result for a single key.
type Result struct {
	Key    string            `json:"key"`
	Status KeyStatus         `json:"status"`
	Values map[string]string `json:"values,omitempty"`
}

// Compare takes a map of environment name -> parsed key/value pairs and
// returns a slice of Result describing missing or mismatched keys.
func Compare(envs map[string]map[string]string) []Result {
	allKeys := collectKeys(envs)
	envNames := envNames(envs)

	var results []Result

	for _, key := range allKeys {
		values := make(map[string]string)
		presentIn := 0

		for _, env := range envNames {
			if val, ok := envs[env][key]; ok {
				values[env] = val
				presentIn++
			}
		}

		if presentIn < len(envNames) {
			results = append(results, Result{Key: key, Status: StatusMissing, Values: values})
			continue
		}

		if hasMismatch(values) {
			results = append(results, Result{Key: key, Status: StatusMismatch, Values: values})
		}
	}

	return results
}

func collectKeys(envs map[string]map[string]string) []string {
	seen := make(map[string]struct{})
	var keys []string
	for _, kv := range envs {
		for k := range kv {
			if _, ok := seen[k]; !ok {
				seen[k] = struct{}{}
				keys = append(keys, k)
			}
		}
	}
	sort.Strings(keys)
	return keys
}

func envNames(envs map[string]map[string]string) []string {
	names := make([]string, 0, len(envs))
	for name := range envs {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func hasMismatch(values map[string]string) bool {
	var first string
	set := false
	for _, v := range values {
		if !set {
			first = v
			set = true
			continue
		}
		if v != first {
			return true
		}
	}
	return false
}
