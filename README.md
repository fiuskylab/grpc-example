# grpc-example
Simple gRPC implementation with Go

### Project Architecture
- [Whimsical](https://whimsical.com/6FVrohtLKkLL32dh9gYod4)

#### Folders

##### `proto`
- _auth.proto_
  - This is the proto file
  - It's in the project root because it is a common file for both services

##### `services`
- `auth/`
  - Auth handler service
- `api/`
  - Service available for users via REST api
  - Communicates with `auth` service via gRPC
