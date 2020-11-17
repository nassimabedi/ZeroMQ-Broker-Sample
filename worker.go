//
// workers
//

package main

import (
	"fmt"
	zmq "github.com/pebbe/zmq4"
	"log"
	"strconv"
	"time"
)

func workerTask() {
	worker, _ := zmq.NewSocket(zmq.REQ)
	defer worker.Close()
	worker.Connect("ipc://backend.ipc")

	//  Tell broker we're ready for work
	worker.Send("READY", 0)

	//receivedInt := 0
	var receivedInt int
	value, err := GetValue("received")
	if err != nil {
		log.Println(err)
	}
	if value == "" {
		receivedInt = 0
	}

	if value != "" {
		receivedInt, err = strconv.Atoi(value)
		if err != nil {
			log.Println("There is an Error: ", err)

		}
	}

	for {
		//  Read and save all frames until we get an empty frame
		//  In this example there is only 1 but it could be more
		identity, _ := worker.Recv(0)
		empty, _ := worker.Recv(0)
		if empty != "" {
			panic(fmt.Sprintf("empty is not \"\": %q", empty))
		}

		//  Get request, send reply
		request, _ := worker.Recv(0)

		id, _ := worker.Recv(0)

		log.Println("Worker:", request)

		receivedInt = receivedInt + 1

		//save received id to redis
		err = SetValue("received", receivedInt)
		log.Println(err)

		log.Println("real id :", id)
		log.Println("received id :", receivedInt)

		worker.Send(identity, zmq.SNDMORE)
		worker.Send("", zmq.SNDMORE)
		worker.Send("OK", 0)
	}
}

func main() {

	for true {
		go workerTask()
		time.Sleep(1000 * time.Millisecond)
	}
}
