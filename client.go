package main

import (
  "fmt"
  //"log"
  //"bytes"
   "io"
  "net"
  "bufio"
  //"strings"
  "github.com/veandco/go-sdl2/sdl"
  
  gotime "time"
  
)
//tcp server
type GameClient struct {
  
  GameThread *GoGameThread
  
}

func NewGameClient() *GameClient {
  p := &GameClient{}
  
  p.GameThread = NewGoGameThread()
  
  return p
}


func start_tcp_client(gs *GameClient) {
  conn, err := net.Dial("tcp", "127.0.0.1:8081")
  if err != nil {
    panic(err)
  }
  
  var ret string
  reader := bufio.NewReader(conn)
  for {
    message, err := reader.ReadString('\n')
    //fmt.Println( len(message))
    if len(message) > 0 {
      ret = gs.GameThread.ProcessCmd([]byte(message))
      conn.Write([]byte(ret+"\n"))
    }

    if err != nil {
      if err != io.EOF {
	panic(err)
      }
    }
  }
}

func start_udp_client(gs *GameClient) {
    
    gotime.Sleep(1000 * gotime.Millisecond)

    conn, err := net.Dial("udp", "127.0.0.1:8081")
    if err != nil {
      panic(fmt.Sprintln("tcp Dial error %v", err))
    }  
    
    defer conn.Close()
    
    conn.Write([]byte("ping"))
    gs.GameThread.UdpConn = conn
    reader := bufio.NewReader(conn)

    for {
      message,_ := reader.ReadString('\n')
      
      if len(message) > 0 {
	gs.GameThread.ProcessLispCmds(message)

      }     
    }

}

func main() {
  
  gs := NewGameClient()
  
  go start_tcp_client(gs)

  go start_udp_client(gs)

  sdl.Main(func() {
    gs.GameThread.Run()
  })
  
  return
}
