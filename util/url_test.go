package util

import (
	"github.com/stretchr/testify/assert"
	"net/url"

	"testing"
)

func TestUpdateQueryParams(t *testing.T) {

	testcases := []struct {
		name             string
		currentUrl       string
		key              string
		value            string
		internalUrlParse func(rawURL string) (*url.URL, error)
		expected         string
		expectErr        error
		hasError         bool
	}{
		{
			name:       "add edit-profile to /explore",
			currentUrl: "/explore",
			key:        "modal",
			value:      "edit-profile",
			expected:   "/explore?modal=edit-profile",
			expectErr:  nil,
			hasError:   false,
		},
		{
			name:       "add filters to /search",
			currentUrl: "/search",
			key:        "modal",
			value:      "filters",
			expected:   "/search?modal=filters",
			expectErr:  nil,
			hasError:   false,
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			originalUrlParse := urlParse

			defer func() {
				urlParse = originalUrlParse
			}()

			urlParse = func(rawURL string) (IURL, error) {
				mockURL := NewMockURL(t)
				var urlValues url.Values = url.Values{}
				mockURL.EXPECT().Query().Return(urlValues)
				return mockURL, nil
			}

			got, _ := UpdateQueryParam(tt.currentUrl, tt.key, tt.value)

			assert.Equal(t, tt.expected, got)

		})

	}
}
