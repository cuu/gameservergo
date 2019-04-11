local UDP = {}
local socket = require("socket")
local udp = assert(socket.udp())

function UDP.init()
	assert(udp:setpeername("127.0.0.1",8080))
	udp:settimeout(0.1)
	udp:send("ping")
end

function UDP.data()
	local data
  data = udp:receive()
	return data
end

return UDP
