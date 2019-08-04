
SRC_DIR = notification_proto
DST_DIR = notification_proto
PROTO_FILE_NAME = notification.proto

# create auto generated codes for server and client interfaces
proto:
	protoc -I $(SRC_DIR)/ $(SRC_DIR)/$(PROTO_FILE_NAME) --go_out=plugins=grpc:$(DST_DIR)

# build the server code
build-server:
	go build gonotification/notification/server

# build the client code
build-client:
	go build gonotification/notification/client

# run the server
run-server:
	./server

# run the client
run-client:
	./client

