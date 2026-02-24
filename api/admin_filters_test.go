package api

import "testing"

func TestNormalizeMonitorGroupFilter(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "empty", input: "", expected: ""},
		{name: "all", input: "all", expected: ""},
		{name: "all uppercase", input: "ALL", expected: ""},
		{name: "trimmed", input: "  Core ", expected: "core"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if actual := normalizeMonitorGroupFilter(testCase.input); actual != testCase.expected {
				t.Fatalf("expected %q, got %q", testCase.expected, actual)
			}
		})
	}
}

func TestNormalizeAdminAuditResultFilter(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "empty", input: "", expected: ""},
		{name: "all", input: "all", expected: ""},
		{name: "all uppercase", input: "ALL", expected: ""},
		{name: "success", input: "success", expected: "success"},
		{name: "failure", input: "failure", expected: "failure"},
		{name: "invalid", input: "unknown", expected: ""},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if actual := normalizeAdminAuditResultFilter(testCase.input); actual != testCase.expected {
				t.Fatalf("expected %q, got %q", testCase.expected, actual)
			}
		})
	}
}
