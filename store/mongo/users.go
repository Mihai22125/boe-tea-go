package mongo

import (
	"context"
	"time"

	"github.com/VTGare/boe-tea-go/store"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type userStore struct {
	client *mongo.Client
	db     *mongo.Database
	col    *mongo.Collection
}

func UserStore(client *mongo.Client, database string) store.UserStore {
	db := client.Database(database)
	col := db.Collection("users")

	return &userStore{
		client: client,
		db:     db,
		col:    col,
	}
}

func (u *userStore) User(ctx context.Context, userID string) (*store.User, error) {
	res := u.col.FindOneAndUpdate(
		ctx,
		bson.M{
			"user_id": userID,
		},
		bson.M{
			"$setOnInsert": store.DefaultUser(userID),
		},
		options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After),
	)

	var user store.User
	err := res.Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *userStore) CreateUser(ctx context.Context, id string) (*store.User, error) {
	user := store.DefaultUser(id)

	_, err := u.col.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userStore) CreateCrosspostGroup(ctx context.Context, userID string, group *store.Group) (*store.User, error) {
	res := u.col.FindOneAndUpdate(
		ctx,
		bson.M{"user_id": userID, "channel_groups.name": bson.M{"$ne": group.Name}, "channel_groups.parent": bson.M{"$ne": group.Parent}},
		bson.M{"$push": bson.M{"channel_groups": group}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	var user store.User
	err := res.Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *userStore) DeleteCrosspostGroup(ctx context.Context, userID, group string) (*store.User, error) {
	res := u.col.FindOneAndUpdate(
		ctx,
		bson.M{"user_id": userID, "channel_groups.name": group},
		bson.M{"$pull": bson.M{"channel_groups": bson.M{"name": group}}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	var user store.User
	err := res.Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *userStore) AddCrosspostChannel(ctx context.Context, userID, group, child string) (*store.User, error) {
	res := u.col.FindOneAndUpdate(
		ctx,
		bson.M{"user_id": userID, "channel_groups.name": group},
		bson.M{"$addToSet": bson.M{"channel_groups.$.children": child}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	var user store.User
	err := res.Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *userStore) DeleteCrosspostChannel(ctx context.Context, userID, group, child string) (*store.User, error) {
	res := u.col.FindOneAndUpdate(
		ctx,
		bson.M{"user_id": userID, "channel_groups.name": group},
		bson.M{"$pull": bson.M{"channel_groups.$.children": child}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	var user store.User
	err := res.Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *userStore) UpdateUser(ctx context.Context, user *store.User) (*store.User, error) {
	user.UpdatedAt = time.Now()
	_, err := u.col.ReplaceOne(
		ctx,
		bson.M{"user_id": user.ID},
		user,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}
