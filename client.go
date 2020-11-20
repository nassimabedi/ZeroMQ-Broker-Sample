//
//  Client : generate package and send to broker
//

package main

import (
	zmq "github.com/pebbe/zmq4"
	"log"
	"math/rand"
	"sync"
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

func clientTask(wg *sync.WaitGroup, m *sync.Mutex) {
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

	for i := 0; i < requests; i++ {
		str := randStringBytes(ir.NextRandom(r))
		client.Send(str, 0)

		result, err := readAndWrite(RealFilename, m)
		if err != nil {
			log.Fatal(err)
		}

		client.Recv(0)
		log.Println("package number: ", result)

	}

	log.Printf("%d calls/second\n", int64(float64(requests)/time.Since(start).Seconds()))
	wg.Done()

}

func main() {
	var wg sync.WaitGroup
	var m sync.Mutex
	for true {
		wg.Add(1)
		go clientTask(&wg, &m)
		time.Sleep(1000 * time.Millisecond)
	}
	wg.Wait()
}
