package main

import (
	"fmt"
	"net"

	"github.com/fiuskylab/grpc-example/common"
	"github.com/fiuskylab/grpc-example/proto"
	"github.com/fiuskylab/grpc-example/services/auth/entity"
	"github.com/fiuskylab/grpc-example/services/auth/repository"
	"github.com/fiuskylab/grpc-example/services/auth/server"
	"github.com/fiuskylab/grpc-example/storage"
	"google.golang.org/grpc"
)

func main() {
	c := common.NewCommon()

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", c.Env.AUTH_PORT))
	if err != nil {
		c.Log.Error(err.Error())
		return
	}

	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)

	sto, err := storage.NewStorage(c)
	if err != nil {
		c.Log.Error(err.Error())
		return
	}

	if err = sto.PGSQL.AutoMigrate(entity.User{}); err != nil {
		c.Log.Error(err.Error())
		return
	}

	r := repository.NewAuthRepo(c, sto)
	sv, err := server.NewServer(c, r)

	if err != nil {
		c.Log.Error(err.Error())
		return
	}

	proto.RegisterAuthServiceServer(grpcServer, sv)

	if err = grpcServer.Serve(lis); err != nil {
		c.Log.Error(err.Error())
		return
	}
}
