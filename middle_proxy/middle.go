package main

import (
    "fmt"
    "net"
    "os"
    "log"
		"io"
		"bufio"
		"os/signal"
)

const (
    CONN_HOST_A = "0.0.0.0"
    CONN_PORT_A = "8080"
    CONN_TYPE_TCP = "tcp"
		CONN_TYPE_UDP ="udp"

		CONN_HOST_B="0.0.0.0"
		CONN_PORT_B = "8081"
		
		MAX_CLIENTS = 2
)

var ResetPeer = 0

type Client struct {
	conn   net.Conn
  
  udp_conn *net.UDPConn
  udp_client_addr *net.UDPAddr
  
  Bind   int

}

func (c *Client) listen_tcp( id int) {
	if id < MAX_CLIENTS {
    TheTCPClients[id] = c
	}
  
	message := make([]byte, 1024)
	reader := bufio.NewReader(c.conn)
	
	for {
		read_number, err := reader.Read(message)
		if err != nil {
			if err != io.EOF {
				fmt.Println("Un EOF error: ",err)
			}
			c.conn.Close()
      for i:=0;i<MAX_CLIENTS;i++ {
        if TheTCPClients[i] != nil {
          TheTCPClients[i].conn.Close()
        }
      }
			break
		}else {
      ExchangeTCP(id,message[:read_number])
    
      for i:=0;i<1024;i++ {
        message[i]=0
      }
    }
	}
}

func (c *Client) listen_udp( id int) {
	if id < MAX_CLIENTS {
    TheUDPClients[id] = c
	}
  message := make([]byte, 1024)
  
  for {
      n, addr, err := c.udp_conn.ReadFromUDP(message)

      fmt.Println("UDP client : ", addr,",Read ", n)

      if err != nil {
             log.Fatal(err)
      }
      c.udp_client_addr = addr
      ExchangeUDP(id,message[:n])
  }
}

// Send text message to client
func (c *Client) Send(message string) error {
	_, err := c.conn.Write([]byte(message))
	return err
}

// Send bytes to client
func (c *Client) SendBytes(b []byte) error {
	_, err := c.conn.Write(b)
	return err
}

func (c *Client) SendUDPBytes(b []byte) error {
  _, err := c.udp_conn.WriteToUDP(b, c.udp_client_addr)
  if err != nil {
    fmt.Println(err)
  }
	return err
}


func (c *Client) Conn() net.Conn {
	return c.conn
}

func (c *Client) Close() error {
	return c.conn.Close()
}

var TheTCPClients [MAX_CLIENTS]*Client
var TheUDPClients [MAX_CLIENTS]*Client

func serverTCP_A() {
    l, err := net.Listen(CONN_TYPE_TCP, CONN_HOST_A+":"+CONN_PORT_A)
    if err != nil {
        fmt.Println("Error listening:", err.Error())
        os.Exit(1)
    }
    defer l.Close()
    fmt.Println("Listening tcp on A " + CONN_HOST_A + ":" + CONN_PORT_A)
    for {
				conn, _ := l.Accept()
				client := &Client{
					conn:   conn,
				}
				
				client.listen_tcp(0)
				
    }
}

func serverTCP_B() {
    l, err := net.Listen(CONN_TYPE_TCP, CONN_HOST_B+":"+CONN_PORT_B)
    if err != nil {
        fmt.Println("Error listening:", err.Error())
        os.Exit(1)
    }
    defer l.Close()
    fmt.Println("Listening tcp on B " + CONN_HOST_B + ":" + CONN_PORT_B)
    for {
				conn, _ := l.Accept()
				client := &Client{
					conn:   conn,
				}

				client.listen_tcp(1)
    }
}

func ExchangeTCP(id int,message []byte) {
	if id == 1 {
		if TheTCPClients[0] != nil {
			TheTCPClients[0].SendBytes(message)
		}
	}else if id == 0{
		
		if TheTCPClients[1] != nil {
			TheTCPClients[1].SendBytes(message)
		}
	}
}


func serverUDP_A() {

    ServerAddr,err := net.ResolveUDPAddr("udp",":"+CONN_PORT_A)
    if err != nil {
      panic(err)
    }
    
    l, err := net.ListenUDP(CONN_TYPE_UDP,ServerAddr)
    if err != nil {
        fmt.Println("Error listening:", err.Error())
        os.Exit(1)
    }
    defer l.Close()
    fmt.Println("Listening udp on A " + CONN_HOST_A + ":" + CONN_PORT_A)
    for {
				
				client := &Client{
					udp_conn: l,
				}
				
				client.listen_udp(0)
				
    }
}

func serverUDP_B() {
    ServerAddr,err := net.ResolveUDPAddr("udp",CONN_HOST_B+":"+CONN_PORT_B)
    if err != nil {
      panic(err)
    }

    l, err := net.ListenUDP(CONN_TYPE_UDP, ServerAddr)
    if err != nil {
        fmt.Println("Error listening:", err.Error())
        os.Exit(1)
    }
    defer l.Close()
    fmt.Println("Listening udp on B " + CONN_HOST_B + ":" + CONN_PORT_B)
    for {
				
				client := &Client{
					udp_conn: l,
				}

				client.listen_udp(1)
    }
}

func ExchangeUDP(id int,message []byte) {
	if id == 1 {
		if TheUDPClients[0] != nil {
      
			TheUDPClients[0].SendUDPBytes(message)
		}
	}else if id == 0{
		
		if TheUDPClients[1] != nil {
			TheUDPClients[1].SendUDPBytes(message)
		}
	}

}


func main(){

	go serverTCP_A()
	go serverTCP_B()

	go serverUDP_A()
	go serverUDP_B()
  
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	
  select {
		case <-signalChan:
			fmt.Println("ctrl +c ,exiting..")
			os.Exit(-1)
  }	

}
