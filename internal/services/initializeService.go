package services

import (
	"log"
	"net"

	"github.com/akshay0074700747/projectandCompany_management_protofiles/pb/companypb"
	"google.golang.org/grpc"
)

type CompanyEngine struct {
	Srv companypb.CompanyServiceServer
}

func NewCompanyEngine(srv companypb.CompanyServiceServer) *CompanyEngine {
	return &CompanyEngine{
		Srv: srv,
	}
}
func (engine *CompanyEngine) Start(addr string) {

	server := grpc.NewServer(
	// grpc.UnaryInterceptor(intersceptors.UnaryInterscaptor),
	)
	companypb.RegisterCompanyServiceServer(server, engine.Srv)

	listener, err := net.Listen("tcp", addr)

	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", addr, err)
	}

	log.Printf("Company Server is listening...")

	if err = server.Serve(listener); err != nil {
		log.Fatalf("Failed to listen on port %s: %v", addr, err)
	}

}
