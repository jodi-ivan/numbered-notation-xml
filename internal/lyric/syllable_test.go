package lyric

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitSyllable(t *testing.T) {
	type args struct {
		word string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "dan",
			args: args{
				word: "dan",
			},
			want: []string{"dan"},
		},
		{
			name: "Kata",
			args: args{
				word: "Kata",
			},
			want: []string{"Ka", "ta"},
		},
		{
			name: "yakin",
			args: args{
				word: "yakin",
			},
			want: []string{"ya", "kin"},
		},
		{
			name: "yakin",
			args: args{
				word: "yakin",
			},
			want: []string{"ya", "kin"},
		},
		{
			name: "lumping",
			args: args{
				word: "lumping",
			},
			want: []string{"lum", "ping"},
		},
		{
			name: "mengantuk",
			args: args{
				word: "mengantuk",
			},
			want: []string{"me", "ngan", "tuk"},
		},
		{
			name: "menyair",
			args: args{
				word: "menyair",
			},
			want: []string{"me", "nya", "ir"},
		},
		{
			name: "yang",
			args: args{
				word: "yang",
			},
			want: []string{"yang"},
		},
		{
			name: "anti",
			args: args{
				word: "anti",
			},
			want: []string{"an", "ti"},
		},
		{
			name: "tanya",
			args: args{
				word: "tanya",
			},
			want: []string{"ta", "nya"},
		},
		{
			name: "kacau",
			args: args{
				word: "kacau",
			},
			want: []string{"ka", "ca", "u"},
		},
		{
			name: "mau",
			args: args{
				word: "mau",
			},
			want: []string{"ma", "u"},
		},
		{
			name: "mengganggu",
			args: args{
				word: "mengganggu",
			},
			want: []string{"meng", "gang", "gu"},
		},
		{
			name: "mempertanggungjawabkan",
			args: args{
				word: "mempertanggungjawabkan",
			},
			want: []string{"mem", "per", "tang", "gung", "ja", "wab", "kan"},
		},
		{
			name: "langit",
			args: args{
				word: "langit",
			},
			want: []string{"la", "ngit"},
		},
		{
			name: "jumat",
			args: args{
				word: "jumat",
			},
			want: []string{"ju", "mat"},
		},
		{
			name: "menyala!",
			args: args{
				word: "menyala!",
			},
			want: []string{"me", "nya", "la!"},
		},
		{
			name: "salah!",
			args: args{
				word: "salah!",
			},
			want: []string{"sa", "lah!"},
		},
		{
			name: "cemerlang!",
			args: args{
				word: "cemerlang!",
			},
			want: []string{"ce", "mer", "lang!"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SplitSyllable(tt.args.word); !assert.Equal(t, tt.want, got) {
				t.Errorf("SplitSyllable() = %v, want %v", got, tt.want)
			}
		})
	}
}
