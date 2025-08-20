package main

import (
	protogenOnes "1c-grpc-gateway/pkg"
	"google.golang.org/protobuf/compiler/protogen"
)

// protoc -I="D:\protobuf\include" -I=".\proto" -I="D:\GIT\googleapis" --go_out=".\proto" example.proto custom.proto
// protoc -I="D:\protobuf\include" -I=".\proto" -I="D:\GIT\googleapis" --1c_out=D:\1 example.proto
// protoc -I="D:\protobuf\include" -I=".\proto" -I="D:\GIT\googleapis" -I="D:\GIT\grpc-gateway" --1c_out="swagger=1,logger=1:." example.proto

func main() {
	p := protogenOnes.NewPlugin()
	protogen.Options{}.Run(p.Process)
}
