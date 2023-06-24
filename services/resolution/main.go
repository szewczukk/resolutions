package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/szewczukk/resolution-service/proto"
	userProto "github.com/szewczukk/user-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type ResolutionModel struct {
	ID        int    `gorm:"primarykey"`
	Name      string `gorm:"unique"`
	UserId    int
	Completed bool
}

type CompleteResolutionPayload struct {
	UserId       int32 `json:"userId"`
	ResolutionId int32 `json:"resolutionId"`
}

type ResolutionServiceServer struct {
	proto.UnimplementedResolutionServiceServer
	Db                *gorm.DB
	UserServiceClient userProto.UserServiceClient
	Ch                *amqp.Channel
}

func NewResolutionServiceServer(
	db *gorm.DB,
	userServiceClient userProto.UserServiceClient,
	ch *amqp.Channel,
) *ResolutionServiceServer {
	return &ResolutionServiceServer{
		Db:                db,
		UserServiceClient: userServiceClient,
		Ch:                ch,
	}
}

func main() {
	db, err := gorm.Open(sqlite.Open("resolutions.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&ResolutionModel{})

	rabbitMqConnection, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}
	defer rabbitMqConnection.Close()

	ch, err := rabbitMqConnection.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	listener, err := net.Listen("tcp", ":3002")
	if err != nil {
		panic(err)
	}

	conn, _ := grpc.Dial(":3001", grpc.WithTransportCredentials(insecure.NewCredentials()))

	client := userProto.NewUserServiceClient(conn)

	grpcServer := grpc.NewServer()
	resolutionServiceServer := NewResolutionServiceServer(db, client, ch)
	proto.RegisterResolutionServiceServer(grpcServer, resolutionServiceServer)
	reflection.Register(grpcServer)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *ResolutionServiceServer) GetAllResolutions(
	ctx context.Context,
	request *proto.GetAllResolutionsRequest,
) (*proto.RepeatedResolutions, error) {
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

	return &proto.RepeatedResolutions{Resolutions: protoResolutions}, nil
}

func (s *ResolutionServiceServer) CreateResolution(
	ctx context.Context,
	request *proto.CreateResolutionRequest,
) (*proto.Resolution, error) {
	result, err := s.UserServiceClient.UserExists(
		context.Background(),
		&userProto.UserServiceUserId{UserId: request.UserId},
	)
	if err != nil {
		return nil, err
	}

	if !result.Exists {
		return nil, errors.New("user doesnt exist")
	}

	resolutionModel := ResolutionModel{
		Name:      request.Name,
		UserId:    int(request.UserId),
		Completed: false,
	}
	err = s.Db.Create(&resolutionModel).Error
	if err != nil {
		return nil, err
	}

	protoResolution := &proto.Resolution{
		Id:        int32(resolutionModel.ID),
		Name:      resolutionModel.Name,
		UserId:    int32(resolutionModel.UserId),
		Completed: resolutionModel.Completed,
	}

	return protoResolution, nil
}

func (s *ResolutionServiceServer) GetResolutionsByUserId(
	ctx context.Context,
	request *proto.UserId,
) (*proto.RepeatedResolutions, error) {
	var resolutionModels []ResolutionModel
	s.Db.Where("user_id = ?", request.UserId).Find(&resolutionModels)

	var protoResolutions []*proto.Resolution

	for _, resolution := range resolutionModels {
		protoResolutions = append(protoResolutions, &proto.Resolution{
			Id:        int32(resolution.ID),
			Name:      resolution.Name,
			UserId:    int32(resolution.UserId),
			Completed: resolution.Completed,
		})
	}

	return &proto.RepeatedResolutions{Resolutions: protoResolutions}, nil
}

func (s *ResolutionServiceServer) CompleteResolution(
	c context.Context,
	request *proto.CompleteResolutionRequest,
) (*proto.Resolution, error) {
	s.Db.Model(
		&ResolutionModel{},
	).Where(
		"id = ?", request.ResolutionId,
	).Update(
		"completed", true,
	)

	resolutionModel := new(ResolutionModel)
	err := s.Db.First(&resolutionModel, ResolutionModel{ID: int(request.ResolutionId)}).Error
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	payload := new(CompleteResolutionPayload)
	payload.ResolutionId = request.ResolutionId
	payload.UserId = int32(resolutionModel.UserId)

	serialized, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	err = s.Ch.PublishWithContext(
		ctx,
		"",
		"resolutionCompleted",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        serialized,
		},
	)
	if err != nil {
		return nil, err
	}

	return &proto.Resolution{
		Id:        int32(resolutionModel.ID),
		Name:      resolutionModel.Name,
		UserId:    int32(resolutionModel.UserId),
		Completed: resolutionModel.Completed,
	}, nil
}
