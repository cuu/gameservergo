package main
import (
	"fmt"
	"strings"
	"strconv"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"

	"github.com/cuu/gogame2/surface"
	"github.com/cuu/gogame2/rect"

	"github.com/cuu/gogame2/color"
	"github.com/cuu/gogame2/font"
	"github.com/cuu/gogame2/display"
	"github.com/cuu/gogame2/transform"
	"github.com/cuu/gogame2/draw"
	//"github.com/cuu/gogame2/time"



)

type Pico8 struct{
	Width int // 128
	Height int // 128
	Version int 
	
	MapMatrix [64*128]int
	
	CanvasHWND    *sdl.Surface
	HWND          *sdl.Surface
	DisplayCanvas *sdl.Surface
	DrawCanvas    *sdl.Surface
	GfxSurface    *sdl.Surface
	
	bg_color   *color.Color
	
	Font  *ttf.Font
	
	pal_colors [] sdl.Color
	draw_colors []sdl.Color
	display_colors []sdl.Color
	
	
	DrawPaletteIdx [16]int
	DrawPalette *sdl.Palette
	
	DisplayPalette *sdl.Palette
	
	PalTransparent [16]int
	
	ClipRect *sdl.Rect
	PenColor int // 0 - 15
	
	Cursor [2]int
	
	CameraDx int 
	CameraDy int
	PaletteModified bool
	
	
	SpriteFlags [256]int
	
	Uptime int

	Resource map[string]string
	
	ToFlip bool

	frames int
	curr_time int
	prev_time  int
	fps        int
}

func NewPico8() *Pico8 {
	p := &Pico8{}
	p.Width = 128
	p.Height = 128
	p.Version = 8 
	
	p.ClipRect = rect.NewRect(0,0,0,0)
	
	p.DrawPalette,_       = sdl.AllocPalette(16)
	p.DisplayPalette,_    = sdl.AllocPalette(16)
	
	p.pal_colors = []sdl.Color{
					sdl.Color{0,0,0,255},
					sdl.Color{29,43,83,255},
					sdl.Color{126,37,83,255},
					sdl.Color{0,135,81,255},
					sdl.Color{171,82,54,255},
					sdl.Color{95,87,79,255},
					sdl.Color{194,195,199,255},
					sdl.Color{255,241,232,255},
					sdl.Color{255,0,77,255},
					sdl.Color{255,163,0,255},
					sdl.Color{255,240,36,255},
					sdl.Color{0,231,86,255},
					sdl.Color{41,173,255,255},
					sdl.Color{131,118,156,255},
					sdl.Color{255,119,168,255},
					sdl.Color{255,204,170,255}}
	
	p.DisplayCanvas = surface.ASurface(p.Width,p.Height,0,32)
	p.DrawCanvas    = surface.ASurface(p.Width,p.Height,0,32)
	p.GfxSurface    = surface.ASurface(p.Width,p.Height,0,32)
	surface.Fill(p.GfxSurface, &color.Color{0,0,0,255})
	for i:=0;i<len(p.pal_colors);i++ {
		p.draw_colors = append(p.draw_colors,p.pal_colors[i])
		p.display_colors = append(p.display_colors,p.pal_colors[i])
	}
	p.DrawPalette.SetColors(      p.draw_colors)
	p.DisplayPalette.SetColors(   p.display_colors)	
	
	p.DisplayCanvas.SetPalette(p.DisplayPalette)
	p.DrawCanvas.SetPalette(   p.DrawPalette)
	

	
	for i:=0;i<16;i++ {
		p.DrawPaletteIdx[i]=i
		
		if i==0 {
			p.PalTransparent[i] = 0
		}else{
			p.PalTransparent[i] = 1
		}
	}
	
	p.Font = font.Font("PICO-8.ttf",4)
	
	p.PenColor = 1
	
	p.Resource = make(map[string]string)
	
	p.Version = 8
      
	p.Uptime = int(sdl.GetTicks())
	
	return p
	
}

func (self *Pico8) SyncDrawPal() {
	
	for i:=0;i<16;i++ {
		self.draw_colors[i] = self.pal_colors[ self.DrawPaletteIdx[i]  ]
		
	}
	
	self.DrawPalette.SetColors(self.draw_colors)

}

func (self *Pico8) Flip() {
    
    window   := display.GetWindow()
    _w,_h := window.GetSize()
    window_w := int(_w)
    window_h := int(_h)

    if self.HWND != nil {
	blit_rect := rect.NewRect(self.CameraDx,self.CameraDy)
		
	//surface.Fill(self.DisplayCanvas,color.NewColor(3,5,10,255))
		
	self.DisplayCanvas.SetPalette(self.DisplayPalette)
	surface.Blit(self.DisplayCanvas,self.DrawCanvas,blit_rect,nil)
		
	if window_w > self.Width && window_h > self.Height {
		
		bigger_border := window_w
		if bigger_border > window_h {
			bigger_border = window_h
		}
			
		_blit_x := (window_w - bigger_border)/2
		_blit_y := (window_h - bigger_border)/2
			
		bigger := transform.Scale(self.DisplayCanvas,bigger_border,bigger_border)
			
		_r := rect.NewRect(_blit_x,_blit_y)
		surface.Blit(self.HWND,bigger,_r,nil)
			
	} else {
 		
		_r := rect.NewRect()
		surface.Blit(self.HWND,self.DisplayCanvas,_r,nil)
	}		
		
		//surface.Fill(self.DisplayCanvas,color.NewColor(0,0,0,255))
		
	self.ClipRect = nil
	self.CameraDx = 0
	self.CameraDy = 0
		
    }
    
    self.frames+=1
    self.curr_time = int(sdl.GetTicks())
    if self.curr_time - self.prev_time > 10000 {
	self.fps = self.frames / 10
	fmt.Println("pico8 fps is ",self.fps)
	self.frames = 0
	self.prev_time = self.curr_time
    }
}


func (self *Pico8) color(c ...int) int {

  if len(c) == 0 {
    return self.PenColor
  }
  
  p := c[0]
  if p < 0 {
    return self.PenColor
  }

  if p < 16 && p >=0 {
    self.PenColor = p
  }
  return self.PenColor

}

func (self *Pico8) Cls(args []CmdArg) {
	
	if len(args) == 0 {
		surface.Fill(self.DrawCanvas,color.NewColor(0))
	}
	
	_color := args[0].GetInt()
	
	if _color >=0 && _color < 16 {
		surface.Fill(self.DrawCanvas, color.NewColorSDL(self.draw_colors[self.DrawPaletteIdx[_color ]]) )
	}
	
	self.Cursor=[2]int{0,0}	
}

func (self *Pico8) Print(args []CmdArg){
	
	text := ""
	x:=self.Cursor[0]
	y:=self.Cursor[1]
	c:=1
	
	if len(args) == 0 || len(args) == 2 {
		return
	}
	
	if len(args) == 1 {
		text = args[0].GetStr()
		self.Cursor[1]+=6
	}
	
	if len(args) == 3 {
		text = args[0].GetStr()
		x = args[1].GetInt()
		y = args[2].GetInt()
		self.Cursor = [2]int{x,y}
	}
	
	if len(args) == 4 {
		text = args[0].GetStr()

		x = args[1].GetInt()
		y = args[2].GetInt()
		c = args[3].GetInt()
		
		self.Cursor = [2]int{x,y}
	}
	
	//fmt.Println(text,x,y,c)
	self.color(c)
	
	imgText := font.Render(self.Font,text,false,color.NewColorSDL(self.draw_colors[self.DrawPaletteIdx[self.PenColor]]),nil)
	
	imgText.SetColorKey(true,0)
	
	_r := rect.NewRect(x,y,0,0)
	surface.Blit(self.DrawCanvas,imgText,_r,nil)
}

func (self *Pico8) SetVersion(args []CmdArg) {
  if len(args) < 1 {
    return
  }
  
  version := args[0].GetInt()
  self.Version = version

}

func (self *Pico8) Res(args []CmdArg) {
  
  res_name := ""
  res_content := ""
  if len(args) < 2 {
    return
  }

  res_name = args[0].GetStr()
  res_content = args[1].GetStr()
  
  if res_name == "gff" {
    sprite := 0
    data_array := strings.Split(res_content,"\n")
    for i:=0;i<len(data_array);i++ {
      rowpixel := data_array[i]
      if self.Version <= 2 {
	for j:=0;j<len(rowpixel);j++ {
	  if s, err := strconv.ParseInt(string(rowpixel[j]), 16, 32); err == nil {
	    self.SpriteFlags[sprite] = int(s)
	    sprite+=1
	  }
	}
      }else {
	for j:=0;j<len(rowpixel);j+=2 {
	  if s, err := strconv.ParseInt(string(rowpixel[j])+string(rowpixel[j+1]), 16, 32); err == nil {
	    self.SpriteFlags[sprite] = int(s)
	    sprite+=1
	    if sprite > 255 {
	      break
	    }
	  }
	}
      }
    }
  }

  if res_name =="gfx" {
    col := 0
    row := 0
    data_array := strings.Split(res_content,"\n")
    for i:=0;i<len(data_array);i++ {
      rowpixel:= data_array[i]
      for j:=0;j<len(rowpixel);j++{
	digi := string(rowpixel[j])
	if v, err := strconv.ParseInt(digi, 16, 16); err == nil {
	    _color := color.NewColor(int(v*16),int(v*16),int(v*16),1)
	    self.GfxSurface.Set(col,row,_color)
	    col+=1
	    if col == self.Width {
		col = 0
		row+=1
	    }
	}
      }
    }
    // set_shared_map
    if self.Version > 3 {
      tx :=0
      ty := 32
      shared := 0
      for sy:=64;sy<128;sy++ {
	for sx:=0;sx<128;sx+=2{
	  R1,_,_,_:= self.GfxSurface.At(sx,sy).RGBA()
	  lo := int(R1/255/16)
	  R2,_,_,_ := self.GfxSurface.At(sx+1,sy).RGBA()
	  hi := int(R2/255/16)
	  //fmt.Println(sx,sy,R1,R2,hi,lo)
	  v := (hi << 4) | lo

	  self.MapMatrix[ ty+tx*64 ] = v
	  shared +=1
	  tx +=1
	  if tx == 128 {
	    tx = 0
	    ty+=1
	  }
	}
      }
      fmt.Println("Map Shared: ",shared)
    }

  }

  if res_name == "map" { //set_map
    col :=0
    row :=0
    tiles := 0
    data_array := strings.Split(res_content,"\n")
    for i:=0;i<len(data_array);i++ {
      rowpixel := data_array[i]
      for j:=0;j<len(rowpixel);j+=2 {
	if v, err := strconv.ParseInt(string(rowpixel[j])+string(rowpixel[j+1]), 16, 16); err == nil {
	  
	  self.MapMatrix[row+col*64] = int(v)

	  tiles+=1
	  col +=1
	  if col == self.Width {
	    col = 0
	    row +=1
	  }
	  
	}
      }
    }
    
    fmt.Println("set_map ",tiles)
  }
}

func (self *Pico8) Spr(args []CmdArg) {

  var n,x,y,w,h int // 
  var flip_x,flip_y int
  
  if len(args) == 0 {
    return
  }

  if len(args) > 2 {
    n = args[0].GetInt()
    x = args[1].GetInt()
    y = args[2].GetInt()
  }
  if len(args) > 3 {
    w = args[3].GetInt()
  }
  if len(args) > 4 {
    h = args[4].GetInt()
  }
  
  
  if len(args) >  5{
    flip_x = args[5].GetInt()
  }

  if len(args) > 6 {
    flip_y = args[6].GetInt()
  }
  
  idx := n%16
  idy := n/16
  
  start_x := int(idx*8)
  start_y := int(idy*8)
  
  _w := w*8
  _h := h*8
  _sw := _w
  _sh := _h
  
  if start_x >= self.Width || start_y >= self.Height {
    fmt.Println("spr start_x or start_y illegl")
    return
  }

  if start_x +_w > self.Width {
    _sw = self.Width - start_x
  }

  if start_y + _h > self.Height {
    _sh = self.Height - start_y
  }
        
  if _sw == 0 || _sh == 0 {
    fmt.Println("spr _sw or _sh is zero",args)
    return
  }
  
  gfx_piece := surface.ASurface(_sw,_sh,0,32)
  gfx_piece.SetPalette(self.DrawPalette)
  
  gfx_piece.SetColorKey(true,0)
  
  for _x:=0;_x<_sw;_x++ {
    for _y:=0;_y<_sh;_y++ {
      R,_,_,_ := self.GfxSurface.At(start_x+_x,start_y+_y).RGBA()
      v   := int(R/255/16)
      gfx_piece.Set(_x,_y, color.NewColorSDL(self.draw_colors[v]))
    }
  }

  
  xflip := false
  yflip := false
  
  if flip_x > 0 {
    xflip = true
  }
  if flip_y > 0 {
    yflip = true
  }
  
  gfx_piece_new := transform.Flip(gfx_piece,xflip,yflip)
  
  for i:=0;i<16;i++ {
    if self.PalTransparent[i] == 0 {
      gfx_piece_new.SetColorKey(true,uint32(i))
    }
  }
  
  _r := rect.NewRect(x,y,0,0)
  surface.Blit(self.DrawCanvas,gfx_piece_new,_r,nil)
  
}

func (self *Pico8) draw_map(n,x,y int) {
  idx := n % 16
  idy := n / 16
  start_x := idx*8
  start_y := idy*8
  
  w_ := 8
  h_ := 8
  gfx_piece := surface.ASurface(w_,h_,0,32)
  gfx_piece.SetPalette(self.DrawPalette)  
  gfx_piece.SetColorKey(true,0)
  
  for _x:=0; _x< w_;_x++ {
    for _y:=0; _y <h_;_y++ {
      R,_,_,_ := self.GfxSurface.At(start_x+_x,start_y+_y).RGBA()
      v := int(R/255/16)
      gfx_piece.Set(_x,_y,color.NewColorSDL(self.draw_colors[v]))
    }
  }

  _r := rect.NewRect(x,y,0,0)
  surface.Blit(self.DrawCanvas,gfx_piece,_r,nil)

}


func (self *Pico8) Map(args []CmdArg) {
  
  var cel_x,cel_y,sx,sy,cel_w,cel_h,bitmask int
  if len(args) == 0 {
    return
  }

  if len(args) > 0 {
    cel_x = args[0].GetInt()
  }
  
  if len(args) > 1 {
    cel_y = args[1].GetInt()

  }  
  if len(args) > 2 {
    sx = args[2].GetInt()
  }
  if len(args) > 3 {
    sy = args[3].GetInt()
  }

  if len(args) > 4 {
    cel_w = args[4].GetInt()
  }
  
  if len(args) > 5 {
    cel_h = args[5].GetInt()
  }

  if len(args) > 6 {
    bitmask = args[6].GetInt()
  }

  addr := 0
  for y:=0;y<cel_h;y++ {
    for x:=0;x<cel_w;x++ {
      addr = cel_y + y +(cel_x+x)*64
      if addr < 8192 { //
	v := self.MapMatrix[addr]
	if v > 0 {
	  if bitmask == 0 {
	    self.draw_map(v,sx+x*8,sy+y*8)
	  }else {
	    if self.SpriteFlags[v] & bitmask != 0 {
	      self.draw_map(v,sx+x*8,sy+y*8)
	    }
	  }
	}
      }else {
	fmt.Println("addr >= 8192,exceeds ",addr)
      }
    }
  }
}

func (self *Pico8) MGet(args []CmdArg) string {
  
  var x,y int
  
  if len(args) < 2 {
    return "0"
  }
  
  x = args[0].GetInt()
  y = args[1].GetInt()

  if y > 63 || x < 0 || x > 127 || y < 0 {
    return "0"
  }
  
  return fmt.Sprintf("%d",self.MapMatrix[y+x*64])
}

func (self *Pico8) Circ(args ...CmdArg) {
  var ox,oy,r,col int
  col = -1

  if len(args) < 3 {
    return
  }
  
  if len(args) > 2 {
    ox = args[0].GetInt()
    oy = args[1].GetInt()
    r  = args[2].GetInt()
  }
  if len(args) > 3 {
    col = args[3].GetInt()
  }
  
  self.color(col)

  draw.Circle(self.DrawCanvas, color.NewColorSDL(self.draw_colors[self.PenColor]),ox,oy,r,1)

}




func (self *Pico8) Circfill(args ...CmdArg) {
  var cx,cy,r,col int
  col = -1
  if len(args) < 3 {
    return
  }
  
  if len(args) > 2 {
    cx = args[0].GetInt()
    cy = args[1].GetInt()
    r  = args[2].GetInt()
  }
  if len(args) > 3 {
    col = args[3].GetInt()
  }
  self.color(col)
  draw.Circle(self.DrawCanvas, color.NewColorSDL(self.draw_colors[self.PenColor]),cx,cy,r,0)

}

func (self *Pico8) Line(args ...CmdArg) {
  var x0,y0,x1,y1,col int

  if len(args) < 4 {
    return
  }

  if len(args) > 3 {
    x0 = args[0].GetInt()
    y0 = args[1].GetInt()
    x1 = args[2].GetInt()
    y1 = args[3].GetInt()
  }
  if len(args) > 4 {
    col = args[4].GetInt()
  }
  
  self.color(col)
  
  draw.Line(self.DrawCanvas,color.NewColorSDL(self.draw_colors[self.PenColor]),x0,y0,x1,y1,1)

}


func (self *Pico8) Rect(args []CmdArg) {
  var x0,y0,x1,y1,col int
  col = -1
  if len(args) < 4 {
    return
  }

  if len(args) > 3 {
    x0 = args[0].GetInt()
    y0 = args[1].GetInt()
    x1 = args[2].GetInt()
    y1 = args[3].GetInt()
  }

  if len(args) > 4 {
    col = args[4].GetInt()
  }

  self.color(col)
  rect_ := rect.NewRect(x0+1,y0+1,x1-x0,y1-y0)
  draw.Rect(self.DrawCanvas, color.NewColorSDL(self.draw_colors[self.PenColor]), rect_,1)
  
}

func (self *Pico8) Rectfill(args []CmdArg) {
  var x0,y0,x1,y1,col int
  col = -1
  if len(args) < 4 {
    return
  }

  if len(args) > 3 {
    x0 = args[0].GetInt()
    y0 = args[1].GetInt()
    x1 = args[2].GetInt()
    y1 = args[3].GetInt()
  }

  if len(args) > 4 {
    col = args[4].GetInt()
  }
  
  self.color(col)
  
  w := (x1-x0)+1
  h := (y1-y0)+1
  
  if w < 0 {
    w = -w
    x0 = x0-w
  }

  if h < 0 {
    h = -h
    y0=y0-h  
  }
  rect_ := rect.NewRect(x0,y0,w,h)
  draw.Rect(self.DrawCanvas, color.NewColorSDL(self.draw_colors[self.PenColor]), rect_,0)
  
}

func (self *Pico8) Palt(args ...CmdArg) {
  var c,t int

  if len(args) == 0 {
    for i:=0;i<16;i++ {
      if i==0 {
	self.PalTransparent[i] = 0
      }else{
	self.PalTransparent[i] = 1
      }
    }
  }

  if len(args) == 2 {
    c = c % 16 
    if t == 1 {
      self.PalTransparent[c] = 0
    }else {
      self.PalTransparent[c] = 1
    }
  }
}

func (self *Pico8) Pal(args []CmdArg) {
  var c0,c1,p int
  if len(args) == 0 {
    if self.PaletteModified == false {
      return
    }
    
    for i:=0;i<16;i++ {
      self.DrawPaletteIdx[i] = i
      self.display_colors[i] = self.pal_colors[i]
      self.draw_colors[i]    = self.pal_colors[i]
    }
    
    self.DrawPalette.SetColors(      self.draw_colors)
    self.DisplayPalette.SetColors(   self.display_colors)
    
    self.Palt()
    self.SyncDrawPal()
    self.DisplayCanvas.SetPalette(self.DisplayPalette)
    self.DrawCanvas.SetPalette(   self.DrawPalette) 
  
    self.PaletteModified = false
  }
  
  if len(args) == 3 {
    c0 = args[0].GetInt()
    c1 = args[1].GetInt()
    p  = args[2].GetInt()
    
    c0 = c0 % 16
    c1 = c1 % 16
    
    if p == 1 {
      self.display_colors[c0] = self.pal_colors[c1]
      self.PaletteModified = true
      self.DisplayPalette.SetColors(   self.display_colors)
      self.DisplayCanvas.SetPalette(self.DisplayPalette)
    }
    if p == 0 {
      self.DrawPaletteIdx[c0] = c1
      self.PaletteModified = true
      self.SyncDrawPal()
      self.DrawCanvas.SetPalette(   self.DrawPalette) 
    }

  }
}

