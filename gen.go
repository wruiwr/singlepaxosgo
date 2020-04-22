package singlepaxos

// Invoke protoc in order to compile our protobuf definition:
//go:generate protoc -I=$GOPATH/src/:. --gorums_out=plugins=grpc+gorums:. singlepaxos.proto
