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
    - **Client Streaming**: Stream of Requests --> Response (client push)
    - **Bi Directional Streaming**: Stream of Requests <--> Stream of Responses. It doesn't need to be ordered. It can be asynchronously. You can receive responses at the same time you're sending new requests (client and server push)
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
- The `stream` keyword is used in proto files to specify that some request/response is going to be a stream of data
- Server streaming API makes sense when the server needs to send a huge amount of data to the clients, then dividing the data into streams or chunks of data. It also makes a lot of sense to use server streaming when you need to push data to the clients, without the need of new requests. E.g live feed, chats, etc
- Client streaming API makes sense when the client needs to send a huge amount of data to the server, then dividing the data into streams or chunks of data. Also when the server processing may be expensive and should happen as the client sends data or when the client needs to push data to the server without really expecting a response
- Bi directional streaming API makes sense when both the client and the server need to send a huge amount of data between them. The number of responses doesn't need to match the number of requests. Can make sense for example in any "chat" protocol or long running connections
- gRPC uses its own error codes: [https://grpc.io/docs/guides/error/](https://grpc.io/docs/guides/error/)
- Deadlines define how long the client will wait for a gRPC call to complete. If we reach the defined deadline we will get a DEADLINE_EXCEEDED error code. It's recommended to set a deadline for each RPC call [https://grpc.io/blog/deadlines/](https://grpc.io/blog/deadlines/)
- Deadline should also be passed (chained) through gRPC services. E.g if client A calls gRPC service B and service B calls gRPC service C, C should know about the initial deadline defined by client A
- gRPC SSL/TLS or token based auth can be found here: [https://grpc.io/docs/guides/auth/](https://grpc.io/docs/guides/auth/)
- Enabling reflection on a gRPC server is very useful for API discovery. We can set reflection with
    ```golang
    import "google.golang.org/grpc/reflection"
    ...
    reflection.Register(s)
    ```
    see [https://github.com/grpc/grpc-go/tree/master/reflection](https://github.com/grpc/grpc-go/tree/master/reflection)
- Reflection can be used to create gRPC clients without the need ofa protobuf. For example, [evans](https://github.com/ktr0731/evans) it's a pretty cool CLI that uses reflection to talk to a gRPC server as a gRPC client. It can also use .proto files and other goodies. I'll be definitely using it
- we can capture os signals for user interaction (e.g CTRL+C to exit)
    ```golang
    ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	<-ch
    os.Exit(0)
    ```
- Examples from Google's protos have a lot of comments documenting the rpc error codes and behaviours. This is very important to the server who implements the proto but specially for clients that consume it, so that they can use it properly. e.g [pubsub.proto](https://github.com/googleapis/googleapis/blob/master/google/pubsub/v1/pubsub.proto), [spanner.proto](https://github.com/googleapis/googleapis/blob/master/google/spanner/v1/spanner.proto)

# Questions

- Even if we set a context deadline in our client, it seems like we need to hanle it on the server. I was hoping it to be more straight forward - e.g the server would return `GRPC_STATUS_DEADLINE_EXCEEDED` automatically once the context deadline was exceeded. Also, I didn't understood the difference between `context.WithDeadline` and `context.WithTimeout`
- What's the difference between calling the `type.GetField()` method from gRPC generated code, isntead of the field name `type.Field` (`Field` name used just as an example to any other message field)
