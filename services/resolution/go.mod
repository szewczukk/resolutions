module github.com/szewczukk/resolution-service

go 1.20

require (
	github.com/rabbitmq/amqp091-go v1.8.1
	github.com/szewczukk/user-service v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.56.0
	google.golang.org/protobuf v1.30.0
	gorm.io/driver/sqlite v1.5.2
	gorm.io/gorm v1.25.2-0.20230530020048-26663ab9bf55
)

require (
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/mattn/go-sqlite3 v1.14.17 // indirect
	golang.org/x/net v0.10.0 // indirect
	golang.org/x/sys v0.9.0 // indirect
	golang.org/x/text v0.10.0 // indirect
	google.golang.org/genproto v0.0.0-20230410155749-daa745c078e1 // indirect
)

replace github.com/szewczukk/user-service => ../user
