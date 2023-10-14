// Package util .
package util

// RemoveDuplicates removes duplicate strings from a slice of strings.
func RemoveDuplicates(input []string) []string {
	encountered := map[string]bool{}
	var result []string

	for _, value := range input {
		if encountered[value] == false {
			encountered[value] = true
			result = append(result, value)
		}
	}

	return result
}
