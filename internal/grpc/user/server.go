package usergrpc

import (
	"context"
	"ecomUser/internal/domain/models"
	"ecomUser/internal/storage/postgres"
	"errors"

	user1 "github.com/KuranovNikita/ecomProto/gen/go/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserAuth interface {
	Login(ctx context.Context, userID int64, password string) (token string, err error)
	SaveUser(ctx context.Context, email string, login string, password string) (int64, error)
	GetUser(ctx context.Context, userID int64) (models.User, error)
	GetUserLogin(ctx context.Context, login string) (models.User, error)
}

type serverAPI struct {
	user1.UnimplementedUserServiceServer
	userAuth UserAuth
}

func Register(gRPCServer *grpc.Server, userAuth UserAuth) {
	user1.RegisterUserServiceServer(gRPCServer, &serverAPI{userAuth: userAuth})
}

func (s *serverAPI) Login(
	ctx context.Context,
	in *user1.LoginRequest,
) (*user1.LoginResponse, error) {
	if in.Login == "" {
		return nil, status.Error(codes.InvalidArgument, "login is required")
	}

	if in.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	user, err := s.userAuth.GetUserLogin(ctx, in.Login)

	token, err := s.userAuth.Login(ctx, user.ID, in.GetPassword())
	if err != nil {

		return nil, status.Error(codes.Internal, "failed to login")
	}

	return &user1.LoginResponse{Token: token}, nil
}

func (s *serverAPI) Register(
	ctx context.Context,
	in *user1.RegisterRequest,
) (*user1.RegisterResponse, error) {
	if in.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	if in.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	if in.Login == "" {
		return nil, status.Error(codes.InvalidArgument, "login is required")
	}

	uid, err := s.userAuth.SaveUser(ctx, in.GetEmail(), in.GetLogin(), in.GetPassword())
	if err != nil {
		if errors.Is(err, postgres.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}

		return nil, status.Error(codes.Internal, "failed to register user")
	}

	return &user1.RegisterResponse{UserId: uid}, nil
}

func (s *serverAPI) GetUser(ctx context.Context, in *user1.GetUserRequest) (*user1.GetUserResponse, error) {
	user, err := s.userAuth.GetUser(ctx, in.GetUserId())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get user")
	}

	grpcUserDetails := &user1.UserDetails{
		UserId: user.ID,
		Email:  user.Email,
		Login:  user.Login,
	}

	return &user1.GetUserResponse{UserDetails: grpcUserDetails}, nil
}
