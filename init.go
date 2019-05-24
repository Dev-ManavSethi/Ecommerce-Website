package main

import (
	"log"
	"net"
	"text/template"

	"github.com/DevManavSethi/EcommerceWebsite/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct{}

func init() {
	var err001 error

	//------------------------------------------------------------------------

	Tpl, err001 = template.ParseGlob("./templates/*")
	if err001 != nil {
		FatalOnError("Error parsing glob templates", err001)
		return
	}

	//-------------------------------------------------------------------

	go func() {
		log.Println("Starting gRPC Server")

		lis, err := net.Listen("tcp", "0.0.0.0:50051")
		if err != nil {
			log.Fatalf("Failed to listen: %v", err)
		}

		s := grpc.NewServer()
		service.RegisterEcommerceServer(s, &server{})

		reflection.Register(s)

		log.Println("gRPC Server started!")

		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
}
