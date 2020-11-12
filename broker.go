//
// broker.
//

package main

import (
	"fmt"
	zmq "github.com/pebbe/zmq4"
	"log"
	"os"
	"time"
)

func writeFile(content string) {
	filename := "access.log"
	info, err := os.Stat(filename)
	fileSize := info.Size()
	fmt.Println(fileSize)
	if fileSize > 1000000 {
		newName := filename + "_" + time.Now().String()
		err = os.Rename(filename, newName)
		if err != nil {
			log.Fatal(err)
		}
	}

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := f.Write([]byte(content)); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}

}

func main() {
	//  Prepare our sockets
	frontend, _ := zmq.NewSocket(zmq.ROUTER)
	backend, _ := zmq.NewSocket(zmq.ROUTER)
	defer frontend.Close()
	defer backend.Close()
	frontend.Bind("ipc://frontend.ipc")
	backend.Bind("ipc://backend.ipc")

	//  Queue of available workers
	worker_queue := make([]string, 0, 10)

	poller1 := zmq.NewPoller()
	poller1.Add(backend, zmq.POLLIN)
	poller2 := zmq.NewPoller()
	poller2.Add(backend, zmq.POLLIN)
	poller2.Add(frontend, zmq.POLLIN)

	for true {
		//  Poll frontend only if we have available workers
		var sockets []zmq.Polled
		if len(worker_queue) > 0 {
			sockets, _ = poller2.Poll(-1)
		} else {
			sockets, _ = poller1.Poll(-1)
		}
		for _, socket := range sockets {
			switch socket.Socket {
			case backend:

				//  Handle worker activity on backend
				//  Queue worker identity for load-balancing
				worker_id, _ := backend.Recv(0)
				//if !(len(worker_queue) < NBR_WORKERS) {
				//	panic("!(len(worker_queue) < NBR_WORKERS)")
				//}
				worker_queue = append(worker_queue, worker_id)

				//  Second frame is empty
				empty, _ := backend.Recv(0)
				if empty != "" {
					panic(fmt.Sprintf("empty is not \"\": %q", empty))
				}

				//  Third frame is READY or else a client reply identity
				client_id, _ := backend.Recv(0)

				//  If client reply, send rest back to frontend
				if client_id != "READY" {
					empty, _ := backend.Recv(0)
					if empty != "" {
						panic(fmt.Sprintf("empty is not \"\": %q", empty))
					}
					reply, _ := backend.Recv(0)
					frontend.Send(client_id, zmq.SNDMORE)
					frontend.Send("", zmq.SNDMORE)
					frontend.Send(reply, 0)
					//client_nbr--
				}

			case frontend:
				//  Here is how we handle a client request:

				//  Now get next client request, route to last-used worker
				//  Client request is [identity][empty][request]
				client_id, _ := frontend.Recv(0)
				empty, _ := frontend.Recv(0)
				if empty != "" {
					panic(fmt.Sprintf("empty is not \"\": %q", empty))
				}

				request, _ := frontend.Recv(0)
				fmt.Println("broker:", client_id, request)
				go writeFile(request)

				backend.Send(worker_queue[0], zmq.SNDMORE)
				backend.Send("", zmq.SNDMORE)
				backend.Send(client_id, zmq.SNDMORE)
				backend.Send("", zmq.SNDMORE)
				backend.Send(request, 0)

				//  Dequeue and drop the next worker identity
				worker_queue = worker_queue[1:]

			}
		}
	}

	time.Sleep(100 * time.Millisecond)
}
