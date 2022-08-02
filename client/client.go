package main

import (
        "context"
        "crypto/tls"
        "crypto/x509"
        "io/ioutil"
        "log"
        "time"

        pb "google.golang.org/grpc/examples/helloworld/helloworld"
        "google.golang.org/grpc"
        // "google.golang.org/grpc/credentials/insecure"
        "google.golang.org/grpc/credentials"
)

func main() {
        

        //here, loading in certificate authorities certificate. The process followed is that the client sends a cert to the server,
        //the server then uses the certificate authorities key to make sure that the cert sent is valid
        //this is used from the clients side this time to amke sure that the server is verified 
        caCert, err := ioutil.ReadFile("../certs/ca-cert.pem")
        if err != nil {
                log.Fatal(caCert)
        }

        //cert pool understanding reference https://golang.hotexamples.com/examples/crypto.x509/-/NewCertPool/golang-newcertpool-function-examples.html
        //append is used to add
        certPool := x509.NewCertPool()
        certPool.AppendCertsFromPEM(caCert)

        //reads in the client cert
        clientCert, err := tls.LoadX509KeyPair("../certs/client-cert.pem", "../certs/client-key.pem")
        if err != nil {
                log.Fatal(err)
        }

        configuration := &tls.Config{
                Certificates: []tls.Certificate{clientCert},
                RootCAs:      certPool,
        }

        tls := credentials.NewTLS(configuration)

        //establishes the client connection, only works because of the server and client -ext.conf files that specify an IP 
        connection, err := grpc.Dial(
                "0.0.0.0:9000",
                grpc.WithTransportCredentials(tls),
        )
        if err != nil {
                log.Fatal(err)
        }
        defer connection.Close()

        client := pb.NewGreeterClient(connection)
                log.Printf("Connection created!")

        //This is where the server is communicated with and response is printed if all worked
        //using context here to get the time to show it below 
        //for best practice with using the cancel function, it must get called at some point 
        ctx, cancel := context.WithTimeout(context.Background(), time.Second)
       //called here
        defer cancel()

        resp, err := client.SayHello(ctx, &pb.HelloRequest{Name: "Saoirse O'Donovan"})
        if err != nil {
                log.Fatal(err)
        }
        log.Printf("Communication secured using mTLS...")
        log.Printf("Greeting: %s", resp.GetMessage())
}
