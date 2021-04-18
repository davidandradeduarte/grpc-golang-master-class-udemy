# Course notes

- gRPC solves communication between micro service architectures, which may be written in different programming languages
- Why not REST? REST is time consuming and hard to mantain. There's a full contract you need to agree on with your API consumers. There's endpoints, data models, error codes, latency, scalability, you name it
- gRPC is an open source framework (developed by Google) and allows you to define requests and responses for RPC (Remote Procedure Calls). It's fast, built on top of HTTP/2 and it's language independent
- gRPC Server handles all the RPC calls from gRPC Clients with proto requests/responses
- We define `.proto` files with our contracts and gRPC will generate the protocol buffers with the generated code we need for our gRPC clients to make RPC calls to the gRPC server
- Protocol Buffers are used to define the Messages, Service and RPC endpoints
    ```protobuf
        syntax = "proto3";

        message SomeRequestModel {
            string some_field = 1;
        }

        message SomeResponseModel {
            string some_field = 1;
        }

        service SomeService {
            rpc SomeMethod(SomeRequestModel) returns (SomeResponseModel) {};
        }
    ```
- Why not JSON - Protocol Buffers is much smaller (because it's binary) and saves bandwith. JSON is also more CPU intensive
- [gRPC](https://grpc.io) currently supports 11 programming languages
- gRPC has native implementations for Java, Go and C. Most of the other languages rely on the gRPC C implementation or implement it natively
- Usually HTTP/2 uses only one TCP connection and the server will push responses to the client when needed
- 4 types of API in gRPC: Unary, Server Streaming, Client Streaming and Bi Directional Streaming
    - **Unary**: Request --> Response
    - **Server Streaming**: Request --> Stream of Responses (server push)
    - **Client Streaming**: Stream of Requests --> Response
    - **Bi Directional Streaming**: Stream of Requests <--> Stream of Responses. It doesn't need to be ordered. It can be asynchronously. You can receive responses at the same time you're sending new requests
- gRPC leverages HTTP/2 multiplexing for async streaming communication
- gRPC servers are async by default
- gRPC clients can be async or sync
- gRPC clients can perform client side load-balancing
- gRPC Interceptors ...
- Packages for working with gRPC in go:
    ```bash
    go get -u google.golang.org/grpc
    go get -u github.com/golang/protobuf/protoc-gen-go
    ```
- Creating a new gRPC server in go:
    ```golang
    type server struct{}

    func main() {

        lis, err := net.Listen("tcp", "0.0.0.0:50051")

        if err != nil {
            log.Fatalf("failed to listen: %v", err)
        }

        s := grpc.NewServer()
        greetpb.RegisterGreetServiceServer(s, &server{})

        if err := s.Serve(lis); err != nil {
            log.Fatalf("failed to serve: %v", err)
        }
    }
    ```
- Creating a new gRPC client in go:
    ```golang
    func main() {
        cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

        if err != nil {
            log.Fatalf("could not connect: %v", err)
        }

        defer cc.Close()

        c := greetpb.NewGreetServiceClient(cc)

        fmt.Printf("created client: %f", c)
    }
    ```
- Server streaming API makes sense when the server needs to send a huge amount of data back, then diving the data into streams or chunks of data. It also makes a lot of sense to use server streaming when you need to push data to the clients, without the need of new requests. E.g live feed, chats, etc
- The `stream` keyword is used in proto files to specify that some request/response is going to be a stream of data
- 

## Additional notes

