package bot

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_sensorIDFromCallback(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		data    string
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "positive",
			data:    "sensor_test",
			want:    "test",
			wantErr: assert.NoError,
		},
		{
			name:    "error",
			data:    "_test",
			want:    "",
			wantErr: assert.Error,
		},
		{
			name:    "two prefixes",
			data:    "sensor_sensor_test",
			want:    "sensor_test",
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := sensorIDFromCallback(tt.data)
			tt.wantErr(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
