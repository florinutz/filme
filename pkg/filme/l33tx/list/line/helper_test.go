package line

import (
	"math"
	"testing"
)

func Test_ToBytes(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    uint64
		wantErr bool
	}{
		{
			name:    "simple",
			input:   "7096kb",
			want:    uint64(7096 * 1024),
			wantErr: false,
		},
		{
			name:    "multiple with space",
			input:   "12 GB",
			want:    uint64(12 * 1024 * 1024 * 1024),
			wantErr: false,
		},
		{
			name:    "decimals with dot",
			input:   "12.4 kb",
			want:    uint64(math.Floor(12.4 * 1024)),
			wantErr: false,
		},
		{
			name:    "decimals with comma",
			input:   "17,3 MiB",
			want:    uint64(math.Floor(17.3 * 1024 * 1024)),
			wantErr: false,
		},
		{
			name:    "bytes",
			input:   "17,3B",
			want:    uint64(math.Floor(17.3)),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToBytes(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ToBytes() got = %v, want %v", got, tt.want)
			}
		})
	}
}
