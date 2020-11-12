//
// workers are shown here in-process
//

package main

import (
	zmq "github.com/pebbe/zmq4"

	"fmt"
	//"math/rand"
	"time"
)

func workerTask() {
	worker, _ := zmq.NewSocket(zmq.REQ)
	defer worker.Close()
	// set_id(worker)
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
		//request, _:=worker.RecvMessage(0)
		//a,_:= worker.GetMaxmsgsize()
		a, _ := zmq.GetMaxMsgsz()
		b, _ := worker.GetSndbuf()
		fmt.Println("Worker:", request, a, b)

		worker.Send(identity, zmq.SNDMORE)
		worker.Send("", zmq.SNDMORE)
		worker.Send("OK", 0)
	}
}

func main() {

	for true {
		//go worker_task()
		workerTask()
		time.Sleep(1000 * time.Millisecond)
	}
	time.Sleep(100 * time.Millisecond)
}
