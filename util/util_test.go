package util

import "testing"

func Test_NextUtil(t *testing.T) {
	type args struct {
		slice []int
	}
	tests := []struct {
		name string
		args args
		want int
		ok   bool
	}{
		{
			name: "ok",
			args: args{
				slice: []int{1, 2, 3},
			},
			want: 1,
			ok:   true,
		},
		{
			name: "empty",
			args: args{
				slice: []int{},
			},
			want: 0,
			ok:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := Next(&tt.args.slice)
			if got != tt.want {
				t.Errorf("Next() got = %v, want %v", got, tt.want)
			}
			if ok != tt.ok {
				t.Errorf("Next() ok = %v, want %v", ok, tt.ok)
			}
		})
	}
}
