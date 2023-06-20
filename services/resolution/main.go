package main

import (
	"context"
	"errors"
	"log"
	"net"

	"github.com/szewczukk/resolution-service/proto"
	userProto "github.com/szewczukk/user-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type ResolutionModel struct {
	ID     int    `gorm:"primarykey"`
	Name   string `gorm:"unique"`
	UserId int
}

type ResolutionServiceServer struct {
	proto.UnimplementedResolutionServiceServer
	Db                *gorm.DB
	UserServiceClient userProto.UserServiceClient
}

func NewResolutionServiceServer(
	db *gorm.DB,
	userServiceClient userProto.UserServiceClient,
) *ResolutionServiceServer {
	return &ResolutionServiceServer{
		Db:                db,
		UserServiceClient: userServiceClient,
	}
}

func main() {
	db, err := gorm.Open(sqlite.Open("resolutions.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&ResolutionModel{})

	listener, err := net.Listen("tcp", ":3001")
	if err != nil {
		panic(err)
	}

	conn, _ := grpc.Dial(":3000", grpc.WithTransportCredentials(insecure.NewCredentials()))

	client := userProto.NewUserServiceClient(conn)

	grpcServer := grpc.NewServer()
	resolutionServiceServer := NewResolutionServiceServer(db, client)
	proto.RegisterResolutionServiceServer(grpcServer, resolutionServiceServer)
	reflection.Register(grpcServer)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *ResolutionServiceServer) GetAllResolutions(
	ctx context.Context,
	request *proto.GetAllResolutionsRequest,
) (*proto.GetAllResolutionsResponse, error) {
	var resolutionModels []ResolutionModel
	s.Db.Find(&resolutionModels)

	var protoResolutions []*proto.Resolution

	for _, resolution := range resolutionModels {
		protoResolutions = append(protoResolutions, &proto.Resolution{
			Id:     int32(resolution.ID),
			Name:   resolution.Name,
			UserId: int32(resolution.UserId),
		})
	}

	return &proto.GetAllResolutionsResponse{Resolutions: protoResolutions}, nil
}

func (s *ResolutionServiceServer) CreateResolution(
	ctx context.Context,
	request *proto.CreateResolutionRequest,
) (*proto.Resolution, error) {
	result, err := s.UserServiceClient.UserExists(
		context.Background(),
		&userProto.UserExistsRequest{Id: request.UserId},
	)
	if err != nil {
		return nil, err
	}

	if !result.Exists {
		return nil, errors.New("user doesnt exist")
	}

	resolutionModel := ResolutionModel{Name: request.Name, UserId: int(request.UserId)}
	err = s.Db.Create(&resolutionModel).Error
	if err != nil {
		return nil, err
	}

	protoUser := &proto.Resolution{
		Id:     int32(resolutionModel.ID),
		Name:   resolutionModel.Name,
		UserId: int32(resolutionModel.UserId),
	}

	return protoUser, nil
}
