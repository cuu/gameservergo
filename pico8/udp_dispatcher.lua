local unpack = table.unpack

local coroutine_scheduler = require("coroutine_scheduler")

sched = coroutine_scheduler.Scheduler()

socket = require("socket")

function sleep(sec)
    socket.sleep(sec)
end


local UDP = {}
local udp = assert(socket.udp())
local remote_ip = "127.0.0.1"

function UDP.connect()
	assert(udp:setpeername(remote_ip,8080))
	udp:settimeout()
--	udp:send("ping")
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
  
  sched:suspend(udp)
  
	return nil
end

UDP.connect()

function upload(host)
		local count = 0

		while true do
			UDP.send("ping"..tostring(count))
			count = count+1
		end
end

function download(host)
			
      local count = 0    -- counts number of bytes read	
			local framerate = 1/120
      udp:send("ping")
			while true do
				udp:settimeout( framerate )
	      local s, status = udp:receive()
				if s ~= nil then
	  	    count = count + string.len(s)
  	 			print("received: ",s)
				end
				
				if status == "timeout" then
					sched:suspend(udp)
				end
				if status == "closed" then
					print("closed....")
					break
				end
			end
      print(count)
end



sched:spawn(download)
sched:spawn(upload)

while true do
    local worked, t = sched:select()
    if worked then
        if t and t ~= 0 then
            if socket then socket.sleep(t) end
        end
--				if t and t == 0 then
--					  if socket then socket.sleep(1/60) end
--				end
    else
        print(worked, t)
        break
    end
end

