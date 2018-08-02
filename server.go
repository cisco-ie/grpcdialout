package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"

	kafka "./kafka-producer"
	dialout "./mdt_dialout"
	telemetryBis "./telemetry_bis"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	peer "google.golang.org/grpc/peer"
)

type grpcLocalServer struct {
	// nothing yet
}

type dummyPeerType struct {
	// nothing yet
}

func (d *dummyPeerType) String() string {
	return "Unknown addr"
}

func (d *dummyPeerType) Network() string {
	return "Unknown net"
}

var dummyPeer dummyPeerType

func decrypt(data *dialout.MdtDialoutArgs) error {
	var err error
	ProtoItem := new(telemetryBis.Telemetry)
	err = proto.Unmarshal(data.Data, ProtoItem)
	if err != nil {
		return err
	}
	err = printer(ProtoItem)
	if err != nil {
		return err
	}

	return err
}

func printer(ProtoItem proto.Message) error {
	var err error
	var jsonpbObject jsonpb.Marshaler
	jsonString, err := jsonpbObject.MarshalToString(ProtoItem)
	if err != nil {
		return (err)
	}
	buf := new(bytes.Buffer)
	json.Indent(buf, []byte(jsonString), "", "  ")
	fmt.Println(buf)
	return err
}

func kafkaProducer(ProtoItem proto.Message) error {
	var err error
	data, err := json.Marshal(ProtoItem)
	if err != nil {
		return (err)
	}
	producer := kafka.NewProducer("TelemetryTest", []string{"172.31.96.200:9094"})
	producer.Produce(data)
	return err
}

func (s *grpcLocalServer) MdtDialout(stream dialout.GRPCMdtDialout_MdtDialoutServer) error {

	var endpoint *peer.Peer
	var ok bool

	if endpoint, ok = peer.FromContext(stream.Context()); !ok {
		endpoint = &peer.Peer{
			Addr: &dummyPeer,
		}
	}
	fmt.Printf("Receiving dialout stream from %s!\n", endpoint.Addr.String())

	for {
		var in *dialout.MdtDialoutArgs
		var err error

		in, err = stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		go decrypt(in)
		newerr := stream.Send(&dialout.MdtDialoutArgs{ReqId: in.ReqId})
		if newerr != nil {
			return newerr
		}
	}
}

func newServer() *grpcLocalServer {
	s := &grpcLocalServer{}
	return s
}

func main() {
	fmt.Printf("Hello, world.\n")

	lis, _ := net.Listen("tcp", ":57501")
	grpcServer := grpc.NewServer()

	dialout.RegisterGRPCMdtDialoutServer(grpcServer, newServer())

	grpcServer.Serve(lis)
}
