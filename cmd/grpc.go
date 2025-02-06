package cmd

import (
	"ewallet-wallet/helpers"
	"log"
	"net"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func ServeGRPC() {
	// init dependency
	// dependency := dependencyInject()

	s := grpc.NewServer()
	// list method
	// tokenvalidation.RegisterTokenValidationServer(s, dependency.TokenValidationAPI)

	lis, err := net.Listen("tcp", ":"+helpers.GetEnv("GRPC_PORT", "7000"))
	if err != nil {
		log.Fatal("failed to listen grpc port: ", err)
	}

	logrus.Info("start listening grpc on port:" + helpers.GetEnv("GRPC_PORT", "7000"))
	if err := s.Serve(lis); err != nil {
		log.Fatal("failed to serve grpc port: ", err)
	}

}
