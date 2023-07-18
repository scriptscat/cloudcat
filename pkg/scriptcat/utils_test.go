package scriptcat

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvCron(t *testing.T) {
	type args struct {
		cron string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{"case1", args{"* * * * *"}, "0 * * * * *", assert.NoError},
		{"case2", args{"* 10-23 once * *"}, "0 0 10 * * *", assert.NoError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvCron(tt.args.cron)
			if !tt.wantErr(t, err, fmt.Sprintf("ConvCron(%v)", tt.args.cron)) {
				return
			}
			assert.Equalf(t, tt.want, got, "ConvCron(%v)", tt.args.cron)
		})
	}
}
