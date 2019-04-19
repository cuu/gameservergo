require 'strict'

local _player = {
    [1] = {[0]='Left',[1]=-1},
    [2] = {[0]='Right',[1] = -1},
    [3] = {[0]='Up', [1] = -1},
    [4] = {[0]='Down',[1] = -1},
    [5] = {[0]='U', [1] = -1 },
    [6] = {[0]='I', [1] = -1 },
    [7] = {[0]='Return',[1]=-1},
    [8] = {[0]='Escape',[1]=-1},
}

local __keymap = {
	[0] = _player
}

local coroutine_scheduler = require("coroutine_scheduler")
sched = coroutine_scheduler.Scheduler()

socket = require("socket")

function sleep(sec)
    socket.sleep(sec)
end

function time_ms() 
	return socket.gettime()*1000
end

--local api = require 'libpico8_unix_socket'

local server = require("server")
local api = require("libpico8")

local remote_host = "127.0.0.1"

api.server = server 

local UDP = { data = {} ,curr_time = 0}
local udp = assert(socket.udp())

function UDP.connect()
	assert(udp:setpeername(remote_host,8080))
	udp:settimeout()
	udp:send("ping")
end

function UDP.send() -- must inside lua's coroutine
  local ret,msg
  local content = ''
  
  -- print("safe_tcp_send data is " ,data ,#data)
  if #UDP.data == 0 then 
    print("data is zero",UDP.data)
    return nil
  end

  local piece = {}
  local divid = 3
  if #UDP.data % 3 == 0 then 
    divid = 3
  end
  
  if #UDP.data % 4 == 0 then 
    divid = 4
  end  
  
  for i=1,#UDP.data,divid do 
    
    for j=1,divid do
      if UDP.data[i+j-1] ~= nil then 
        piece[j] = UDP.data[i+j-1]
      end
    end
    
    content = table.concat(piece,"|")
    ret,msg = udp:send(content.."\n")

  end
  
 
  UDP.data = {}
  
  return nil
end

function UDP.cache(data)
  --UDP.data = UDP.data..data.."\n"
  UDP.data[#UDP.data+1] = data
  
end

function UDP_SendLoop()


  local now

  while true do
	  now = time_ms()
	  if now - UDP.curr_time >= 33 then
	    UDP.send()
	    UDP.curr_time = now
	  end

	  sched:suspend(udp)
  end

end


UDP.connect()

server.Network = UDP

local TCP= {}

local tcp = assert(socket.tcp())

function TCP.connect()
  local host, port = remote_host, 8080
  tcp:connect(host, port);
  tcp:settimeout(5)
end

function TCP.send(data)
  local ret,msg
  local ret2
  -- print("safe_tcp_send data is " ,data ,#data)
  if #data == 0 then 
    print("data is zero",data)
    return nil
  end
  
  data = data.."\n"
  
  ret,msg = tcp:send(data)
  if(ret == nil) then
	  print("exiting...",msg)
  	os.exit()
  end
  
	return nil
end

TCP.connect()

server.NetworkTCP = TCP

log = print

local function setfenv(fn, env)
  local i = 1
  while true do
    local name = debug.getupvalue(fn, i)
    if name == "_ENV" then
      debug.upvaluejoin(fn, i, (function()
        return env
      end), 1)
      break
    elseif not name then
      break
    end

    i = i + 1
  end

  return fn
end

local function getfenv(fn)
  local i = 1
  while true do
    local name, val = debug.getupvalue(fn, i)
    if name == "_ENV" then
      return val
    elseif not name then
      break
    end
    i = i + 1
  end
end

function new_sandbox()
  return {

    clip=api.clip,
    pget=api.pget,
    pset=api.pset,

    fget=api.fget,    
    fset=api.fset,

    flip=api.flip,
    print=api.print,
    printh=log,
    cursor=api.cursor,
    color=api.color,
    cls=api.cls,
    camera=api.camera,
    circ=api.circ,
    circfill=api.circfill,
    line=api.line,
    rect=api.rect,
    rectfill=api.rectfill,
    
    pal=api.pal,
    palt=api.palt,
    spr=api.spr,
    sspr=api.sspr,
    add=api.add,
    del=api.del,
    foreach=api.foreach,
    count=api.count,
    all=api.all,
    btn=api.btn,
    btnp=api.btnp,
    sfx=api.sfx,
    music=api.music,
    
    mget=api.mget,
    mset=api.mset,
    map=api.map,

    max=api.max,
    min=api.min,
    mid=api.mid,
    flr=api.flr,
    cos=api.cos,
    sin=api.sin,
    atan2=api.atan2,
    sqrt=api.sqrt,
    abs=api.abs,
    rnd=api.rnd,
    srand=api.srand,
    sgn=api.sgn,
    band=api.band,
    bor=api.bor,
    bxor=api.bxor,
    bnot=api.bnot,
    shl=api.shl,
    shr=api.shr,
    exit=api.shutdown,
    shutdown=api.shutdown,
    sub=api.sub,
    stat=api.stat,
    time = api.time, 
		reboot = api.reboot,
		printh = api.printh,
		tostr  = api.tostr,
    mapdraw = api.map,
		run    = api.run,
		string = string,
    flip_network = api.flip_network

   }
end


function set_keymap(data,keymap)

	local data_array
	local key 
	local action
	local pos
  
  pos = data:find(",")
  if pos ~= nil then
    key = data:sub(0,pos-1)
    action = data:sub(pos+1,#data-1)
  end
  
  for i,v in ipairs(keymap[0]) do
		if v[0] == key then
			if action == "Down" then
				keymap[0][i][1] = 1
			end
			if action == "Up" then
				keymap[0][i][1] = -1
			end
      break
		end
	end 
  
  return nil

end


function GetBtnLoop()
  local count = 0 
  local framerate = 1/api.pico8.fps
  udp:send("ping")
  while true do
    udp:settimeout( framerate )
    local s, status = udp:receive()
    if s ~= nil then
      --count = count + string.len(s)
      --print("received: ",s)
      set_keymap(s,__keymap)
    end
				
    if status == "timeout" then
      sched:suspend(udp)
    end
    if status == "closed" then
      print("closed....")
      break
    end
  end
end
	
function draw(cart)

  local frames = 0
  local frame_time = 1/api.pico8.fps
  local prev_time = 0
  local curr_time = 0
  local skip = 0  
  
  curr_time = time_ms()
  prev_time = time_ms()
  local frame_time_ms = frame_time*1000
  
  while true do

    if cart._update then cart._update() end
    if cart._update60 then cart._update60() end

    if cart._draw   then cart._draw() end

      curr_time = time_ms()

      if curr_time - prev_time < frame_time_ms then
	sleep(frame_time)
      end

      api.flip()
      api.flip_network()
      prev_time = curr_time
      sched:suspend(udp)
    end

end

function api.btn(i,p)
	local thing
	local ret
	
	if type(i) == 'number' then
		i = i +1
		p = p or 0
		if __keymap[p] and __keymap[p][i] then
			ret =  __keymap[p][i][1]
			if ret >= 0 then
				return true
			else
				return false
			end
		end
	end
	
	return false

end

function api.run()
  
end

function api.flip_network()

  UDP.send()

end

function RunLoop(file)
  
  api.load_p8_text(file)
  
	local cart = new_sandbox()
	local ok,f,e = pcall(load,api.loaded_code)
  if not ok or f==nil then
    log('=======8<========')
    log(api.loaded_code)
    log('=======>8========')
    error('Error loading lua: '..tostring(e))
  else
		local result
		setfenv(f,cart)
    ok,result = pcall(f)
    if not ok then
      error('Error running lua: '..tostring(result))
    else
      log('lua completed')
    end

		if cart._init then cart._init() end
    
		 draw(cart)
  end
end



--------------------------------------------
function main(file)

  sched:spawn(RunLoop,file)
  
  sched:spawn(GetBtnLoop)
  --sched:spawn(UDP_SendLoop)

  while true do
    local worked, t = sched:select()
    if worked then
        if t and t ~= 0 then
	    print(t)
            if socket then socket.sleep(t) end
        end
    else
        print(worked, t)
        break
    end
  end


end

if #arg > 1 then
	main(arg[2])
else
	log("No arguments")
	log(arg[0])
end

