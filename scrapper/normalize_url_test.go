package scrapper

import "testing"

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name        string
		inputURL    string
		expectError bool

		expectedURL string
	}{
		{
			name:        "remove schema",
			inputURL:    "https://www.boot.dev/blog/path",
			expectedURL: "www.boot.dev/blog/path",
		},
		{
			name:        "remove query params",
			inputURL:    "https://www.boot.dev/blog/path?q=param-bad",
			expectedURL: "www.boot.dev/blog/path",
		},
		{
			name:        "invalid url",
			inputURL:    "sdfasdfas@asassaas.com",
			expectedURL: "",
			expectError: true,
		},
		{
			name:        "empty url",
			inputURL:    "",
			expectError: true,
			expectedURL: "",
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := NormalizeURL(tc.inputURL)
			if err != nil && !tc.expectError {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, tc.name, err)
				return
			}
			if actual != tc.expectedURL {
				t.Errorf("Test %v - %s FAIL: expected %v, actual: %v", i, tc.name, tc.expectedURL, actual)
			}
		})
	}
}
