# ZeroMQ-Broker-Sample


A simple sample of client,broker and worker

**client** : create package and send to broker

**broker**: save package to file and send to worker

**worker**: get package from broker and display count of received and real package number


Pre installation:

`go mod vendor`

to run:

`go run client.go handler.go`

`go run broker.go handler.go`

`go run worker.go handler.go`
