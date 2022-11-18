package utils

import "testing"

func TestGetNextHalfStep(t *testing.T) {
	type args struct {
		pitch string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "C",
			args: args{
				pitch: "C",
			},
			want: "C#",
		},
		{
			name: "C#",
			args: args{
				pitch: "C#",
			},
			want: "D",
		},
		{
			name: "D",
			args: args{
				pitch: "D",
			},
			want: "D#",
		},
		{
			name: "D#",
			args: args{
				pitch: "D#",
			},
			want: "E",
		},
		{
			name: "E",
			args: args{
				pitch: "E",
			},
			want: "F",
		},
		{
			name: "B",
			args: args{
				pitch: "B",
			},
			want: "C",
		},
		{
			name: "B#",
			args: args{
				pitch: "B#",
			},
			want: "C#",
		},
		{
			name: "Bx",
			args: args{
				pitch: "Bx",
			},
			want: "D",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetNextHalfStep(tt.args.pitch); got != tt.want {
				t.Errorf("getNextHalfStep() = %v, want %v", got, tt.want)
			}
		})
	}
}
