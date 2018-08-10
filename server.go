package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	kafka "github.com/cisco-ie/grpcdialout/kafka-producer"
	dialout "github.com/cisco-ie/grpcdialout/mdt_dialout"
	telemetryBis "github.com/cisco-ie/grpcdialout/telemetry_bis"
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

func decrypt(data *dialout.MdtDialoutArgs) {
	ProtoItem := new(telemetryBis.Telemetry)
	err := proto.Unmarshal(data.Data, ProtoItem)
	if err != nil {
		log.Fatal(err)
	}
	if Configuration.Dump {
		if Configuration.Raw {
			go printer(data.Data)
		} else {
			var jsonpbObject jsonpb.Marshaler
			jsonString, err := jsonpbObject.MarshalToString(ProtoItem)
			if err != nil {
				log.Fatal(err)
			}
			buf := new(bytes.Buffer)
			json.Indent(buf, []byte(jsonString), "", "  ")
			go printer(buf.Bytes())
		}

	}
	if Configuration.Kafka.Brokers != nil {
		if Configuration.Raw {
			go kafkaProducer(data.Data, Configuration.Kafka.Topic, Configuration.Kafka.Brokers)
		} else {
			marshaled, err := json.Marshal(ProtoItem)
			if err != nil {
				log.Fatal(err)
			}
			go kafkaProducer(marshaled, Configuration.Kafka.Topic, Configuration.Kafka.Brokers)
		}
	}

}

func printer(data []byte) {
	f := File{Filename: Configuration.File}
	f.Write(data)
}

func kafkaProducer(data []byte, topic string, brokers []string) {
	producer := kafka.NewProducer(topic, brokers)
	producer.Produce(data)
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
	fmt.Printf("Starting gRPC Dialout Collector.\n")
	ConfigLoader()
	lis, _ := net.Listen("tcp", Configuration.Port)
	fmt.Printf("gRPC Server starting at: %s \n", Configuration.Port)
	grpcServer := grpc.NewServer()

	dialout.RegisterGRPCMdtDialoutServer(grpcServer, newServer())

	grpcServer.Serve(lis)
}
