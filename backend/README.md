grpcurl -plaintext -d '{}' localhost:50051 user.UserService/ListUsers

protoc --go_out=services/user-service/api/gen --go-grpc_out=services/user-service/api/gen proto/user.proto

protoc --go_out=services/course-service/api/gen --go-grpc_out=services/course-service/api/gen proto/course.proto

docker build -t amrimuf/course-service:latest .     