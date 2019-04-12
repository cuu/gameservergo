-- Goroutines for ComputerCraft!
-- Made by 1lann (Jason Chu)
-- Last updated: 31st July 2015

--[[
Licensed under the MIT License:
The MIT License (MIT)

Copyright (c) 2015 1lann (Jason Chu)

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
]]--


-- Goroutine manager variables

local activeGoroutines = {}
local channels = {}
local termCompatibilityMode = false
local quitDispatcherEvent = "quit_goroutine_dispatcher"
local channelEventHeader = "goroutine_channel_event_"
local waitGroupEvent = "goroutine_wait_group_event"
local goEnv = {}
local nativeTerm = term.current()
local dispatcherRunning = false
local dispatcherQuitFunc = function() end

local currentGoroutine = 0

-- Goroutine utility functions

local function keyOfGoroutineId(id)
	for k, v in pairs(activeGoroutines) do
		if v.id == id then
			return k
		end
	end

	return nil
end

local function currentKey()
	local key = keyOfGoroutineId(currentGoroutine)
	if key == nil then
		return goEnv.error("Cannot store term context outside of goroutine.", 1)
	end

	return key
end

-- Stacktracer
-- A very small portion was taken from CoolisTheName007
-- Origin: http://pastebin.com/YWwLUUpk

local function stacktrace(depth, isInvoke)
	local trace = {}
	local i = depth + 2
	local first = true

	while true do
		i = i + 1
		_, err = pcall(error, "" , i)
		if err:match("^[^:]+") == "bios" or
			#err == 0 then
			break
		end

		if first then
			first = false
			if isInvoke then
				table.insert(trace, "created by " .. err:sub(1, -3))
			else
				table.insert(trace, "at " .. err:sub(1, -3))
			end
		else
			table.insert(trace, "from " .. err:sub(1, -3))
		end
	end

	if currentGoroutine == -1 then
		table.insert(trace, "created by ? (trace unavailable)")
	else
		local goTrace
		for _, v in pairs(activeGoroutines) do
			if v.id == currentGoroutine then
				goTrace = v.stacktrace
			end
		end

		if not goTrace then
			table.insert(trace, "created by ? (trace unavailable)")
		else
			for _, v in pairs(goTrace) do
				table.insert(trace, v)
			end
		end
	end

	return trace
end

function goEnv.error(err, depth)
	if not depth then
		depth = 1
	end

	if currentGoroutine == -1 then
		return error(err, depth)
	end

	local trace = stacktrace(depth)
	table.insert(activeGoroutines[currentKey()].errors, {
		err = err,
		trace = trace,
	})

	if #(activeGoroutines[currentKey()].errors) > 10 then
		table.remove(activeGoroutines[currentKey()].errors, 1)
	end

	return error(err, depth)
end

local function traceError(err)
	local location, msg = err:match("([^:]+:%d+): (.+)")

	local recentErrors = activeGoroutines[currentKey()].errors

	local tracedTrace = nil
	for _, v in pairs(recentErrors) do
		if v.err == msg then
			tracedTrace = v.trace
			break
		end
	end

	local _, y = nativeTerm.getCursorPos()
	nativeTerm.setCursorPos(1, y)
	nativeTerm.setTextColor(colors.red)

	if tracedTrace then
		print("goroutine runtime error:")
		print(msg)
		for _, v in pairs(tracedTrace) do
			print("  " .. v)
		end

		return
	end

	if not msg then
		msg = err
	end

	print("goroutine runtime error:")
	print(msg)

	if location then
		print("  at " .. location)
	else
		print("  at ? (location unavailable)")
	end

	print("  from ? (trace unavailable)")

	if currentGoroutine == -1 then
		print("  created by ? (trace unavailable)")
		print("  from goroutine start")
	else
		local goTrace
		for _, v in pairs(activeGoroutines) do
			if v.id == currentGoroutine then
				goTrace = v.stacktrace
			end
		end

		if not goTrace then
			print("  created by ? (trace unavailable)")
			print("  from goroutine start")
		else
			for _, v in pairs(goTrace) do
				print("  " .. v)
			end
		end
	end
end

-- Term wrapper to allow for saving terminal states

local emulatedTerm = {}

for k, v in pairs(nativeTerm) do
	emulatedTerm[k] = v
end

emulatedTerm.setCursorBlink = function(blink)
	activeGoroutines[currentKey()].termState.blink = blink
	return nativeTerm.setCursorBlink(blink)
end

emulatedTerm.setBackgroundColor = function(color)
	if type(color) ~= "number" then
		return goEnv.error("Argument to term.setBackgroundColor must be a number")
	end

	activeGoroutines[currentKey()].termState.bg_color = color
	return nativeTerm.setBackgroundColor(color)
end

emulatedTerm.write = function(text)
	goEnv.emitChannel("term_events", "write")
	return nativeTerm.write(text)
end

emulatedTerm.setTextColor = function(color)
	if type(color) ~= "number" then
		return goEnv.error("Argument to term.setTextColor must be a number")
	end

	activeGoroutines[currentKey()].termState.txt_color = color
	return nativeTerm.setTextColor(color)
end

emulatedTerm.scroll = function(...)
	goEnv.emitChannel("term_events", "scroll")
	return nativeTerm.scroll(...)
end

emulatedTerm.clearLine = function(...)
	goEnv.emitChannel("term_events", "clearLine")
	return nativeTerm.clearLine(...)
end

emulatedTerm.clear = function(...)
	goEnv.emitChannel("term_events", "clear")
	return nativeTerm.clear(...)
end

local function restoreTermState(termState)
	nativeTerm.setTextColor(termState.txt_color)
	nativeTerm.setBackgroundColor(termState.bg_color)
	nativeTerm.setCursorBlink(termState.blink)
	nativeTerm.setCursorPos(unpack(termState.cursor_pos))
end

function goEnv.goroutineId()
	return currentGoroutine
end

-- Goroutines and channel functions

function goEnv.go(...)
	local args = {...}

	if type(args[1]) ~= "function" then
		return goEnv.error("First argument to go must be a function.")
	end

	local funcArgs = {}

	if #args > 1 then
		for i = 2, #args do
			table.append(funcArgs, args[i])
		end
	end

	local idsInUse = {}
	for k,v in pairs(activeGoroutines) do
		idsInUse[tostring(v.id)] = true
	end

	local newId = -1
	for i = 1, 1000000 do
		if not idsInUse[tostring(i)] then
			newId = i
			break
		end
	end

	if newId < 0 then
		return goEnv.error("Reached goroutine limit, cannot spawn new goroutine")
	end

	activeGoroutines[currentKey()].termState.cursor_pos =
		{nativeTerm.getCursorPos()}

	local parent = activeGoroutines[currentKey()]
	local copyTermState = {}

	for k, v in pairs(parent.termState) do
		copyTermState[k] = v
	end

	local func = args[1]

	local env = {}

	local global = getfenv(0)
	local localEnv = getfenv(1)

	for k,v in pairs(global) do
		env[k] = v
	end

	for k,v in pairs(localEnv) do
		env[k] = v
	end

	for k, v in pairs(goEnv) do
		env[k] = v
	end

	local newFunc = setfenv(func, env)

	table.insert(activeGoroutines, {
		id = newId,
		func = coroutine.create(newFunc),
		arguments = funcArgs,
		stacktrace = stacktrace(1, true),
		termState = copyTermState,
		suspended = false,
		forceResume = false,
		errors = {},
		filter = nil,
		firstRun = true,
	})

	return newId
end

function goEnv.suspend(goroutineId)
	local key = keyOfGoroutineId(goroutineId)
	if key == nil then
		return goEnv.error("Attempt to suspend non-existent goroutine.")
	end

	activeGoroutines[key].suspended = true
	activeGoroutines[key].forceResume = false
end

function goEnv.resume(goroutineId)
	local key = keyOfGoroutineId(goroutineId)
	if key == nil then
		return goEnv.error("Attempt to resume non-existent goroutine.")
	end

	activeGoroutines[key].suspended = false
	activeGoroutines[key].forceResume = true
end

function goEnv.emitChannel(channel, data, wait)
	if type(channel) ~= "string" then
		return goEnv.error("First argument to emitChannel must be a string.")
	end

	if data == nil then
		return goEnv.error("Second argument (data) to emitChannel, cannot be nil.")
	end

	if channel == "term_events" and termCompatibilityMode and data then
		return
	end

	channels[channel] = {data, goEnv.goroutineId()}

	os.queueEvent(channelEventHeader .. channel)

	if wait then
		while true do
			if channels[channel] == nil then
				return
			end

			coroutine.yield(channelEventHeader .. channel)
		end
	end
end

function goEnv.waitChannel(channel, allowPrev, timeout)
	if type(channel) ~= "string" then
		return goEnv.error("First argument to waitChannel must be a string.")
	end

	if timeout and type(timeout) ~= "number" then
		return goEnv.error("Third argument to waitChannel must be a number or nil.")
	end

	local stillAlive = true

	if timeout then
		goEnv.go(function()
			goEnv.sleep(timeout)
			if stillAlive then
				goEnv.emitChannel(channel, false)
			end
		end)
	end

	if not allowPrev then
		channels[channel] = nil
	end

	os.queueEvent(channelEventHeader .. channel)

	while true do
		if channels[channel] ~= nil then
			stillAlive = false
			local value = channels[channel]
			channels[channel] = nil
			return unpack(value)
		end

		coroutine.yield(channelEventHeader .. channel)
	end
end

-- Custom overrides

function goEnv.sleep(sec)
	local timer = os.startTimer(sec)
	local start = os.clock()
	while true do
		local event, timerId = os.pullEvent()

		if event == "timer" and timerId == timer then
			return
		end

		if os.clock() - start >= sec then
			return
		end
	end
end

-- Wait groups

goEnv.WaitGroup = {}
goEnv.WaitGroup.__index = goEnv.WaitGroup

function goEnv.WaitGroup.new()
	local self = setmetatable({}, goEnv.WaitGroup)
	self:setZero()
	return self
end

function goEnv.WaitGroup:setZero()
	self.incrementer = 0
	os.queueEvent(waitGroupEvent)
end

function goEnv.WaitGroup:done()
	if self.incrementer > 0 then
		self.incrementer = self.incrementer - 1
		os.queueEvent(waitGroupEvent)
	end
end

function goEnv.WaitGroup:wait()
	while true do
		if self.incrementer == 0 then
			return
		end

		coroutine.yield(waitGroupEvent)
	end
end

function goEnv.WaitGroup:add(amount)
	self.incrementer = self.incrementer + amount
	if self.incrementer < 0 then
		self.incrementer = 0
	end
	os.queueEvent(waitGroupEvent)
end

function goEnv.WaitGroup:value()
	return self.incrementer
end

-- Runner

local function cleanUp()
	channels = {}
	activeGoroutines = {}
	dispatcherRunning = false
	termCompatibilityMode = false
	term.redirect(nativeTerm)

	local ret, err = pcall(dispatcherQuitFunc)
	if not ret then
		local _, y = term.getCursorPos()
		term.setCursorPos(1, y)
		term.setTextColor(colors.red)

		print("user dispatcher quit error:")
		print(err)

		term.setTextColor(colors.white)
	end
end

function runDispatcher(programFunction)
	if dispatcherRunning then
		error("Dispatcher already running.")
	end
	dispatcherRunning = true

	term.redirect(emulatedTerm)

	local env = {}
	local global = getfenv(0)
	local localEnv = getfenv(1)

	for k, v in pairs(global) do
		env[k] = v
	end

	for k, v in pairs(localEnv) do
		env[k] = v
	end

	for k, v in pairs(goEnv) do
		env[k] = v
	end

	local main = setfenv(programFunction, env)

	table.insert(activeGoroutines, {
		func = coroutine.create(main),
		arguments = {},
		id = 0,
		termState = {
			txt_color = colors.white,
			bg_color = colors.black,
			blink = false,
			cursor_pos = {nativeTerm.getCursorPos()},
		},
		errors = {},
		stacktrace = {"from dispatcher start"},
		suspended = false,
		forceResume = false,
		filter = nil,
		firstRun = true,
	})

	local events = {}

	while true do
		for k, v in pairs(activeGoroutines) do
			if coroutine.status(v.func) ~= "dead" and not v.suspended then
				if (v.filter and events and #events > 0 and
					events[1] == v.filter) or (not v.filter and
					((events and #events > 0 and
					not events[1]:find("goroutine_channel_event_")) or
					(not events or #events <= 0))) or
					(events and #events > 0 and events[1] == "terminate") or
					v.forceResume then

					if v.forceResume then
						if v.filter and events[1] ~= v.filter then
							events = {"goroutine_force_resume"}
						end

						activeGoroutines[k].forceResume = false
					end

					activeGoroutines[k].filter = nil
					local resp, err

					currentGoroutine = v.id

					restoreTermState(v.termState)

					if v.firstRun then
						resp, err = coroutine.resume(v.func, unpack(v.arguments))
						activeGoroutines[k].firstRun = false
					else
						resp, err = coroutine.resume(v.func, unpack(events))
					end

					activeGoroutines[k].termState.cursor_pos =
						{nativeTerm.getCursorPos()}

					if resp then
						if err == quitDispatcherEvent then
							cleanUp()
							return
						end

						if type(err) == "string" then
							activeGoroutines[k].filter = err
						end
					else
						traceError(err)
						cleanUp()
						return
					end

					currentGoroutine = -1
				end
			end
		end

		local sweeper = {}

		for k,v in pairs(activeGoroutines) do
			if coroutine.status(v.func) ~= "dead" then
				table.insert(sweeper, v)
			end
		end

		activeGoroutines = {}

		for k,v in pairs(sweeper) do
			if v.termState.blink then
				nativeTerm.setCursorBlink(true)
				nativeTerm.setTextColor(v.termState.txt_color)
				nativeTerm.setCursorPos(unpack(v.termState.cursor_pos))
			end
			table.insert(activeGoroutines, v)
		end

		if #activeGoroutines == 0 then
			cleanUp()
			return
		end

		events = {os.pullEventRaw()}

		if events and #events > 0 and events[1] == quitDispatcherEvent then
			cleanUp()
			return
		end
	end
end

function goEnv.quitDispatcher()
	coroutine.yield(quitDispatcherEvent)
end

function termCompatibility()
	term.redirect(nativeTerm)
	termCompatibilityMode = true
end

-- You should not manipulate the terminal with the
-- dispatch quit function, as it will be also be
-- called on dirty quits, such as errors.
function onDispatcherQuit(func)
	dispatcherQuitFunc = func
end

function quitDispatcher()
	os.queueEvent(quitDispatcherEvent)
end

goEnv.onDispatcherQuit = onDispatcherQuit
goEnv.termCompatibility = termCompatibility