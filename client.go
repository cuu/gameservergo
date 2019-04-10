package main

import (
  //"fmt"
   "io"
  "net"
  "bufio"
  "github.com/veandco/go-sdl2/sdl"

  
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
      ret = gs.GameThread.ProcessCmd(message)
      conn.Write([]byte(ret+"\n"))
    }

    if err != nil {
      if err != io.EOF {
	panic(err)
      }
    }
  }
}

func main() {
  
  gs := NewGameClient()
  
  go start_tcp_client(gs)
  
  sdl.Main(func() {
    gs.GameThread.Run()
  })
  
  return
}
