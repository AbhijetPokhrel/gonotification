// Package main implements a server for Notification service
package main

import (
	"bufio"
	"fmt"
	pb "gonotification/notification_proto"
	"log"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
)

const (
	address = "localhost:5001"
)

// Server implement the  Notification service
// clients is the map of connected clients
// clientStreams is the map of connected clients stream
type Server struct {
	clients       map[string]*pb.ClientDetail
	clientStreams map[string]*pb.Notification_ConnectToServerServer
}

func (server *Server) init() {
	server.clients = make(map[string]*pb.ClientDetail)
	server.clientStreams = make(map[string]*pb.Notification_ConnectToServerServer)
}

// ConnectToServer is called when clietn make connection to the server
// this function will add the client to the servers clienst map and stores the client stream
// the stream should not be killed so we do not return from this server
// for this purpose the infinite loop is used
func (server *Server) ConnectToServer(in *pb.ClientDetail, stream pb.Notification_ConnectToServerServer) error {
	server.addNewClient(in, &stream)
	// loop infinitely to keep stream alive
	// else this stream will be closed
	for {
	}
	return nil
}

// adds new client to map
func (server *Server) addNewClient(in *pb.ClientDetail, stream *pb.Notification_ConnectToServerServer) {
	log.Printf("adding new client")
	server.clientStreams[in.Name] = stream
	server.clients[in.Name] = in
}

// send notification to specific client
func (server *Server) sendNotification(clientID string, msg string) {
	client := server.clients[clientID]
	stream := server.clientStreams[clientID]

	notificationMessage := &pb.NotificationMessage{
		Message: fmt.Sprintf("%s(age : %d) currently living in %s :: %s", client.Name, client.Age, client.Address, msg),
		Time:    time.Now().UnixNano(),
	}

	(*stream).Send(notificationMessage)

}

func main() {

	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	server := &Server{}
	server.init()
	options := []grpc.ServerOption{}
	options = append(options, grpc.MaxMsgSize(100*1024*1024))
	options = append(options, grpc.MaxRecvMsgSize(100*1024*1024))
	s := grpc.NewServer(options...)

	pb.RegisterNotificationServer(s, server)
	// go routine to get server notification message from stdin
	go waitForMessage(server)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}

func waitForMessage(server *Server) {
	for { // get the server notification message
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Notification Msg : ")
		msg, _ := reader.ReadString('\n')
		// send the message to all the clients
		for clientID := range server.clients {
			server.sendNotification(clientID, msg)
		}
	}

}
