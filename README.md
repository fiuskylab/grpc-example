# grpc-example
Simple gRPC implementation with Go

## Summary
1. [References](#references)
2. [Project Architecture](#project-architecture)
  1. [Folders](#folders)
    1. [proto/](#proto)
    2. [services/](#services)

## References
- [Proto Guide](https://developers.google.com/protocol-buffers/docs/proto3)
- [gRPC Docs](https://grpc.io/docs/languages/go/quickstart/)

## Project Architecture
- [Whimsical](https://whimsical.com/6FVrohtLKkLL32dh9gYod4)

### Folders

#### `proto`
- _auth.proto_
  - This is the proto file
  - It's in the project root because it is a common file for both services

#### `services`
- `auth/`
  - Auth handler service
- `api/`
  - Service available for users via REST api
  - Communicates with `auth` service via gRPC
