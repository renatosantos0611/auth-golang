package repositories

import (
	"auth-golang/internal/database"
	"auth-golang/internal/models"
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UsersRepositoryInterface interface {
	Create(ctx context.Context, user *models.User) error
	FindByID(ctx context.Context, id string) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByRefreshToken(ctx context.Context, refreshToken string) (*models.User, error)
	Update(ctx context.Context, user *models.User) (*models.User, error)
}

type UserRepository struct {
	collection *mongo.Collection
}

func NewUsersRepository(db *database.Service) *UserRepository {
	return &UserRepository{
		collection: db.Database.Collection("users"),
	}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	_, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %v", err)
	}

	err = r.collection.FindOne(ctx, bson.D{{Key: "_id", Value: objectID}}).Decode(&user)

	if err != nil {
		if errors.Is((err), mongo.ErrNoDocuments) {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User

	err := r.collection.FindOne(ctx, bson.D{{Key: "email", Value: email}}).Decode(&user)

	if err != nil {
		if errors.Is((err), mongo.ErrNoDocuments) {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) FindByRefreshToken(ctx context.Context, token string) (*models.User, error) {
	var user models.User

	err := r.collection.FindOne(ctx, bson.D{{Key: "refresh_token", Value: token}}).Decode(&user)

	if err != nil {
		if errors.Is((err), mongo.ErrNoDocuments) {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *models.User) (*models.User, error) {

	if user.ID.IsZero() {
		return nil, fmt.Errorf("user ID is required")
	}

	user.UpdatedAt = time.Now()

	_, err := r.collection.UpdateOne(
		ctx,
		bson.D{{Key: "_id", Value: user.ID}},
		bson.D{
			{Key: "$set", Value: user},
		},
	)

	if err != nil {
		if errors.Is((err), mongo.ErrNoDocuments) {
			return nil, nil
		}

		return nil, err
	}

	return user, nil
}
