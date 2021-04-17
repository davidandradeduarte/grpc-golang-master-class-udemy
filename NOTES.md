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
    1. **Unary**: Request --> Response
    1. **Server Streaming**: Request --> Stream of Responses (server push)
    1. **Client Streaming**: Stream of Requests --> Response
    1. **Bi Directional Streaming**: Stream of Requests <--> Stream of Responses. It doesn't need to be ordered. It can be asynchronously. You can receive responses at the same time you're sending new requests
- gRPC leverages HTTP/2 multiplexing for async streaming communication
