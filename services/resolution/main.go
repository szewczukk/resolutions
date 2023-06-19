package main

import (
	"context"
	"log"
	"net"

	"github.com/szewczukk/resolution-service/proto"
	"google.golang.org/grpc"
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
	Db *gorm.DB
}

func NewResolutionServiceServer(db *gorm.DB) *ResolutionServiceServer {
	return &ResolutionServiceServer{
		Db: db,
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

	grpcServer := grpc.NewServer()
	resolutionServiceServer := NewResolutionServiceServer(db)
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
	resolutionModel := ResolutionModel{Name: request.Name, UserId: int(request.UserId)}
	err := s.Db.Create(&resolutionModel).Error
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
