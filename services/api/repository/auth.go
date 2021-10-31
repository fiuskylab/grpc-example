package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/fiuskylab/grpc-example/common"
	"github.com/fiuskylab/grpc-example/proto"
	"github.com/fiuskylab/grpc-example/services/api/entity"
	grpc "google.golang.org/grpc"
)

type AuthRepo interface {
	Create(u entity.User) (entity.UserResponse, error)
	Login(u entity.User) (entity.UserResponse, error)
	Check(token string) error
}

type AuthRepoCtx struct {
	Ctx        *common.Common
	AuthClient proto.AuthServiceClient
}

func NewAuthRepoCtx(c *common.Common) *AuthRepoCtx {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", c.Env.AUTH_PORT, c.Env.AUTH_PORT), grpc.WithInsecure())

	if err != nil {
		c.Log.Error(err.Error())
		return &AuthRepoCtx{}
	}

	client := proto.NewAuthServiceClient(conn)

	return &AuthRepoCtx{
		Ctx:        c,
		AuthClient: client,
	}
}

func (a *AuthRepoCtx) Create(u entity.User) (entity.UserResponse, error) {
	createReq := &proto.CreateUserRequest{
		Username: u.Username,
		Password: u.Password,
	}

	c, cancel := context.WithTimeout(context.Background(), time.Second*2)

	defer cancel()

	createRes, err := a.AuthClient.CreateUser(c, createReq)

	if err != nil {
		return entity.UserResponse{
			Error: err.Error(),
		}, err
	}

	userResp := entity.UserResponse{
		Token: createRes.Token,
		Error: createRes.ErrorMessage,
	}

	return userResp, nil
}

func (a *AuthRepoCtx) Login(u entity.User) (entity.UserResponse, error) {
	loginReq := &proto.LoginRequest{
		Username: u.Username,
		Password: u.Password,
	}

	c, cancel := context.WithTimeout(context.Background(), time.Second*2)

	defer cancel()

	loginRes, err := a.AuthClient.Login(c, loginReq)

	if err != nil {
		return entity.UserResponse{
			Error: err.Error(),
		}, err
	}

	userResp := entity.UserResponse{
		Token: loginRes.Token,
		Error: loginRes.ErrorMessage,
	}

	return userResp, nil
}

func (a *AuthRepoCtx) Check(token string) error {
	checkReq := &proto.CheckTokenRequest{
		Token: token,
	}

	c, cancel := context.WithTimeout(context.Background(), time.Second*2)

	defer cancel()

	checkRes, err := a.AuthClient.CheckToken(c, checkReq)

	if err != nil {
		return err
	}

	if checkRes.ErrorMessage != "" {
		return fmt.Errorf(checkRes.ErrorMessage)
	}

	return nil
}
