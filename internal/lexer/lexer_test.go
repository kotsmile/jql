package lexer

import (
	"reflect"
	"testing"

	"github.com/kotsmile/jql/internal/lexer/token"
	"github.com/kotsmile/jql/util"
)

func Test_NextWord(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name  string
		args  args
		want0 string
		want1 string
	}{
		{
			name: "hello world",
			args: args{
				s: "hello world",
			},
			want0: "hello",
			want1: " world",
		},
		{
			name: "single word",
			args: args{
				s: "hello",
			},
			want0: "hello",
			want1: "",
		},
		{
			name: "with leading whitespaces",
			args: args{
				s: "     hello world",
			},
			want0: "hello",
			want1: " world",
		},
		{
			name: "with another whitespace",
			args: args{
				s: "     \thello world",
			},
			want0: "hello",
			want1: " world",
		},
		{
			name: "empty",
			args: args{
				s: "   ",
			},
			want0: "",
			want1: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got0, got1 := nextWord(tt.args.s, separators, symbols)

			if got0 != tt.want0 || got1 != tt.want1 {
				t.Errorf(
					"nextWord(\"%s\") = \"%s\",\"%s\" , want \"%s\",\"%s\"",
					tt.args.s, got0, got1, tt.want0, tt.want1,
				)
			}
		})
	}
}

func Test_Next(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    *token.Token
		wantErr bool
	}{
		{
			name: "string",
			args: args{
				s: "\"hello world\"",
			},
			want: token.New(
				token.String,
				"hello world",
			),
			wantErr: false,
		},
		{
			name: "string with whitespaces",
			args: args{
				s: "\"hello world \"",
			},
			want: token.New(
				token.String,
				"hello world ",
			),
			wantErr: false,
		},
		{
			name: "unterminated string",
			args: args{
				s: "\"hello world",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "word",
			args: args{
				s: "hello",
			},
			want: token.New(
				token.Word,
				"hello",
			),
			wantErr: false,
		},
		{
			name: "word2",
			args: args{
				s: "hello world",
			},
			want: token.New(
				token.Word,
				"hello",
			),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := New(util.NewLoggerTest())
			lexer.Lex(tt.args.s)

			got, err := lexer.Next()
			if (err != nil) != tt.wantErr {
				t.Errorf("Next() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Next() = %v, want %v", got, tt.want)
			}
		})
	}
}
