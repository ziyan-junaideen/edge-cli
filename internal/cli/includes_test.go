package cli

import (
	"reflect"
	"testing"
)

func TestParseIncludesSplitsDeduplicatesAndPreservesOrder(t *testing.T) {
	includes, err := parseIncludes(
		[]string{"addresses,merchant", "addresses"},
		[]string{"payment_demands"},
		includeSet("addresses", "merchant", "payment_demands"),
	)
	if err != nil {
		t.Fatalf("parseIncludes returned error: %v", err)
	}

	expectedIncludes := []string{"addresses", "merchant", "payment_demands"}
	if !reflect.DeepEqual(includes, expectedIncludes) {
		t.Fatalf("expected %v, got %v", expectedIncludes, includes)
	}
}

func TestParseIncludesRejectsUnsupportedInclude(t *testing.T) {
	_, err := parseIncludes([]string{"addresses"}, nil, includeSet("merchant"))
	if err == nil {
		t.Fatal("expected unsupported include to return error")
	}
}
