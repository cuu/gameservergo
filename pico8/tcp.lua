
local TCP = {}

local host, port = "127.0.0.1", 8080
local socket = require("socket")
local tcp = assert(socket.tcp())

function TCP.connect()
  tcp:connect(host, port);
  tcp:settimeout(5)
end

--local host = "/tmp/gs"
--local socket = require("socket")
--socket.unix =  require("socket.unix")
--local tcp = assert(socket.unix())
--function TCP.connect()
--  tcp:connect(host)
--  tcp:settimeout(5)
--end
function TCP.get()
	local ret
	ret,msg = tcp:receive("*l")
	if ret == nil then
		print(msg)
	end
	return ret
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

return TCP



