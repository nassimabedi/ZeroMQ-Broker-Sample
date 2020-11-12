//
//  Clients
//

package main

import (
	zmq "github.com/pebbe/zmq4"
	"math/rand"

	"fmt"
	//"math/rand"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// range specification, note that min <= max
type IntRange struct {
	min, max int
}

// get next random value within the interval including min and max
func (ir *IntRange) NextRandom(r *rand.Rand) int {
	return r.Intn(ir.max-ir.min+1) + ir.min
}

func client_task() {

	//a:= zmq.SetMaxMsgsz(100)
	client, _ := zmq.NewSocket(zmq.REQ)
	defer client.Close()
	// set_id(client) //  Set a printable identity
	client.Connect("ipc://frontend.ipc")

	//  Send request, get reply
	//Nassim set RATE ZMQ_RATE
	client.SetRate(1000)
	//client.SetMaxmsgsize(100)
	client.SetSndbuf(80)

	start := time.Now()
	requests := 10000
	r := rand.New(rand.NewSource(55))
	//ir := IntRange{50,8000}
	ir := IntRange{50, 100}
	for i := 0; i < requests; i++ {
		str := randStringBytes(ir.NextRandom(r))
		//client.Send("HELLO1111111111",0)
		client.Send(str, 0)
		//client.SendMessage(i,str)

		//client.Send([]byte("hello"), 0)

		client.Recv(0)
	}
	fmt.Printf("%d calls/second\n", int64(float64(requests)/time.Since(start).Seconds()))
}

//  This is the main task. It starts the clients and workers, and then
//  routes requests between the two layers. Workers signal READY when
//  they start; after that we treat them as ready when they reply with
//  a response back to a client. The load-balancing data structure is
//  just a queue of next available workers.

func main() {

	go client_task()
	time.Sleep(100 * time.Millisecond)

}
