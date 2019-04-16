json = require("json")

local server = {Network=nil,NetworkTCP=nil}

function safe_format_json(funcname,...)
  local ret = {Fc=funcname}
  local args = {}
 
  for i, v in ipairs{...} do
    local arg = {Typ="",Val=""}
    if type(v) == "string" then
      arg.Typ="S"
    end 
    if type(v) == "number" then
      arg.Typ="I"
    end
    if type(v) == "boolean" then
      arg.Typ="B"
    end
    arg.Val=v

    table.insert(args,arg)
  end 
  
  if #args > 0 then
    ret.Ags = args
  end 

  return json.encode(ret)
end

-- lisp style 
function safe_format(funcname,...) 
  local args = {}
  local arg

  table.insert(args,funcname)
  for i,v in ipairs{...} do
    if type(v) == "string" then
      arg = "\""..v.."\""
    end
    if type(v) == "number" then
      arg = tostring(v)
    end
    if type(v) == "boolean" then
      arg = tostring(v)
    end
    table.insert(args,arg)
  end

  if #args > 0 then
    return "("..table.concat(args," ")..")"
  end

  return ""

end

function server.down()

end

function server.scroll(dy)
	dy = dy or 0
	local thing = safe_format("scroll",dy)
	return server.Network.cache(thing)
end

function server.print(str,x,y,col)
	local thing
	if x == nil then
		thing = safe_format("print",str)
	end

	if x ~= nil and col == nil  then 
		thing = safe_format("print",str,math.floor(x),math.floor(y))
	end

	if x ~= nil and col ~= nil then
		thing = safe_format("print",str,math.floor(x),math.floor(y),math.floor(col))
	end

	return server.Network.cache(thing)

end

function server.cls(frame)
	local thing = safe_format("cls",frame)
	return server.Network.cache(thing)

end

function server.flip()
	local thing = safe_format("flip")
	return server.Network.cache(thing)
end

function server.btn(codestr,playernumber)
	local thing = safe_format("btn", codestr,playernumber)
	return server.Network.send(thing)
end

function server.btnp(codestr,playernumber)
	local thing = safe_format("btnp", codestr,playernumber)
	return server.Network.send(thing)
end

function server.sspr(sx,sy,sw,sh,dx,dy,dw,dh,flip_x,flip_y)
  local thing = safe_format("sspr", sx,sy,sw,sh,dx,dy,dw,dh,flip_x,flip_y)
  return server.Network.cache(thing)
end

function server.spr(n,x,y,w,h,flip_x,flip_y)
  local thing = safe_format("spr", n,x,y,w,h,flip_x,flip_y)
  return server.Network.cache(thing)

end

function server.map(cel_x,cel_y,sx,sy,cel_w,cel_h,bitmask)
  local thing = safe_format("map",cel_x,cel_y,sx,sy,cel_w,cel_h,bitmask)
  return server.Network.cache(thing)
end

function server.color(c)
  local thing = safe_format("color",c)
  return server.Network.cache(thing)
end

function server.pset(x,y,c)
  local thing = safe_format("pset",x,y,c)
  return server.Network.send(thing)
end

function server.cursor(x,y)
  local thing = safe_format("cursor",x,y,c)
  return server.Network.cache(thing)
end

function server.mget(x,y)
 local thing = safe_format("mget",x,y)
 return server.Network.send(thing)
end

function server.mset(x,y,v)
 local thing = safe_format("mset",x,y,v)
 server.Network.send(thing)
end

function server.rect(x0,y0,x1,y1,col)
  local thing

  if col == nil then
    thing = safe_format("rect",x0,y0,x1,y1)
  else
    thing = safe_format("rect",x0,y0,x1,y1,col)
  end
  server.Network.cache(thing)
end

function server.rectfill(x0,y0,x1,y1,col)
  local thing
  x0 = math.floor(x0)
  y0 = math.floor(y0)
  x1 = math.floor(x1)
  y1 = math.floor(y1)

  if col == nil then
    thing = safe_format("rectfill",x0,y0,x1,y1)
  else
    col = math.floor(col)
    thing = safe_format("rectfill",x0,y0,x1,y1,col)
  end
  server.Network.cache(thing)
end

function server.circ(ox,oy,r,col)
  local thing
  ox = math.floor(ox)
  oy = math.floor(oy)
  r =  math.floor(r)

  if col == nil then
    thing = safe_format("circ",ox,oy,r)
  else
    col = math.floor(col)
    thing = safe_format("circ",ox,oy,r,col)
  end
  server.Network.cache(thing)
end

function server.circfill(cx,cy,r,col)
  local thing
  cx = math.floor(cx)
  cy = math.floor(cy)
  r = math.floor(r)

  if col == nil then
    thing = safe_format("circfill",cx,cy,r)
  else
    col = math.floor(col)
    thing = safe_format("circfill",cx,cy,r,col)
  end
  server.Network.cache(thing)
end

function server.line(x0,y0,x1,y1,col)
  local thing
  if col == nil then
    thing = safe_format("line",x0,y0,x1,y1)
  else
    thing = safe_format("line",x0,y0,x1,y1,col)
  end
  server.Network.cache(thing)
end

function server.time()
  local thing = safe_format("time")
  return server.Network.send(thing)
end

function server.pal(c0,c1,p)
  local thing
  if type(c0) ~= 'number' then
    thing=safe_format("pal")
	end

  if c1 ~= nil then
    thing = safe_format("pal",c0,c1,p)
  end

  server.Network.cache(thing)

end

function server.palt(c,t)
  local thing
  if type(c) ~= 'number' then
   thing=safe_format("palt")

  else
   t = t or false
   if t == true then
     thing = safe_format("palt",c,1)
   else
     thing = safe_format("palt",c,0)
   end
  end

  server.Network.cache(thing)

end

function server.fget(n,f)
  local thing
  if f == nil then
    thing = safe_format("fget",n)
  else
    thing = safe_format("fget",n,f)
  end
 	
	local ret = server.Network.send(thing)
	
	return ret 
end

function server.reboot()
	local thing = safe_format("reboot")
	server.Network.cache(thing)
end

function server.clip(x,y,w,h)
	local thing
	if type(x) == 'number' then
		thing = safe_format("clip",x,y,w,h)
	else
		thing = safe_format("clip")
	end

	server.Network.cache(thing)

end

function server.restore_camera(x,y)
	local thing = safe_format("restore_camera",x,y)
	server.Network.cache(thing)
end


function server.printh(text)
	local thing = safe_format("printh",text)
	server.Network.cache(thing)
end

function server.music(n,fade_len,channel_mask)

end
function server.sfx(n,channel,mask)
end

function server.send_pico8_version(version)
  local thing = safe_format_json("pico8", version)
  server.NetworkTCP.send(thing)
end

function server.send_resource_done()
  local thing = safe_format_json("resdone")
  server.NetworkTCP.send(thing)

end

function server.send_resource(res_type,res_data)
  if res_data == nil or #res_data == 0  then 
    return 
  end
  
  local thing = safe_format_json("res", res_type,res_data)
  server.NetworkTCP.send(thing)

end



return server
