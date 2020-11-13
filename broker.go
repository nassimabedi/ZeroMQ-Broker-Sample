//
// broker.
//

package main

import (
	"fmt"
	zmq "github.com/pebbe/zmq4"
	"log"
	"os"
	"strconv"
	"time"
)

const (
	Filename    = "access.log" // name of the file to write file
	MaxFileSize = 1000000      //mac file size to write to another file
)

// write message to file if file size bigger than MaxFileSize Older file rename
func writeFile(content string) {
	filename := Filename
	info, err := os.Stat(filename)
	fileSize := info.Size()
	//  if filesize mare than MaxFileSize the file rename
	if fileSize > MaxFileSize {
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
	frontend, _ := zmq.NewSocket(zmq.ROUTER)
	backend, _ := zmq.NewSocket(zmq.ROUTER)
	defer frontend.Close()
	defer backend.Close()
	frontend.Bind("ipc://frontend.ipc")
	backend.Bind("ipc://backend.ipc")

	workerQueue := make([]string, 0, 10)

	poller1 := zmq.NewPoller()
	poller1.Add(backend, zmq.POLLIN)
	poller2 := zmq.NewPoller()
	poller2.Add(backend, zmq.POLLIN)
	poller2.Add(frontend, zmq.POLLIN)

	for true {
		//  Poll frontend only if we have available workers
		var sockets []zmq.Polled
		if len(workerQueue) > 0 {
			sockets, _ = poller2.Poll(-1)
		} else {
			sockets, _ = poller1.Poll(-1)
		}
		for _, socket := range sockets {
			switch socket.Socket {
			case backend:

				//  Handle worker activity on backend
				//  Queue worker identity for load-balancing
				workerId, _ := backend.Recv(0)
				workerQueue = append(workerQueue, workerId)

				//  Second frame is empty
				empty, _ := backend.Recv(0)
				if empty != "" {
					panic(fmt.Sprintf("empty is not \"\": %q", empty))
				}

				//  Third frame is READY or else a client reply identity
				clientId, _ := backend.Recv(0)

				//  If client reply, send rest back to frontend
				if clientId != "READY" {
					empty, _ := backend.Recv(0)
					if empty != "" {
						panic(fmt.Sprintf("empty is not \"\": %q", empty))
					}
					reply, _ := backend.Recv(0)
					frontend.Send(clientId, zmq.SNDMORE)
					frontend.Send("", zmq.SNDMORE)
					frontend.Send(reply, 0)
				}

			case frontend:
				clientId, _ := frontend.Recv(0)
				empty, _ := frontend.Recv(0)
				if empty != "" {
					panic(fmt.Sprintf("empty is not \"\": %q", empty))
				}

				request, _ := frontend.Recv(0)
				frameID, _ := frontend.Recv(0)
				if i, err := strconv.Atoi(frameID); err == nil {
					log.Println("======>", i)
				} else {
					log.Println("There is an Error to convert type: ", err)
				}

				fmt.Println("broker:", clientId, request)
				//  write message to file
				go writeFile(request)

				backend.Send(workerQueue[0], zmq.SNDMORE)
				backend.Send("", zmq.SNDMORE)
				backend.Send(clientId, zmq.SNDMORE)
				backend.Send("", zmq.SNDMORE)
				backend.Send(request, zmq.SNDMORE)
				backend.Send(frameID, 0)

				//  Dequeue and drop the next worker identity
				workerQueue = workerQueue[1:]

			}
		}
	}

	time.Sleep(100 * time.Millisecond)
}
