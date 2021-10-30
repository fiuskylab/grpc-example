package server

import (
	"context"

	"github.com/fiuskylab/grpc-example/common"
	"github.com/fiuskylab/grpc-example/proto"
	"github.com/fiuskylab/grpc-example/services/auth/repository"
)

type Server struct {
	proto.UnimplementedAuthServiceServer
	common *common.Common
	Repo   *repository.AuthRepoCtx
}

func NewServer(c *common.Common, r *repository.AuthRepoCtx) (*Server, error) {
	s := Server{
		common: c,
		Repo:   r,
	}

	return &s, nil
}

func (s *Server) CreateUser(ctx context.Context, req *proto.CreateUserRequest) (*proto.CreateUserResponse, error) {
	return s.Repo.Create(req), nil
}

func (s *Server) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	return s.Repo.Login(req), nil
}

func (s *Server) CheckToken(ctx context.Context, req *proto.CheckTokenRequest) (*proto.CheckTokenResponse, error) {
	return s.Repo.Check(req), nil
}
