//
//  Clients
//

package main

import (
	"fmt"
	zmq "github.com/pebbe/zmq4"
	"log"
	"math/rand"
	"strconv"
	"sync"

	//"math/rand"

	"time"
)

const (
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	//ir := IntRange{50,8000}
	MinMessageSize = 50
	MaxMessageSize = 100
	SendRate       = 10000
)

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

func clientTask(wg *sync.WaitGroup) {
	time.Sleep(1 * time.Second)
	client, _ := zmq.NewSocket(zmq.REQ)
	defer client.Close()
	client.Connect("ipc://frontend.ipc")
	//  Send request, get reply
	//  set RATE ZMQ_RATE
	client.SetRate(SendRate)
	client.SetMaxmsgsize(MaxMessageSize)

	start := time.Now()
	requests := SendRate
	r := rand.New(rand.NewSource(55))
	ir := IntRange{MinMessageSize, MaxMessageSize}
	realInt := 0
	for i := 0; i < requests; i++ {
		str := randStringBytes(ir.NextRandom(r))
		client.Send(str, zmq.SNDMORE)
		client.Send(strconv.Itoa(i), 0)
		client.Recv(0)
		fmt.Println(realInt)
	}

	fmt.Printf("%d calls/second\n", int64(float64(requests)/time.Since(start).Seconds()))

	value, err := GetValue("real")
	if err != nil {
		log.Println(err)
	}
	if value == "" {
		realInt = 0
	}

	if value != "" {
		realInt, err = strconv.Atoi(value)
		if err != nil {
			log.Println("There is an Error: ", err)

		} else {
			realInt = realInt + requests
		}

	} else {
		realInt = realInt + requests
	}

	err = SetValue("real", realInt)
	fmt.Println(err)
	wg.Done()

}

func main() {
	var wg sync.WaitGroup
	for true {
		wg.Add(1)
		go clientTask(&wg)
		time.Sleep(1000 * time.Millisecond)
	}
	wg.Wait()
}
