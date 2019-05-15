package main

import (
    "fmt"
    //"net"
    "os"
//    "log"
//    "io"
//    "bufio"
    "os/signal"
)


type MiddleInterface interface {
	Ex(id int, message []byte)
}

type TCPMiddle struct {

	GUI_Clients []*Client
	LUA_Clients []*Client

}

func (self*TCPMiddle) ClearZombie() {
	for i:=0;i<len(self.LUA_Clients);i++ {
		if self.LUA_Clients[i].HasError != nil {
			self.LUA_Clients[i] = nil
		}
	}
	for i:=0;i<len(self.GUI_Clients);i++ {
		if self.GUI_Clients[i].HasError != nil {
			self.GUI_Clients[i] = nil
		}
	}
}

func (self *TCPMiddle) Ex(id int, message []byte) {

    if id == GUI {
      for i:=0;i<len(self.LUA_Clients);i++ {
        if self.LUA_Clients[i] != nil {
          err := self.LUA_Clients[i].SendBytes(message)
          if err != nil {
            self.LUA_Clients[i].HasError = err
          }
        }
      }
    }else if id == LUA {
      for i:=0;i<len(self.GUI_Clients);i++ {
        if self.GUI_Clients[i] != nil {
          err := self.GUI_Clients[i].SendBytes(message)
          if err != nil {
            self.GUI_Clients[i].HasError = err
          }
        }
      }
    }

    self.ClearZombie()
}



type UDPMiddle struct {

	GUI_Client *Client
	LUA_Client *Client

}

func (self*UDPMiddle) ClearZombie() {
  
}

func (self *UDPMiddle) Ex(id int, message []byte) {

	if id == GUI {
      if self.LUA_Client != nil {
        self.LUA_Client.SendUDPBytes(message)
      }

	}else if id == LUA {
      if self.GUI_Client != nil {
        self.GUI_Client.SendUDPBytes(message)
      }
	}

	self.ClearZombie()
}


func main(){


	tcp_middle := &TCPMiddle{}

	udp_middle := &UDPMiddle{}

	go tcp_middle.serverTCP_GUI()
	go tcp_middle.serverTCP_LUA()

	go udp_middle.serverUDP_GUI()
	go udp_middle.serverUDP_LUA()
  

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	

	select {
	case <-signalChan:
		fmt.Println("ctrl +c ,exiting..")
		os.Exit(-1)
	}

}
