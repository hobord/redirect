//protoc -I ./redirect/ -I ../../../ --go_out=plugins=grpc:./redirect/ ./redirect/redirect.proto
// go get github.com/fullstorydev/grpcurl
// go install github.com/fullstorydev/grpcurl/cmd/grpcurl
/*
// docker run -p 6379:6379 --name redis-redisjson redislabs/rejson:latest
grpcurl.exe -plaintext localhost:50051 list

grpcurl -plaintext -d '{"url": "http://index.hu/path/subpath/?foo=bar&toremove=xyz&other=ok#bookmark"}' localhost:50051 github.com.hobord.redirect.RedirectService/GetRedirection

*/

package main

import (
	"fmt"
	"log"
	"net"
	"os"

	pb "github.com/hobord/redirect/redirect"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = ":50052"
	}
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Println("Server listen: ", port)

	s := grpc.NewServer()
	reflection.Register(s)

	srv := pb.CreateGrpcServer()
	pb.RegisterRedirectServiceServer(s, srv)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
