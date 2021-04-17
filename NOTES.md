# Course notes

- gRPC solves communication between micro service architectures, which may be written in different programming languages
- Why not REST? REST is time consuming and hard to mantain. There's a full contract you need to agree on with your API consumers. There's endpoints, data models, error codes, latency, scalability, you name it
- gRPC is an open source framework (developed by Google) and allows you to define requests and responses for RPC (Remote Procedure Calls). It's fast, built on top of HTTP/2 and it's language independent
- gRPC Server handles all the RPC calls from gRPC Clients with proto requests/responses
- We define `.proto` files with our contracts and gRPC will generate the protocol buffers with the generated code we need for our gRPC clients to make RPC calls to the gRPC server