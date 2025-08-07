package service

import (
	"context"
	"time"

	"go-grpcgateway/internal/db"
	"go-grpcgateway/internal/models"
	"go-grpcgateway/pkg/pb"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// UserService implements the UserService gRPC service
type UserService struct {
	pb.UnimplementedUserServiceServer
	db         *db.MongoDB
	collection *mongo.Collection
}

// NewUserService creates a new UserService
func NewUserService(database *db.MongoDB) *UserService {
	return &UserService{
		db:         database,
		collection: database.GetCollection("users"),
	}
}

// CreateUser creates a new user
func (s *UserService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.User, error) {
	logrus.WithFields(logrus.Fields{
		"name":  req.Name,
		"email": req.Email,
	}).Info("Creating new user")

	user := models.NewUser(req.Name, req.Email, req.Phone)

	result, err := s.collection.InsertOne(ctx, user)
	if err != nil {
		logrus.WithError(err).Error("Failed to create user")
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	user.ID = result.InsertedID.(primitive.ObjectID)

	return s.modelToProto(user), nil
}

// GetUser retrieves a user by ID
func (s *UserService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
	logrus.WithField("id", req.Id).Info("Getting user")

	objectID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID: %v", err)
	}

	var user models.User
	err = s.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		logrus.WithError(err).Error("Failed to get user")
		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}

	return s.modelToProto(&user), nil
}

// ListUsers retrieves a list of users with pagination
func (s *UserService) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	logrus.WithFields(logrus.Fields{
		"page":      req.Page,
		"page_size": req.PageSize,
	}).Info("Listing users")

	page := req.Page
	pageSize := req.PageSize

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	skip := (page - 1) * pageSize
	limit := pageSize

	// Count total documents
	total, err := s.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		logrus.WithError(err).Error("Failed to count users")
		return nil, status.Errorf(codes.Internal, "failed to count users: %v", err)
	}

	// Find users with pagination
	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(limit))
	findOptions.SetSort(bson.D{{"created_at", -1}})

	cursor, err := s.collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		logrus.WithError(err).Error("Failed to find users")
		return nil, status.Errorf(codes.Internal, "failed to find users: %v", err)
	}
	defer cursor.Close(ctx)

	var users []*pb.User
	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			logrus.WithError(err).Error("Failed to decode user")
			continue
		}
		users = append(users, s.modelToProto(&user))
	}

	return &pb.ListUsersResponse{
		Users:    users,
		Total:    int32(total),
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// UpdateUser updates an existing user
func (s *UserService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.User, error) {
	logrus.WithField("id", req.Id).Info("Updating user")

	objectID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID: %v", err)
	}

	update := bson.M{
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	if req.Name != nil {
		update["$set"].(bson.M)["name"] = req.Name
	}
	if req.Email != nil {
		update["$set"].(bson.M)["email"] = req.Email
	}
	if req.Phone != nil {
		update["$set"].(bson.M)["phone"] = req.Phone
	}

	var user models.User
	err = s.collection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": objectID},
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		logrus.WithError(err).Error("Failed to update user")
		return nil, status.Errorf(codes.Internal, "failed to update user: %v", err)
	}

	return s.modelToProto(&user), nil
}

// DeleteUser deletes a user by ID
func (s *UserService) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*emptypb.Empty, error) {
	logrus.WithField("id", req.Id).Info("Deleting user")

	objectID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID: %v", err)
	}

	result, err := s.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		logrus.WithError(err).Error("Failed to delete user")
		return nil, status.Errorf(codes.Internal, "failed to delete user: %v", err)
	}

	if result.DeletedCount == 0 {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}

	return &emptypb.Empty{}, nil
}

// modelToProto converts a models.User to pb.User
func (s *UserService) modelToProto(user *models.User) *pb.User {
	return &pb.User{
		Id:        user.ID.Hex(),
		Name:      user.Name,
		Email:     user.Email,
		Phone:     user.Phone,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
}
