# phase-02
grpc client and server secure communication secured using mutual TLS (mTLS).

1. move to certs directory
2. run ./gen_certs.sh which will generate all of the keys and certificates 
3. create client-ext.conf and server-ext.conf files and add the IP so that a connection can be developed on the local side (subjectAltName=IP:0.0.0.0)
4. add -extfile server-ext.conf to line 30 and -extfile client-ext.conf to line 54 of the gen_certs.sh file otherwise we get this error when running 
"2022/08/02 10:20:19 rpc error: code = Unavailable desc = connection error: desc = "transport: authentication handshake failed: x509: cannot validate certificate for 0.0.0.0 because it doesn't contain any IP SANs"
exit status 1"
5. then open a split terminal, one within the server directory, and the other within the client directory
6. use the following command on the server side first: go run server.go
7. followed by: go run client.go on the client side and the connection should be established properly and a greeting should be printed
