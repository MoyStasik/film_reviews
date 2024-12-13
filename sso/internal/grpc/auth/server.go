package auth

import (
	"context"

	ssov1 "github.com/Lesha222/protos/gen/go/sso"
	"google.golang.org/grpc"
)

type serverAPI struct {
	ssov1.UnimplementedAuthServer
}

func Register(gRPC *grpc.Server) {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{})
}

func (s *serverAPI) Login(ctx context.Context, req *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	panic("inplement me ")
}

func (s *serverAPI) Register(ctx context.Context, req *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	panic("inplement me ")
}

func (s *serverAPI) IdAdmin(ctx context.Context, req *ssov1.AdminRequest) (*ssov1.AdminResponse, error) {
	panic("inplement me ")
}
