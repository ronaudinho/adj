package adj

import (
	"testing"
)

func TestSet(t *testing.T) {
	tests := []struct {
		name     string
		parallel int
		max      int
		want     int
	}{
		{
			name:     "less than",
			parallel: 5,
			max:      10,
			want:     5,
		},
		{
			name:     "equal",
			parallel: 10,
			max:      10,
			want:     10,
		},
		{
			name:     "more than",
			parallel: 15,
			max:      10,
			want:     10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Set(tt.parallel, tt.max)
			if got != tt.want {
				t.Errorf("want %d, got %d", tt.want, got)
			}
		})
	}
}
