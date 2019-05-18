package main

import (
    "fmt"
    "net"
//    "os"
    "log"
    "io"
    "bufio"
//    "os/signal"
    gotime "time"
)

type UDP_Addr struct {
	client_info  string
	client_addr *net.UDPAddr
	LastActive gotime.Time // this is for udp dead connection

}

type Client struct {

	tcp_conn   net.Conn
  	udp_conn  *net.UDPConn
	udp_client_addr []*UDP_Addr
	Bind   int
  	LastActive gotime.Time //IsZero
  	HasError error  // this is for tcp socket error mark
  	Parent MiddleInterface
}

func (c *Client) ReadTCP( id int) {  
  
	message := make([]byte, 1024)
	reader := bufio.NewReader(c.tcp_conn)
	
	for {
		read_number, err := reader.Read(message)
		if err != nil {
			if err != io.EOF {
				fmt.Println("Un EOF error: ",err)
			}
			c.tcp_conn.Close()
      		c.HasError = err
			break
		}else {

	      c.Parent.Ex(id,message[:read_number])
    
	      for i:=0;i<1024;i++ {
    	    message[i]=0
	      }
    	  c.LastActive = gotime.Now()
	    }
	}
  
}

func (c *Client) ReadUDP( id int) {
  
  message := make([]byte, 4096)
  
  for {
  	n, addr, err := c.udp_conn.ReadFromUDP(message)
    //fmt.Println("UDP client" ,id , " : ", addr,",Read ", n)

    if err != nil {
    	log.Fatal(err)
    }
  	hasit := false
    for i:=0;i<len(c.udp_client_addr);i++ {
    	if c.udp_client_addr[i] != nil {
			if c.udp_client_addr[i].client_info == addr.String() {
				hasit = true
				c.udp_client_addr[i].client_addr = addr
				c.udp_client_addr[i].LastActive = gotime.Now()
				break
			}
		}
    }

    if hasit == false {
    	a_new_udp_client := &UDP_Addr{}
    	a_new_udp_client.client_info = addr.String()
    	a_new_udp_client.client_addr= addr
    	a_new_udp_client.LastActive = gotime.Now()
    	c.udp_client_addr = append(c.udp_client_addr,a_new_udp_client)
	}

    c.Parent.Ex(id,message[:n])
  }
}

// Send text message to client
func (c *Client) Send(message string) error {
	_, err := c.tcp_conn.Write([]byte(message))
	return err
}

// Send bytes to client
func (c *Client) SendBytes(b []byte) error {
	_, err := c.tcp_conn.Write(b)
	return err
}

func (c *Client) SendUDPBytes(b []byte) error {

	for i:=0;i< len(c.udp_client_addr); i++ {
		if c.udp_client_addr[i] != nil {
			_, err := c.udp_conn.WriteToUDP(b, c.udp_client_addr[i].client_addr)
			if err != nil {
				fmt.Println(err)
				return err
			}else {
				c.udp_client_addr[i].LastActive = gotime.Now()
			}
		}
	}
	return nil
}


func (c *Client) Conn() net.Conn {
	return c.tcp_conn
}

func (c *Client) Close() error {
	return c.tcp_conn.Close()
}




