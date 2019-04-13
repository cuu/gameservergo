package main

import (
  "fmt"
  //"log"
  "bytes"
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

func dropCR(data []byte) []byte {
    if len(data) > 0 && data[len(data)-1] == '\r' {
        return data[0 : len(data)-1]
    }
    return data
}


func ScanCRLF(data []byte, atEOF bool) (advance int, token []byte, err error) {
        if atEOF && len(data) == 0 {
            return 0, nil, nil
        }
        if i := bytes.Index(data, []byte{'\r','\n'}); i >= 0 {
            // We have a full newline-terminated line.
            return i + 2, dropCR(data[0:i]), nil
        }
        // If we're at EOF, we have a final, non-terminated line. Return it.
        if atEOF {
            return len(data), dropCR(data), nil
        }
        // Request more data.
        return 0, nil, nil
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
	gs.GameThread.ProcessCmds(message)

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
