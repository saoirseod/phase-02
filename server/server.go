package main

import (
        "context"
        "crypto/tls"
        "crypto/x509"
        "io/ioutil"
        "log"
        "net"

        pb "google.golang.org/grpc/examples/helloworld/helloworld"
        "google.golang.org/grpc"
        "google.golang.org/grpc/credentials"
)

type server struct {
        pb.UnimplementedGreeterServer
}

//%v is the value in a default format when printing structs
//SayHello method implements helloworlds' GreeterServer
//built off of golang's grpc hello world example
func (s *server) SayHello(ctx context.Context, request *pb.HelloRequest) (*pb.HelloReply, error) {
        log.Printf("Name received: %v", request.GetName())
        return &pb.HelloReply{Message: "Hello " + request.GetName()}, nil
}

func main() {
        //assigns listening port
        lis, err := net.Listen("tcp", "0.0.0.0:9000")
        if err != nil {
                log.Fatalf("Oops! Failed to listen: %v", err)
        }

        
        //here, loading in certificate authorities certificate. The process followed is that the client sends a cert to the server,
        //the server then uses the certificate authorities key to make sure that the cert sent is valid
        caPem, err := ioutil.ReadFile("../certs/ca-cert.pem")
        if err != nil {
                log.Fatal(err)
        }

        //cert pool understanding reference https://golang.hotexamples.com/examples/crypto.x509/-/NewCertPool/golang-newcertpool-function-examples.html
        //cert pool made that adds certificate authorities cert, x509 is the type of cert
        //x509 certs generally contain a public key
        certPool := x509.NewCertPool()
        //.AppendCertsFromPEM is a golang function, it parses PEM encoded certificates and appends them (adds) to the pool
        certPool.AppendCertsFromPEM(caPem) 

        //here the .LoadX509KeyPair method is used to load the servers cert and key
        serverCert, err := tls.LoadX509KeyPair("../certs/server-cert.pem", "../certs/server-key.pem")
        if err != nil {
                log.Fatal(err)
        }

        // configuration of the certificate what we want to
        configuration := &tls.Config{
                //pulls in fetched server cert from the method above
                Certificates: []tls.Certificate{serverCert},
                //.RequireAndVerifyClientCert indicates that a client's cert should be requested upon handshake 
                ClientAuth:   tls.RequireAndVerifyClientCert,
                ClientCAs:    certPool,
        }

        //creates and holds all of the tls credentials
        tls := credentials.NewTLS(configuration)

        //initiates the grpc server with the stated grpc tls credentials 
        grpcServerSetup := grpc.NewServer(grpc.Creds(tls))

        //link back to code at the top, same as in hello-world go example 
        //the service is now registered in the server
        pb.RegisterGreeterServer(grpcServerSetup, &server{})

        //Print out the port being listened at so that we know its working 
        log.Printf("Currently listening at port %v", lis.Addr())
        if err := grpcServerSetup.Serve(lis); err != nil {
                log.Fatalf("Uh oh! There's an error with the grpc server: %v", err)
        }
}
