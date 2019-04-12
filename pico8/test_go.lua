local go = require("goroutines")

socket = require("socket")

function sleep(sec)
    socket.sleep(sec)
end

function time_ms() 
	return socket.gettime()*1000
end

function test() 
	while true do
		print("hi")
		sleep(1)
	end
end


go.goEnv.go(test)

while true do
	print("yeah")
	sleep(0.5)
end


