package cli

import (
	"fmt"
	"strings"
)

func parseIncludes(includeValues []string, preloadValues []string, allowedIncludes map[string]struct{}) ([]string, error) {
	seenIncludes := map[string]struct{}{}
	parsedIncludes := []string{}

	for _, rawValue := range append(includeValues, preloadValues...) {
		for _, includeName := range strings.Split(rawValue, ",") {
			includeName = strings.TrimSpace(includeName)
			if includeName == "" {
				continue
			}
			if _, ok := allowedIncludes[includeName]; !ok {
				return nil, fmt.Errorf("include %q is not supported", includeName)
			}
			if _, ok := seenIncludes[includeName]; ok {
				continue
			}
			seenIncludes[includeName] = struct{}{}
			parsedIncludes = append(parsedIncludes, includeName)
		}
	}

	return parsedIncludes, nil
}

func includeSet(includeNames ...string) map[string]struct{} {
	includes := map[string]struct{}{}
	for _, includeName := range includeNames {
		includes[includeName] = struct{}{}
	}
	return includes
}
