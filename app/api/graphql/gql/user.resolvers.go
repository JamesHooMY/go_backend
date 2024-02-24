package gql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.44

import (
	"context"
	"fmt"
	"go_backend/app/api/graphql/gql/generated"
	"go_backend/app/api/graphql/gql/model"
)

// CreateUser is the resolver for the createUser field.
func (r *mutationResolver) CreateUser(ctx context.Context, input model.UserInput) (*model.User, error) {
	// panic(fmt.Errorf("not implemented: CreateUser - createUser"))
	tempID := "1"
	record := &model.User{
		ID:       &tempID,
		Email:    input.Email,
		Password: &input.Password,
	}

	return record, nil
}

// UpdateUser is the resolver for the updateUser field.
func (r *mutationResolver) UpdateUser(ctx context.Context, input model.UserInput) (*model.User, error) {
	panic(fmt.Errorf("not implemented: UpdateUser - updateUser"))
}

// DeleteUser is the resolver for the deleteUser field.
func (r *mutationResolver) DeleteUser(ctx context.Context, id string) (bool, error) {
	panic(fmt.Errorf("not implemented: DeleteUser - deleteUser"))
}

// Users is the resolver for the users field.
func (r *queryResolver) Users(ctx context.Context, id *string) ([]*model.User, error) {
	// panic(fmt.Errorf("not implemented: Users - users"))
	tempID := "1"
	tempEmail := "your@email.com"
	records := []*model.User{
		{
			ID:    &tempID,
			Email: tempEmail,
		},
	}
	return records, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
