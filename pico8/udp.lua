local UDP = {}
local socket = require("socket")
local udp = assert(socket.udp())

function UDP.connect()
	assert(udp:setpeername("127.0.0.1",8080))
	udp:settimeout()
	udp:send("ping")
end


function UDP.send(data)
  local ret,msg
  local ret2
  -- print("safe_tcp_send data is " ,data ,#data)
  if #data == 0 then 
    print("data is zero",data)
    return nil
  end
  
  data = data.."\n"
  
  ret,msg = udp:send(data)
  if(ret == nil) then
	  print("exiting...",msg)
  	os.exit()
  end
  
	return nil
end



return UDP
