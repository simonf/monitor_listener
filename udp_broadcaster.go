package monitor_listener

import (
	"fmt"
	"net"
	"os"
	"simonf.net/monitor_db"
	"time"
)

const BroadcastPort = 41238

var db = monitor_db.NewDatabase()

// var db_mutex = &sync.Mutex{}

func PeriodicallyAdvertise() {
	for {
		BroadcastHostName()
		time.Sleep(2 * time.Minute)
	}
}

func BroadcastHostName() {
	BroadcastAddress := net.IPv4(255, 255, 255, 255)

	socket, err := net.DialUDP("udp4", nil, &net.UDPAddr{
		IP:   BroadcastAddress,
		Port: BroadcastPort,
	})

	if err != nil {
		fmt.Println("Error connecting to the network broadcast address")
		return
	}

	defer socket.Close()

	name, err := os.Hostname()

	if err != nil {
		panic(err)
	}

	fmt.Printf("Broadcasting %v\n", name)

	buf := []byte(name)
	_, err = socket.Write(buf)
	if err != nil {
		fmt.Println(name, err)
	}

}
