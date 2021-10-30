package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/fiuskylab/grpc-example/common"
	"github.com/fiuskylab/grpc-example/proto"
	"github.com/fiuskylab/grpc-example/services/auth/entity"
	"github.com/fiuskylab/grpc-example/services/auth/jwt"
	"github.com/fiuskylab/grpc-example/storage"
	"golang.org/x/crypto/bcrypt"
)

type AuthRepo interface {
	Create(req *proto.CreateUserRequest) *proto.CreateUserResponse
	Login(req *proto.LoginRequest) *proto.LoginResponse
	Check(req *proto.CheckTokenRequest) *proto.CheckTokenResponse
}

type AuthRepoCtx struct {
	Common  *common.Common
	JWT     *jwt.JWT
	Storage *storage.Storage
}

func NewAuthRepo(c *common.Common, sto *storage.Storage) *AuthRepoCtx {
	a := AuthRepoCtx{
		Common:  c,
		Storage: sto,
		JWT:     jwt.NewJWT(c, sto),
	}
	return &a
}

const (
	fieldsNotMatch = `'%s' and '%s' don't match`
	emptyField     = `field '%s' can not be empty`
	wrongPassord   = `wrong password`
)

func (a *AuthRepoCtx) Create(req *proto.CreateUserRequest) *proto.CreateUserResponse {
	if req.Username == "" {
		return &proto.CreateUserResponse{
			ErrorMessage: fmt.Sprintf(emptyField, "username"),
		}
	}
	if req.Password == "" {
		return &proto.CreateUserResponse{
			ErrorMessage: fmt.Sprintf(emptyField, "password"),
		}
	}

	encPW, err := a.EncPassword(req.Password)

	if err != nil {
		a.Common.Log.Error(err.Error())
		return &proto.CreateUserResponse{
			ErrorMessage: err.Error(),
		}
	}

	u := entity.User{
		Username: req.Username,
		Password: string(encPW),
	}

	if err = a.Storage.PGSQL.Create(&u).Error; err != nil {
		a.Common.Log.Error(err.Error())
		return &proto.CreateUserResponse{
			ErrorMessage: err.Error(),
		}
	}
	token, err := a.JWT.NewToken(u.ID.String())

	if err != nil {
		a.Common.Log.Error(err.Error())
		return &proto.CreateUserResponse{
			ErrorMessage: err.Error(),
		}
	}

	if res := a.Storage.Redis.Set(context.Background(), token, u.ID.String(), time.Hour*8); res.Err() != nil {
		a.Common.Log.Error(res.Err().Error())
		return &proto.CreateUserResponse{
			ErrorMessage: res.Err().Error(),
		}
	}

	return &proto.CreateUserResponse{
		Token: token,
	}
}

func (a *AuthRepoCtx) Login(req *proto.LoginRequest) *proto.LoginResponse {
	if req.Username == "" {
		return &proto.LoginResponse{
			ErrorMessage: fmt.Sprintf(emptyField, "username"),
		}
	}

	if req.Password == "" {
		return &proto.LoginResponse{
			ErrorMessage: fmt.Sprintf(emptyField, "password"),
		}
	}

	u := new(entity.User)

	if err := a.Storage.PGSQL.
		First(u, "username = ?", req.Username).
		Error; err != nil {
		return &proto.LoginResponse{
			ErrorMessage: err.Error(),
		}
	}

	if !a.CheckPassword(u.Password, req.Password) {
		return &proto.LoginResponse{
			ErrorMessage: wrongPassord,
		}
	}

	token, err := a.JWT.NewToken(u.ID.String())

	if err != nil {
		return &proto.LoginResponse{
			ErrorMessage: err.Error(),
		}
	}

	if res := a.Storage.Redis.Set(context.Background(), token, u.ID.String(), time.Hour*8); res.Err() != nil {
		a.Common.Log.Error(res.Err().Error())
		return &proto.LoginResponse{
			ErrorMessage: res.Err().Error(),
		}
	}

	return &proto.LoginResponse{
		Token: token,
	}
}

func (a *AuthRepoCtx) Check(req *proto.CheckTokenRequest) *proto.CheckTokenResponse {
	var errMsg string
	if err := a.JWT.CheckToken(req.Token); err != nil {
		errMsg = err.Error()
	}

	res := a.Storage.Redis.Get(context.Background(), req.Token)

	if err := res.Err(); err != nil {
		return &proto.CheckTokenResponse{
			ErrorMessage: errMsg,
		}
	}

	return &proto.CheckTokenResponse{
		Id:           res.Val(),
		ErrorMessage: errMsg,
	}
}

func (a *AuthRepoCtx) EncPassword(pw string) (string, error) {
	encPW, err := bcrypt.GenerateFromPassword([]byte(pw), 14)
	return string(encPW), err
}

func (a *AuthRepoCtx) CheckPassword(encPW, pw string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(encPW), []byte(pw))
	return err == nil
}
