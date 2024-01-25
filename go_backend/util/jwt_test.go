package util_test

import (
	"testing"
	"time"

	"go_backend/util"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func Test_GenerateJwtToken(t *testing.T) {
	type args struct {
		id   uint
		name string
	}
	type testCase struct {
		name string
		args args
	}

	testCases := []testCase{
		{
			name: "test1",
			args: args{
				id:   1,
				name: "james",
			},
		},
	}

	for _, tc := range testCases {
		token, err := util.GenerateJwtToken(tc.args.id, tc.args.name)
		assert.NoError(t, err, tc.name)
		assert.NotEmpty(t, token, tc.name)
	}
}

func Test_ParseJwtToken(t *testing.T) {
	type args struct {
		id   uint
		name string
	}
	type expected struct {
		claims util.Claims
		err    error
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
				id:   1,
				name: "james",
			},
			expected: expected{
				claims: util.Claims{
					ID:   1,
					Name: "james",
					RegisteredClaims: jwt.RegisteredClaims{
						ExpiresAt: jwt.NewNumericDate(time.Unix(int64(1706087945), 0)),
						IssuedAt:  jwt.NewNumericDate(time.Unix(int64(1706086145), 0)),
					},
				},
				err: nil,
			},
		},
	}

	for _, tc := range testCases {
		token, _ := util.GenerateJwtToken(tc.args.id, tc.args.name)
		claims, err := util.ParseJwtToken(token)
		assert.NoError(t, err, tc.name+"_ParseJwtToken")
		assert.NotEmpty(t, claims, tc.name+"_ParseJwtToken")
		assert.Equal(t, tc.expected.claims.ID, claims.ID, tc.name+"_ParseJwtToken")
		assert.Equal(t, tc.expected.claims.Name, claims.Name, tc.name+"_ParseJwtToken")

	}
}

func Test_ParseJwtToken_Expired(t *testing.T) {
	type args struct {
		token string
	}
	type expected struct {
		claims util.Claims
		err    error
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
				token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwibmFtZSI6ImphbWVzIiwiZXhwIjoxNzA2MDg3OTQ1LCJpYXQiOjE3MDYwODYxNDV9.4aq7EpflxeOb1J7ukrT8xJjsoxq0732Nil6jB_JIP6A",
			},
			expected: expected{
				claims: util.Claims{
					ID:   1,
					Name: "james",
					RegisteredClaims: jwt.RegisteredClaims{
						ExpiresAt: jwt.NewNumericDate(time.Unix(int64(1706087945), 0)),
						IssuedAt:  jwt.NewNumericDate(time.Unix(int64(1706086145), 0)),
					},
				},
				err: nil,
			},
		},
	}

	for _, tc := range testCases {
		claims, err := util.ParseJwtToken(tc.args.token)
		assert.ErrorIs(t, err, jwt.ErrTokenExpired, tc.name)
		assert.Empty(t, claims, tc.name)
	}
}
