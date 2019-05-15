package main

import (
    "fmt"
    "net"
    "os"
//    "log"
//    "io"
//	"bufio"
//	"os/signal"
)


func (self *TCPMiddle) serverTCP_GUI() {
    l, err := net.Listen(CONN_TCP, GUI_HOST+":"+GUI_PORT)
    if err != nil {
        fmt.Println("Error listening:", err.Error())
        os.Exit(1)
    }
    defer l.Close()
    fmt.Println("Listening tcp on GUI " + GUI_HOST + ":" + GUI_PORT)
    for {
		conn, _ := l.Accept()
		client := &Client{
			tcp_conn:   conn,
			HasError:nil,
			Parent:self,
		}
        
        has_seat := -1
        
        for i:=0;i<len(self.GUI_Clients);i++ {
          if self.GUI_Clients[i]== nil {
            has_seat = i
            break
          }
        }
        
        if has_seat == -1 {
          self.GUI_Clients = append(self.GUI_Clients,client)
          has_seat = len(self.GUI_Clients) - 1
        }else{
          
            if self.GUI_Clients[has_seat] == nil {
              self.GUI_Clients[has_seat] = client // take seat
            }
        }
        
        client.Bind = has_seat
        go client.ReadTCP(GUI)
				
    }
}


func (self *TCPMiddle) serverTCP_LUA() {
    l, err := net.Listen(CONN_TCP, LUA_HOST+":"+LUA_PORT)
    if err != nil {
        fmt.Println("Error listening:", err.Error())
        os.Exit(1)
    }

    defer l.Close()

    fmt.Println("Listening tcp on LUA " + LUA_HOST + ":" + LUA_PORT)
    for {
		conn, _ := l.Accept()
		client := &Client{
			tcp_conn:   conn,
			HasError:nil,
			Parent:self,
		}

        has_seat := -1
        
        for i:=0; i<len(self.LUA_Clients); i++ {
          if self.LUA_Clients[i] == nil {
            has_seat = i
            break
          }
        }
        
        if has_seat == -1 {
          
          self.LUA_Clients = append(self.LUA_Clients,client)

          has_seat = len(self.LUA_Clients) - 1

        }else{

          if self.LUA_Clients[has_seat] == nil {
              self.LUA_Clients[has_seat] = client // take seat
          }
        }
        client.Bind = has_seat
        go client.ReadTCP(LUA)
    }
}

func (self*UDPMiddle) serverUDP_GUI() {

    ServerAddr,err := net.ResolveUDPAddr("udp",":"+GUI_PORT)
    if err != nil {
      panic(err)
    }
    
    l, err := net.ListenUDP(CONN_UDP,ServerAddr)
    if err != nil {
        fmt.Println("Error listening:", err.Error())
        os.Exit(1)
    }
    defer l.Close()
    fmt.Println("Listening udp on GUI " + GUI_HOST + ":" + GUI_PORT)
    for {
				
        client := &Client{
			udp_conn: l,
			Parent:self,
        }

        self.GUI_Client = client
		client.ReadUDP(GUI) // no goroutine
				
    }
}


func (self*UDPMiddle) serverUDP_LUA() {

    ServerAddr,err := net.ResolveUDPAddr("udp",":"+LUA_PORT)
    if err != nil {
      panic(err)
    }
    
    l, err := net.ListenUDP(CONN_UDP,ServerAddr)
    if err != nil {
        fmt.Println("Error listening:", err.Error())
        os.Exit(1)
    }
    defer l.Close()
    fmt.Println("Listening udp on LUA " + LUA_HOST + ":" + LUA_PORT)
    for {
				
		client := &Client{
		    udp_conn: l,
		    Parent:self,
		}
        self.LUA_Client = client
		client.ReadUDP(LUA) // no goroutine

    }
}


