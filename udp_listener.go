package monitor_listener

import (
	"fmt"
	"net"
	"simonf.net/monitor_db"
)

const ServerPort = ":41237"

func checkError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
		panic(err)
	}
}

func ListenForClients() {
	ListenForUDPClients(db)
}

func ListenForUDPClients(db *monitor_db.Database) {
	addr, err := net.ResolveUDPAddr("udp4", ":41237")
	checkError(err)

	socket, err := net.ListenUDP("udp4", addr)

	fmt.Println("Listening for clients on 41237")

	defer socket.Close()

	buf := make([]byte, 20480)

	for {
		n, addr, err := socket.ReadFromUDP(buf)
		fmt.Println("Received ", string(buf[0:n]), " from ", addr)

		if err != nil {
			fmt.Println("Error: ", err)
		} else {
			c := monitor_db.NewComputerFromJSON(buf[0:n])
			fmt.Printf("Unmarshalled computer %s. Adding it", c.Name)
			db.AddComputer(c)
		}
	}
}

// func stripPort(udp_host_and_port string) string {
//   fmt.Println("Server: ", udp_host_and_port)
//   i := strings.Index(udp_host_and_port, ":")
//   pruned := udp_host_and_port[0:i]
//   fmt.Printf("Colon at %v gives %v\n", i, pruned)
//   return string(pruned)
// }
