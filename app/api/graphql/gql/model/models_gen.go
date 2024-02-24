// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type Mutation struct {
}

type Query struct {
}

type User struct {
	ID        *string `json:"id,omitempty"`
	Email     string  `json:"email"`
	Password  *string `json:"password,omitempty"`
	Mobile    *string `json:"mobile,omitempty"`
	Name      string  `json:"name"`
	Age       *int    `json:"age,omitempty"`
	CreatedAt *string `json:"createdAt,omitempty"`
	UpdatedAt *string `json:"updatedAt,omitempty"`
	DeletedAt *string `json:"deletedAt,omitempty"`
}

type UserInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
