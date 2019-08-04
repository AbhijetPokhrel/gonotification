// Package main implements a client for Notification service
package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	pb "gonotification/notification_proto"

	"google.golang.org/grpc"
)

const (
	address = "localhost:5001"
)

func main() {

	reader := bufio.NewReader(os.Stdin)

	clientName, _ := stdin(reader, "Enter client name : ")

	ageStr, _ := stdin(reader, "Enter client age : ")
	age, err := strconv.Atoi(ageStr)
	if err != nil {
		fmt.Print(err)
		log.Fatal("Enter valid age %v", err)
	}

	address, _ := stdin(reader, "Enter client address : ")

	clientDetails := &pb.ClientDetail{
		Name:    clientName,
		Age:     int32(age),
		Address: address,
	}

	connectToServer(clientDetails)
}

func connectToServer(clientDetails *pb.ClientDetail) {

	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	client := pb.NewNotificationClient(conn)

	stream, err := client.ConnectToServer(context.Background(), clientDetails)

	for {

		// listen for streams
		notificationMessage, err := stream.Recv()

		if err == io.EOF { //no more stream to listen
			break
		}
		if err != nil { // some error occured
			log.Fatal("%v", err)
		}
		onNewNotification(notificationMessage)
	}
}

// called when the notification arrives
func onNewNotification(notificationMessage *pb.NotificationMessage) {
	fmt.Printf("%d: New message : %s", notificationMessage.GetTime()/1e6, notificationMessage.Message)
}

func stdin(reader *bufio.Reader, query string) (string, error) {
	fmt.Println(query)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", nil
	}
	return input[:len(input)-1], nil
}
