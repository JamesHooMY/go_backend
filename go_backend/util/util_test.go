package util_test

import (
	"errors"
	"fmt"
	"testing"

	"go_backend/util"

	"github.com/stretchr/testify/assert"
)

func Test_ErrorWrapper(t *testing.T) {
	type args struct {
		existErr error
		newErr   error
	}
	type expected struct {
		err error
	}
	type testCase struct {
		name     string
		args     args
		expected expected
	}

	testCases := []testCase{
		{
			name: "test1",
			args: args{
				existErr: nil,
				newErr:   nil,
			},
			expected: expected{
				err: fmt.Errorf("%w, %w", nil, nil),
			},
		},
		{
			name: "test2",
			args: args{
				existErr: nil,
				newErr:   errors.New("new error"),
			},
			expected: expected{
				err: errors.New("new error"),
			},
		},
		{
			name: "test3",
			args: args{
				existErr: errors.New("exist error"),
				newErr:   nil,
			},
			expected: expected{
				err: errors.New("exist error"),
			},
		},
		{
			name: "test4",
			args: args{
				existErr: errors.New("exist error"),
				newErr:   errors.New("new error"),
			},
			expected: expected{
				err: fmt.Errorf("%w, %w", errors.New("exist error"), errors.New("new error")),
			},
		},
	}

	for _, tc := range testCases {
		err := util.ErrorWrapper(tc.args.existErr, tc.args.newErr)
		assert.Equal(t, tc.expected.err, err, tc.name)
	}
}
