package main

import (
	"fmt"
	"net"
	"strings"
	"time"
)

type User struct {
	name string
	conn net.Conn
}

func (u User) Read() string {
	buff := make([]byte, 1024)
	size, err := u.conn.Read(buff)
	if err != nil {
		return "Error while reading text"
	}
	return strings.TrimSpace(string(buff[:size]))
}

func (u User) write(message string) {
	_, err := u.conn.Write([]byte(message + "\n"))
	if err != nil {
		fmt.Println("Error:", err)
	}
}

func (u User) ReadMessageStream() {
	for {
		msg := u.Read()
		if msg == "" {
			continue
		}
		if strings.HasPrefix(msg, "/send") {
			cmd := strings.SplitN(msg, " ", 3)
			if len(cmd) != 3 {
				u.write("Error: Couldn't understand the command. \nCorrect Usage: /send <recipient> <message>")
				continue
			}
			receiverName, message := cmd[1], cmd[2]
			if receiverName == "" || message == "" {
				u.write("Error: malformed command\n")
				continue
			}
			if receiver, found := Users[receiverName]; found {
				message = fmt.Sprintf("%d:%d:%d [%s] >> %s",
					time.Now().Hour(),
					time.Now().Minute(),
					time.Now().Second(), u.name, message)
				receiver.write(message + "\n")
			} else {
				u.write("Error: Couldn't find " + receiverName + "\n")
				continue
			}
		} else if strings.HasPrefix(msg, "/list") {
			for k := range Users {
				u.write(k)
			}
		} else if strings.HasPrefix(msg, "/help") {
			doc := "Commands are following: \n/list : list all the User available to chat\n/send <person> <message> : Send private message to User listed using /command\n/h : Show help text\n---\n"
			u.write(doc)
		} else {
			for k, v := range Users {
				if k == u.name {

				} else {
					v.write(msg + "\n")
				}
			}
		}
	}
}

const (
	NETWORK string = "tcp"
	PORT    string = "8080"
)

var Users = make(map[string]User, 8)

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
	Host := fmt.Sprintf(":%s", PORT)
	tcpAddr, err := net.ResolveTCPAddr(NETWORK, Host)
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
		Users[name] = User{name: name, conn: conn}
		conn.Write([]byte("Hello " + name + "\n"))
		doc := "Commands are following: \n/list : list all the User available to chat\n/send <person> <message> : Send private message to User listed using /command\n/help : Show help text\n---\n"
		conn.Write([]byte(doc))
		go Users[name].ReadMessageStream()
	}()
}
