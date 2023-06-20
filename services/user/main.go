package main

import (
	"context"
	"errors"
	"log"
	"net"

	"github.com/szewczukk/user-service/proto"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type UserModel struct {
	ID       int    `gorm:"primarykey"`
	Username string `gorm:"unique"`
	Password string
}

type UserServiceServer struct {
	proto.UnimplementedUserServiceServer
	Db *gorm.DB
}

func NewUserServiceServer(db *gorm.DB) *UserServiceServer {
	return &UserServiceServer{
		Db: db,
	}
}

func main() {
	db, err := gorm.Open(sqlite.Open("users.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&UserModel{})

	listener, err := net.Listen("tcp", ":3000")
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	userServiceServer := NewUserServiceServer(db)
	proto.RegisterUserServiceServer(grpcServer, userServiceServer)
	reflection.Register(grpcServer)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *UserServiceServer) UserExists(
	ctx context.Context,
	request *proto.UserExistsRequest,
) (*proto.UserExistsResponse, error) {
	userModel := UserModel{ID: int(request.Id)}
	err := s.Db.First(&userModel).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &proto.UserExistsResponse{Exists: false}, nil
		} else {
			return nil, err
		}
	}

	return &proto.UserExistsResponse{Exists: true}, nil
}

func (s *UserServiceServer) GetAllUsers(
	ctx context.Context,
	request *proto.GetAllUsersRequest,
) (*proto.GetAllUsersResponse, error) {
	var userModels []UserModel
	s.Db.Find(&userModels)

	var protoUsers []*proto.User

	for _, user := range userModels {
		protoUsers = append(protoUsers, &proto.User{Id: int32(user.ID), Username: user.Username})
	}

	return &proto.GetAllUsersResponse{Users: protoUsers}, nil
}

func (s *UserServiceServer) CreateUser(
	ctx context.Context,
	request *proto.UserCredentials,
) (*proto.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	payload := UserModel{Username: request.Username, Password: string(hashedPassword)}
	err = s.Db.Create(&payload).Error
	if err != nil {
		return nil, err
	}

	protoUser := &proto.User{Id: int32(payload.ID), Username: payload.Username}

	return protoUser, nil
}

func (s *UserServiceServer) AuthenticateUser(
	ctx context.Context,
	request *proto.UserCredentials,
) (*proto.AuthenticateUserResponse, error) {
	user := new(UserModel)
	err := s.Db.First(&user, UserModel{Username: request.Username}).Error
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		return nil, errors.New("wrong password")
	}

	return &proto.AuthenticateUserResponse{UserId: int32(user.ID)}, nil
}
