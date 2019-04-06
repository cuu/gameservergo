package main

import (
  //"fmt"
  
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
  
  reader := bufio.NewReader(conn)
  var ret string

  for {
    
    message, _ := reader.ReadString('\n')
    //fmt.Println( len(message))
    ret = gs.GameThread.ProcessCmd(message)
    conn.Write([]byte(ret+"\n"))
    
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
