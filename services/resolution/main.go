package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net"

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

	queue, err := ch.QueueDeclare(
		"resolutionCompleted",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}

	msgs, err := ch.Consume(
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}

	go func() {
		var forever chan struct{}
		for msg := range msgs {
			payload := new(CompleteResolutionPayload)
			err := json.Unmarshal(msg.Body, payload)
			if err != nil {
				panic(err)
			}

			if err != nil {
				panic(err)
			}

			db.Model(
				&ResolutionModel{},
			).Where(
				"id = ?", payload.ResolutionId,
			).Update(
				"completed", true,
			)
		}
		<-forever
	}()

	listener, err := net.Listen("tcp", ":3002")
	if err != nil {
		panic(err)
	}

	conn, _ := grpc.Dial(":3001", grpc.WithTransportCredentials(insecure.NewCredentials()))

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
