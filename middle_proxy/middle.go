package main

import (
    "fmt"
    "net"
    "os"
		"io"
		"bufio"
		"os/signal"
)

const (
    CONN_HOST_A = "localhost"
    CONN_PORT_A = "8080"
    CONN_TYPE = "tcp"
		
		CONN_HOST_B="localhost"
		CONN_PORT_B = "8081"
		
		MAX_CLIENTS = 2
)

var ResetPeer = 0

type Client struct {
	conn   net.Conn
}

func (c *Client) listen( id int) {
	if id < MAX_CLIENTS {
		TheClients[id] = c
	}
  
  ResetPeer+=1
  
	message := make([]byte, 1024)
	reader := bufio.NewReader(c.conn)
	
	for {
		read_number, err := reader.Read(message)
		if err != nil {
			if err != io.EOF {
				fmt.Println("Un EOF error: ",err)
			}
			c.conn.Close()
      if ResetPeer > 0 {
        for i:=0;i<MAX_CLIENTS;i++ {
          if TheClients[i] != nil {
            TheClients[i].conn.Close()
          }
          ResetPeer-=1
        }
      }
			return
		}
		Exchange(id,message[:read_number])
		for i:=0;i<1024;i++ {
			message[i]=0
		}
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

func (c *Client) Conn() net.Conn {
	return c.conn
}

func (c *Client) Close() error {
	return c.conn.Close()
}

var TheClients [MAX_CLIENTS]*Client

func serverA() {
    l, err := net.Listen(CONN_TYPE, CONN_HOST_A+":"+CONN_PORT_A)
    if err != nil {
        fmt.Println("Error listening:", err.Error())
        os.Exit(1)
    }
    defer l.Close()
    fmt.Println("Listening on A " + CONN_HOST_A + ":" + CONN_PORT_A)
    for {
				conn, _ := l.Accept()
				client := &Client{
					conn:   conn,
				}
				
				client.listen(0)
				
    }
}

func serverB() {
    l, err := net.Listen(CONN_TYPE, CONN_HOST_B+":"+CONN_PORT_B)
    if err != nil {
        fmt.Println("Error listening:", err.Error())
        os.Exit(1)
    }
    defer l.Close()
    fmt.Println("Listening on B " + CONN_HOST_B + ":" + CONN_PORT_B)
    for {
				conn, _ := l.Accept()
				client := &Client{
					conn:   conn,
				}

				client.listen(1)
    }
}

func Exchange(id int,message []byte) {
	if id == 1 {
		if TheClients[0] != nil {
			TheClients[0].SendBytes(message)
		}
	}else if id == 0{
		
		if TheClients[1] != nil {
			TheClients[1].SendBytes(message)
		}
	}

}


func main(){

	go serverA()
	go serverB()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	
  select {
		case <-signalChan:
			fmt.Println("ctrl +c ,exiting..")
			os.Exit(-1)
  }	

}
