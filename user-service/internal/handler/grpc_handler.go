package handler

import (
	"context"	

	"E-Commerce/user-service/internal/auth"
	"E-Commerce/user-service/internal/entity"
	"E-Commerce/user-service/internal/repository"
	pb "E-Commerce/user-service/proto"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) RegisterUser(ctx context.Context, req *pb.RegisterUserRequest) (*pb.AuthResponse, error) {
	existing, err := s.repo.GetUserByEmail(req.Email)
	if err == nil && existing != nil {
		return nil, status.Errorf(codes.AlreadyExists, "user already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password")
	}

	id := uuid.New().String()

	role := req.Role
	if role == "" {
		role = "user"
	}
	if role != "admin" && role != "user" {
		return nil, status.Errorf(codes.InvalidArgument, "invalid role")
	}

	user := &models.User{
		ID:       id,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     role,
	}
	err = s.repo.CreateUser(user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user")
	}

	token, err := auth.GenerateToken(user.ID, user.Role)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate token")
	}

	return &pb.AuthResponse{
		Token: token,
		User: &pb.User{
			Id:    user.ID,
			Email: user.Email,
			Role:  user.Role,
		},
	}, nil
}

func (s *UserService) AuthenticateUser(ctx context.Context, req *pb.AuthenticateUserRequest) (*pb.AuthResponse, error) {
	user, err := s.repo.GetUserByEmail(req.Email)
	if err != nil || user == nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid credentials")
	}

	token, err := auth.GenerateToken(user.ID, user.Role)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate token")
	}

	return &pb.AuthResponse{
		Token: token,
		User: &pb.User{
			Id:    user.ID,
			Email: user.Email,
			Role:  user.Role,
		},
	}, nil
}

func (s *UserService) GetUserProfile(ctx context.Context, req *pb.GetUserProfileRequest) (*pb.UserProfile, error) {
	user, err := s.repo.GetUserByID(req.UserId)
	if err != nil || user == nil {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}

	return &pb.UserProfile{
		Id:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	}, nil
}

func (s *UserService) UpdateUserProfile(ctx context.Context, req *pb.UpdateUserProfileRequest) (*pb.UserProfile, error) {
	user, err := s.repo.GetUserByID(req.UserId)
	if err != nil || user == nil {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password")
	}

	user.Password = string(hashedPassword)
	err = s.repo.UpdateUser(user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update user")
	}

	return &pb.UserProfile{
		Id:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	}, nil
}