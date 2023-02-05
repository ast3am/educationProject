package user

import "context"

type Repository interface {
	Create(ctx context.Context, user *UserModel) error
	MakeFriends(ctx context.Context, sourceId, targetId string) (string, error)
	Delete(ctx context.Context, id string) (string, error)
	FindFriend(ctx context.Context, id string) (ufriends []*UserModel, err error)
	UpdateAge(ctx context.Context, id, age string) error
	MakeID() string
}
