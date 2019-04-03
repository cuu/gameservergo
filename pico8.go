package main
import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"

	"github.com/cuu/gogame/surface"
	"github.com/cuu/gogame/rect"

	"github.com/cuu/gogame/color"
	"github.com/cuu/gogame/font"
	"github.com/cuu/gogame/display"
	//"github.com/cuu/gogame/transform"
	"github.com/cuu/gogame/draw"
	//"github.com/cuu/gogame/event"



)

type Pico8 struct{
	Width int
	Height int
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
	
	p.DisplayCanvas = surface.ASurface(p.Width,p.Height,0,8)
	p.DrawCanvas    = surface.ASurface(p.Width,p.Height,0,8)
	p.GfxSurface    = surface.ASurface(p.Width,p.Height,0,24)
	
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
	
	
	return p
	
}

func (self *Pico8) SyncDrawPal() {
	
	for i:=0;i<16;i++ {
		self.draw_colors[i] = self.pal_colors[ self.DrawPaletteIdx[i]  ]
		
	}
	
	self.DrawPalette.SetColors(self.draw_colors)

}

func (self *Pico8) Flip() {
	if self.HWND != nil {
		fmt.Println("pico8 flip now")
		blit_rect := rect.NewRect(self.CameraDx,self.CameraDy)
		/*
		window_surface := display.GetSurface()
		window_w := surface.GetWidth(window_surface)
		window_h := surface.GetHeight(window_surface)
		*/
		surface.Fill(self.DisplayCanvas,color.NewColor(3,5,10,255))
		
		self.DisplayCanvas.SetPalette(self.DisplayPalette)
		
		
		 screen := display.GetSurface()
	      draw.Line(screen,color.NewColor(255,44,255,255), 0,100, 320,100,3)
	      draw.Line(screen,color.NewColor(255,44,255,255), 10, 0, 10,250,4)

	      rect2 := rect.Rect(3,120,200,30)
	      draw.AARoundRect(screen,&rect2,&color.Color{0,213,222,255},10,0, &color.Color{0,213,222,255})
      
		
		surface.Blit(self.DisplayCanvas,self.DrawCanvas,blit_rect,nil)
		
		/*
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
			
		} else {*/
			_r := rect.NewRect()
			surface.Blit(self.HWND,self.DisplayCanvas,_r,nil)
		//}		
		
		surface.Fill(self.DisplayCanvas,color.NewColor(0,0,0,255))
		
		self.ClipRect = nil
		self.CameraDx = 0
		self.CameraDy = 0
		
	}
}


func (self *Pico8) Color(c ...int) int {
	if len(c) == 0 {
		return self.PenColor
	}
	
	p := c[0]
	
	if p < 16 && p >=0 {
		self.PenColor = p
	}
	return self.PenColor
	
}

func (self *Pico8) Cls(color_index ...int) {
	if len(color_index) == 0 {
		surface.Fill(self.DrawCanvas,color.NewColor(0))
	}
	
	_color := color_index[0]
	
	if _color >=0 && _color < 16 {
		surface.Fill(self.DrawCanvas, color.NewColorSDL(self.draw_colors[self.DrawPaletteIdx[_color ]]) )
	}
		
	self.Cursor=[2]int{0,0}	
}

func (self *Pico8) Print(args []CmdArg){
	
	text := ""
	x:=self.Cursor[0]
	y:=self.Cursor[1]
	c:=0
	
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
	
	fmt.Println(text,x,y,c)
	c = 2
	self.Color(c)
	
	imgText := font.Render(self.Font,text,false,color.NewColorSDL(self.draw_colors[self.DrawPaletteIdx[self.PenColor]]),nil)
	
	imgText.SetColorKey(true,0)
	
	_r := rect.NewRect(x,y,0,0)
	surface.Blit(self.DrawCanvas,imgText,_r,nil)
}


