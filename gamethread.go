package main

import(
    "fmt"
    "encoding/json"
    "strconv"
    gotime "time"
//	  "github.com/veandco/go-sdl2/sdl"


  	//"github.com/cuu/gogame"
    

    
  	"github.com/cuu/gogame/color"
  	"github.com/cuu/gogame/event"

  	"github.com/cuu/gogame/time"
    "github.com/cuu/gogame/display"
    "github.com/cuu/gogame/font"
    "github.com/cuu/gogame/draw"
    "github.com/cuu/gogame/rect"


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
}

func NewGoGameThread() *GoGameThread {
  p := &GoGameThread{}
  
  p.Width =  640
  p.Height = 480
  
  p.Inited = false
  p.DT = time.NewClock().Tick(30)
  
  fmt.Println("DT=",p.DT)
  
  p.State = "draw"
  p.BgColor = color.NewColor(0,0,0,255)

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
      if ev.Data["Key"] == "Q" {
        break
			}
      if ev.Data["Key"] == "P" {
        screen := display.GetSurface() 
	      draw.Line(screen,color.NewColor(255,44,255,255), 0,100, 320,100,3)
	      draw.Line(screen,color.NewColor(255,44,255,255), 10, 0, 10,250,4)

	      rect2 := rect.Rect(3,120,200,30)
	      draw.AARoundRect(screen,&rect2,&color.Color{0,213,222,255},10,0, &color.Color{0,213,222,255})
        display.Flip()
        
      }
		}
    
    time.NewClock().Tick(30)
    
	}
}

func (self *GoGameThread) Flip() {
  if self.Frames == 0 {
    self.PrevTime = gotime.Now()
  }  
  
  self.ThePico8.Flip()
  
  display.Flip()
  self.Frames+=1
  
  self.CurrentTime = gotime.Now()
  
  if self.CurrentTime.Sub(self.PrevTime) > 10*gotime.Second {
    fps := self.Frames /10
    print("fps is: ",fps)
    self.Frames = 0
    self.PrevTime = self.CurrentTime
  }
}

func (self *GoGameThread) Run() int {

  self.InitWindow()
  self.EventLoop()  
  
  return 0
  
}

type CmdArg struct {
  Type string `json:"Type"`
  Value string `json:"Value"`
}

func (self *CmdArg) GetInt() int{
  if self.Type == "I" {
    val,err := strconv.Atoi(self.Value)
    if err == nil {
      return val
    }else {
      fmt.Println(err)
    }
  }else {
    fmt.Println("try to get int from not-integer Arg")
  }
  
  return -1
}

func (self *CmdArg) GetStr() string {  
  return self.Value

}

type ACmd struct {
  Func string  `json:"Func"`
  Args []CmdArg `json:"Args"`
}

func (self *GoGameThread) ProcessCmd(cmd string) {
  
  acmd := &ACmd{}
   
  if err := json.Unmarshal([]byte(cmd), &acmd); err != nil {
    fmt.Println(err)
    return
  }
  
  if acmd.Func == "print" {
    self.ThePico8.Print(acmd.Args)
  }
  
  if acmd.Func == "flip" {
    self.Flip()
  }

}




