package main

import (
	"context"
	"flag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	desc "homework/pkg/pvz-service/v1"
	"log"
)

var (
	methodFlag = flag.String("method", "", "The method to call")
	dataFlag   = flag.String("data", "{}", "The data to send")
	hostFlag   = flag.String("host", "localhost:8080", "The host to connect to")
)

func main() {
	flag.Parse()

	conn, err := grpc.NewClient(
		*hostFlag,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	pvzService := desc.NewPvzServiceClient(conn)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var resp proto.Message
	switch *methodFlag {
	case "AcceptOrderDelivery":
		req := &desc.AcceptOrderDeliveryRequest{}
		if err := protojson.Unmarshal([]byte(*dataFlag), req); err != nil {
			log.Fatalf("failed to unmarshal data: %v", err)
		}
		resp, err = pvzService.AcceptOrderDelivery(ctx, req)
	case "AcceptReturn":
		req := &desc.AcceptReturnRequest{}
		if err := protojson.Unmarshal([]byte(*dataFlag), req); err != nil {
			log.Fatalf("failed to unmarshal data: %v", err)
		}
		resp, err = pvzService.AcceptReturn(ctx, req)
	case "GetOrders":
		req := &desc.GetOrdersRequest{}
		if err := protojson.Unmarshal([]byte(*dataFlag), req); err != nil {
			log.Fatalf("failed to unmarshal data: %v", err)
		}
		resp, err = pvzService.GetOrders(ctx, req)
	case "GetReturns":
		req := &desc.GetReturnsRequest{}
		if err := protojson.Unmarshal([]byte(*dataFlag), req); err != nil {
			log.Fatalf("failed to unmarshal data: %v", err)
		}
		resp, err = pvzService.GetReturns(ctx, req)
	case "GiveOrderToClient":
		req := &desc.GiveOrderToClientRequest{}
		if err := protojson.Unmarshal([]byte(*dataFlag), req); err != nil {
			log.Fatalf("failed to unmarshal data: %v", err)
		}
		resp, err = pvzService.GiveOrderToClient(ctx, req)
	case "ReturnOrderDelivery":
		req := &desc.ReturnOrderDeliveryRequest{}
		if err := protojson.Unmarshal([]byte(*dataFlag), req); err != nil {
			log.Fatalf("failed to unmarshal data: %v", err)
		}
		resp, err = pvzService.ReturnOrderDelivery(ctx, req)
	default:
		log.Fatalf("unknown method: %s", *methodFlag)
	}

	if err != nil {
		log.Fatalf("failed to call method: %v", err)
	}

	data, err := protojson.Marshal(resp)
	if err != nil {
		log.Fatalf("failed to marshal response: %v", err)
	}

	log.Printf("response: %s", data)
}
