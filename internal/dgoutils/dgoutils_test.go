package dgoutils

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRange_Array(t *testing.T) {
	type fields struct {
		Low  int
		High int
	}
	tests := []struct {
		name   string
		fields fields
		want   []int
	}{
		{
			name: "range",
			fields: fields{
				Low:  1,
				High: 3,
			},
			want: []int{1, 2, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Range{
				Low:  tt.fields.Low,
				High: tt.fields.High,
			}

			ran := r.Array()
			assert.Equal(t, tt.want, ran)
		})
	}
}

func TestNewRange(t *testing.T) {
	tests := []struct {
		name  string
		arg   string
		want  *Range
		err   error
		isErr bool
	}{
		{
			name: "correct range from 1 to 3",
			arg:  "1-3",
			want: &Range{
				Low:  1,
				High: 3,
			},
			err: nil,
		},
		{
			name: "high is lower than low",
			arg:  "4-3",
			want: nil,
			err:  ErrRangeSyntax,
		},
		{
			name: "low is not an integer",
			arg:  "str-3",
			want: nil,
			err:  &strconv.NumError{},
		},
		{
			name: "high is not an integer",
			arg:  "3-str",
			want: nil,
			err:  &strconv.NumError{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ran, err := NewRange(tt.arg)
			if tt.isErr {
				if assert.Error(t, err) {
					assert.ErrorAs(t, tt.err, err)
				}
			} else {
				assert.Equal(t, tt.want, ran)
			}
		})
	}
}
