package monitor_listener

import (
	"fmt"
	"net"
	"simonf.net/monitor_db"
)

func checkError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
		panic(err)
	}
}

func ListenForClients() {
	ListenForUDPClients(db)
	// s := monitor_db.Service{Name: "test", Status: "ok", Updated: time.Now()}
	// c := monitor_db.NewComputer("pizero", "ok", time.Now())
	// c.SetService(&s)
	// fmt.Println(c.JSON())
	// db.AddComputer(c)
	// db.PrintComputers()
}

func ListenForUDPClients(db *monitor_db.Database) {
	serverAddr, err := net.ResolveUDPAddr("udp", ServerPort)
	checkError(err)

	/* Now listen at selected port */
	serverConn, err := net.ListenUDP("udp", serverAddr)
	checkError(err)

	defer serverConn.Close()

	buf := make([]byte, 20480)

	for {
		n, addr, err := serverConn.ReadFromUDP(buf)
		fmt.Println("Received ", string(buf[0:n]), " from ", addr)

		if err != nil {
			fmt.Println("Error: ", err)
		} else {
			c := monitor_db.NewComputerFromJSON(buf)
			db_mutex.Lock()
			db.AddComputer(c)
			db_mutex.Unlock()
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
