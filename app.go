package main

import (
	"fmt"
	"net"
	"strings"
)

type People struct {
	name string
	conn net.Conn
}

func (p People) Read() string {
	buff := make([]byte, 1024)
	size, err := p.conn.Read(buff)
	if err != nil {
		return "Error while reading text"
	}
	return strings.TrimSpace(string(buff[:size]))
}

func (p People) write(message string) {
	_, err := p.conn.Write([]byte(message + "\n"))
	if err != nil {
		fmt.Println("Error:", err)
	}
}

func (p People) ReadMessageStream() {
	for {
		msg := p.Read()
		if msg == "" {
			continue
		}
		cmd := make([]string, 16)
		if strings.HasPrefix(msg, "/send") {
			cmd = strings.SplitN(msg, " ", 3)
			if len(cmd) != 3 {
				p.write("Error: Couldn't understand the command. \nCorrect Usage: /send <recipient> <message>")
				continue
			}
			receiverName, message := cmd[1], cmd[2]
			if receiverName == "" || message == "" {
				p.write("Error: malformed command\n")
				continue
			}
			if receiver, found := Users[receiverName]; found {
				receiver.write(message + "\n")
			} else {
				p.write("Error: Couldn't find " + receiverName + "\n")
				continue
			}
		} else if strings.HasPrefix(msg, "/list") {
			for k := range Users {
				p.write(k)
			}
		} else {
			for k, v := range Users {
				if k == p.name {

				} else {
					v.write(msg + "\n")
				}
			}
		}
	}
}

const (
	NETWORK string = "tcp"
	IP      string = "0.0.0.0"
	PORT    string = ":8080"
)

var Users = make(map[string]People, 8)

func main() {
	connaction := Connect()
	defer connaction.Close()
	for {
		conn, er := connaction.AcceptTCP()
		if er != nil {
			fmt.Println("couldn't connect via tcp")
			continue
		}
		HandleConnection(conn)
	}
}

func Connect() *net.TCPListener {
	tcpAddr, err := net.ResolveTCPAddr(NETWORK, IP+PORT)
	if err != nil {
		panic(err)
	}
	fmt.Println("localhost listing at prot 8080\nType Ctrl+C to exit")
	listner, err := net.ListenTCP(NETWORK, tcpAddr)
	if err != nil {
		panic(err)
	}
	return listner
}
func HandleConnection(conn net.Conn) {
	go func() {
		conn.Write([]byte("Name: "))
		input := make([]byte, 1024)
		size, err := conn.Read(input)
		if err != nil {
			fmt.Println("Error:", err)
			conn.Close()
			return
		}
		name := strings.TrimSpace(string(input[:size]))
		Users[name] = People{name: name, conn: conn}
		conn.Write([]byte("Hello " + name + "\n"))
		doc := "Commands are following: \n/list : list all the people available to chat\n/send <person> <message> : Send private message to people listed using /command\n/h : Show help text\n---\n"
		conn.Write([]byte(doc))
		go Users[name].ReadMessageStream()
	}()
}
