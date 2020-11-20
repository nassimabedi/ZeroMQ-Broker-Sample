//
// worker get package from broker
//

package main

import (
	"fmt"
	zmq "github.com/pebbe/zmq4"
	"log"
	"sync"
	"time"
)

const (
	ReceivedFilename = "received.txt" // keep received package count number in this file
)

func workerTask(wg *sync.WaitGroup, mReceived *sync.Mutex) {
	time.Sleep(1 * time.Second)
	worker, _ := zmq.NewSocket(zmq.REQ)
	defer worker.Close()
	worker.Connect("ipc://backend.ipc")

	//  Tell broker we're ready for work
	worker.Send("READY", 0)

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
		log.Println("Worker:", request)

		received, err := readAndWrite(ReceivedFilename, mReceived)
		if err != nil {
			log.Println("Error is : ", err)
		}

		real, err := readFile(RealFilename)
		if err != nil {
			log.Println("Error is : ", err)
		}

		log.Println("real number send by client :", real)
		log.Println("received id :", received)

		worker.Send(identity, zmq.SNDMORE)
		worker.Send("", zmq.SNDMORE)
		worker.Send("OK", 0)
	}
	wg.Done()
}

func main() {
	var wg sync.WaitGroup
	var mReceived sync.Mutex
	for true {
		wg.Add(1)
		go workerTask(&wg, &mReceived)
		time.Sleep(1000 * time.Millisecond)
	}
}
