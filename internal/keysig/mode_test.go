package keysig

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMode_GetHumanized(t *testing.T) {
	type args struct {
		fifth int
	}
	tests := []struct {
		name string
		m    Mode
		args args
		want string
	}{
		{
			name: "A minor",
			m:    NewMode("minor"),
			args: args{
				fifth: 0,
			},
			want: "la = a",
		},
		{
			name: "F major",
			m:    NewMode(""),
			args: args{
				fifth: -1,
			},
			want: "do = f",
		},
		{
			name: "D major",
			m: Mode{
				humanized: "do = d",
			},
			args: args{
				fifth: 2,
			},
			want: "do = d",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.GetHumanized(tt.args.fifth); got != tt.want {
				t.Errorf("Mode.GetHumanized() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMode_GetScaleSteps(t *testing.T) {
	tests := []struct {
		name string
		m    Mode
		want []float64
	}{
		{
			name: "Major",
			m:    NewMode("major"),
			want: []float64{1, 1, 0.5, 1, 1, 1, 0.5},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.GetScaleSteps(); !assert.Equal(t, tt.want, got) {
				t.Errorf("Mode.GetScaleSteps() = %v, want %v", got, tt.want)
			}
		})
	}
}
