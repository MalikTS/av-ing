package handler

import "testing"

func TestParsePositiveInt(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		fallback  int
		maximum   int
		want      int
		wantError bool
	}{
		{name: "fallback", fallback: 20, maximum: 100, want: 20},
		{name: "valid", value: "3", fallback: 1, maximum: 100, want: 3},
		{name: "clamped", value: "500", fallback: 1, maximum: 100, want: 100},
		{name: "zero", value: "0", fallback: 1, maximum: 100, wantError: true},
		{name: "invalid", value: "abc", fallback: 1, maximum: 100, wantError: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parsePositiveInt(tc.value, "value", tc.fallback, tc.maximum)
			if (err != nil) != tc.wantError {
				t.Fatalf("parsePositiveInt() error = %v, wantError %v", err, tc.wantError)
			}
			if got != tc.want {
				t.Fatalf("parsePositiveInt() = %d, want %d", got, tc.want)
			}
		})
	}
}
