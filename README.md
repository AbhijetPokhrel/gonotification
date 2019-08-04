# gonotification
Example of simple notification in go lang using gRPC. The exaple illustrates the use of stream to send notification whithout closing it.

## Build

step 1. Create auto generated codes for server and client interfaces
```
make proto
```

step 2. Build the server

```
make build-server
```

step 3. Build the client

```
make build-client
```

Step 4. Run server
```
make run-server
```

Step 5. Run client

```
make run-client
```

## Working
Once the server and clients are connected type some message in server console and press enter the message should pop up on client console. You can connect as many client as you want to the same server. The notification goes to all of them.
