package main

import (
	"log"
	"net"

	"github.com/costinm/ssh-mesh/sshca-grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// This is a minimal gRPC server - signing SSH certificates. It is meant to run behind envoy, so no TLS or authz
// It does extract the identity from the XFCC header.
// ~9M stripped
func main() {

	sshs := &sshca_grpc.SSHSigner{
	}

	err := sshs.Init()
	if err != nil {
		panic(err)
	}

	servicePort := ":8080"
	lis, err := net.Listen("tcp", servicePort)
	if err != nil {
		log.Fatalf("net.Listen(tcp, %q) failed: %v", servicePort, err)
	}

	creds := insecure.NewCredentials()

	grpcOptions := []grpc.ServerOption{
		grpc.Creds(creds),
	}

	grpcServer := grpc.NewServer(grpcOptions...)
	sshca_grpc.RegisterSSHCertificateServiceServer(grpcServer, sshs)

	err = grpcServer.Serve(lis)
	if err != nil {
		panic(err)
	}
}

