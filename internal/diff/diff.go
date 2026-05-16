package diff

import "sort"

// Result holds the comparison outcome for a single key across environments.
type Result struct {
	Key    string            `json:"key"`
	Status string            `json:"status"` // "match", "missing", "mismatch"
	Values map[string]string `json:"values"`
}

// Compare takes a map of environment name -> key/value pairs and returns
// a slice of Result describing the diff across all environments.
func Compare(envs map[string]map[string]string) []Result {
	keys := collectKeys(envs)
	names := envNames(envs)

	var results []Result
	for _, key := range keys {
		values := make(map[string]string, len(names))
		for _, name := range names {
			values[name] = envs[name][key]
		}

		status := "match"
		for _, name := range names {
			if _, ok := envs[name][key]; !ok {
				status = "missing"
				break
			}
		}
		if status != "missing" && hasMismatch(values) {
			status = "mismatch"
		}

		results = append(results, Result{
			Key:    key,
			Status: status,
			Values: values,
		})
	}
	return results
}

func collectKeys(envs map[string]map[string]string) []string {
	seen := make(map[string]struct{})
	for _, kv := range envs {
		for k := range kv {
			seen[k] = struct{}{}
		}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
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
	var ref string
	first := true
	for _, v := range values {
		if first {
			ref = v
			first = false
			continue
		}
		if v != ref {
			return true
		}
	}
	return false
}
