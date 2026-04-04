package lyric

import "testing"

func Test_lyricInteractor_CalculateMarginLeft(t *testing.T) {
	type args struct {
		txt string
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "no margin left",
			args: args{
				txt: "Ha",
			},
			want: 0,
		},
		{
			name: "margin left",
			args: args{
				txt: "1. Ha",
			},
			want: -16.58,
		},
		{
			name: "margin left",
			args: args{
				txt: "15. Be",
			},
			want: -24.19,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			li := &lyricInteractor{}

			if got := li.CalculateMarginLeft(tt.args.txt); got != tt.want {
				t.Errorf("lyricInteractor.CalculateMarginLeft() = %v, want %v", got, tt.want)
			}
		})
	}
}
