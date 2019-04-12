package main

import(
    "fmt"
    "sync"
    "encoding/json"
    "strconv"
    "net"
    gotime "time"
//   "github.com/veandco/go-sdl2/sdl"


     //"github.com/cuu/gogame2"
    
    "github.com/cuu/gogame2/color"
    "github.com/cuu/gogame2/event"

    "github.com/cuu/gogame2/time"
    "github.com/cuu/gogame2/display"
    "github.com/cuu/gogame2/font"
    //"github.com/cuu/gogame2/draw"
    //"github.com/cuu/gogame2/rect"


)

const (
    RUNEVT=1

)

type GoGameThread struct {
  
  Width int
  Height int
  Inited bool
  
  DT int
  
  BgColor *color.Color
  State string //draw,res
  Resource  string //
  ConsoleType string 
  
  ThePico8 *Pico8
  
  Frames int 
  
  PrevTime gotime.Time
  CurrentTime gotime.Time

  KeyLog sync.Map
  
  UdpConn net.Conn
  TcpConn net.Conn
}

func NewGoGameThread() *GoGameThread {
  p := &GoGameThread{}
  
  p.Width =  320
  p.Height = 240
  
  p.Inited = false
  p.DT = time.NewClock().Tick(30)
  
  fmt.Println("DT=",p.DT)
  
  p.State = "draw"
  p.BgColor = color.NewColor(0,0,0,255)
  
  p.KeyLog = sync.Map{}
  p.KeyLog.Store("Left",-1)
  p.KeyLog.Store("Right",-1)
  p.KeyLog.Store("Up",-1)
  p.KeyLog.Store("Down",-1)
  p.KeyLog.Store("U",-1)
  p.KeyLog.Store("I",-1)
  p.KeyLog.Store("Return",-1)
  p.KeyLog.Store("Escape",-1)
  
  return p 
}

func (self*GoGameThread) InitWindow(){
  if self.Inited == false {
      self.Inited = true
      display.Init()
      font.Init()
      screen := display.SetMode(int32(self.Width),int32(self.Height),0,32)
  
      self.ThePico8 = NewPico8()
      self.ThePico8.HWND = screen
      
      event.AllocEvents(8)
      event.AddCustomEvent(RUNEVT)

  }
}

func (self*GoGameThread) QuitWindow(){
  self.Inited = false
  display.Destroy()  
}

func (self*GoGameThread) EventLoop() {
  
  for self.Inited {
    ev := event.Poll()
    if ev.Type == event.QUIT {
      break
    }
    
    if ev.Type == event.USEREVENT {
      fmt.Println(ev.Data["Msg"])
    }
    if ev.Type == event.KEYDOWN {
      if ev.Data["Key"] == "Escape" {
        break
      }
      fmt.Fprintf(self.UdpConn,fmt.Sprintf("%s,%s\n",ev.Data["Key"],"Down"))

    }
    if ev.Type == event.KEYUP {
      self.KeyLog.Store(ev.Data["Key"],-1)
      fmt.Fprintf(self.UdpConn,fmt.Sprintf("%s,%s\n",ev.Data["Key"],"Up"))
    }
    
    time.SDL_Delay(30)
    
  }
}

func (self *GoGameThread) FlipLoop() {
  for {
    if self.Frames == 0 {
      self.PrevTime = gotime.Now()
    }  
  
    //self.ThePico8.Flip()
    display.UpdatePixels()
    display.Flip()
    self.Frames+=1
  
    self.CurrentTime = gotime.Now()
  
    if self.CurrentTime.Sub(self.PrevTime) > 10*gotime.Second {
      fps := self.Frames /10
      println("fps is: ",fps)
      self.Frames = 0
      self.PrevTime = self.CurrentTime
    }
    
    time.NewClock().Tick(30)

  }
}

func (self *GoGameThread) Btn(args []CmdArg) string {

  if len(args) < 2 {
    return "FALSE"
  }
  
  keycode_string := args[0].GetStr()
  //player_idx     := args[1].GetInt() // Not implemented yet
  if val,ok := self.KeyLog.Load(keycode_string); ok {
    if val.(int) >= 0 {
      return "TRUE"
    }
  }
  
  return "FALSE"

}

func (self *GoGameThread) StartTcp() {

  conn, err := net.Dial("tcp", "127.0.0.1:8081")
  if err != nil {
    panic(fmt.Sprintln("tcp Dial error %v", err))
  }  
  self.TcpConn = conn

}

func (self *GoGameThread) Run() int {

  self.InitWindow()

  go self.FlipLoop()
  
  //go self.ThePico8.FlipLoop()
  //self.StartUdp()

  self.EventLoop()  


  return 0
  
}

type CmdArg struct {
  Type string `json:"Type"`
  Value interface{} `json:"Value"`
}

func (self *CmdArg) GetInt() int {
  
  switch v := self.Value.(type) {
  case int:
    return self.Value.(int)
  case float64:
    tmp := self.Value.(float64)
    return int(tmp)
  case string:
    fmt.Printf("String: %v", v)
    tmp,err := strconv.Atoi(self.Value.(string))
    if err != nil {
      fmt.Println(err)
    }else {
      return int(tmp)
    }
  case bool:
    if self.Value.(bool) == true {
      return 1
    }else {
      return 0
    }
  default:
    panic("Value type error")
  }
  return -1
}

func (self *CmdArg) GetStr() string {
  return self.Value.(string)
}

func (self *CmdArg) GetBool() bool {
  return self.Value.(bool)
}

type ACmd struct {
  Role string `json:"Role"` 
  Func string  `json:"Func"`
  Args []CmdArg `json:"Args"`
}

func (self *GoGameThread) ProcessCmd(cmd string) string {
  
  if len(cmd) == 0 {
    return "Error"
  }
  acmd := &ACmd{}

  if err := json.Unmarshal([]byte(cmd), &acmd); err != nil {
    fmt.Println(fmt.Sprintf("%v,%s,%d",err,cmd,len(cmd)))
    return "Error"
  }
  
  if acmd.Func == "res" {
    self.ThePico8.Res(acmd.Args)
  }

  if acmd.Func == "flip" {
    self.ThePico8.Flip()
  }

  if acmd.Func == "print" {
    self.ThePico8.Print(acmd.Args)
  }
  
  if acmd.Func == "pico8" {
    self.ThePico8.SetVersion(acmd.Args)
  }
  if acmd.Func == "map" {
    self.ThePico8.Map(acmd.Args)
  }

  if acmd.Func == "spr" {
    self.ThePico8.Spr(acmd.Args)
  }
  if acmd.Func == "mget" {
    return self.ThePico8.MGet(acmd.Args)
  }

  if acmd.Func == "rect" {
    self.ThePico8.Rectfill(acmd.Args)
  }

  if acmd.Func == "rectfill" {
    self.ThePico8.Rectfill(acmd.Args)
  }

  if acmd.Func == "btn" {
    return self.Btn(acmd.Args)
  }

  if acmd.Func == "pal" {
    self.ThePico8.Pal(acmd.Args)
  }
  
  if acmd.Func == "palt" {
    self.ThePico8.Palt(acmd.Args...)
  }
  
  if acmd.Func == "circ" {
    self.ThePico8.Circ(acmd.Args...)
  }
  if acmd.Func == "circfill" {
    self.ThePico8.Circfill(acmd.Args...)
  }
  
  return "O"
}




