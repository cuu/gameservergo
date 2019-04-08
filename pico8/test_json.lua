json = require("json")


function safe_format(funcname,...)
  local ret = {Func=funcname}
 	local args = {}
 
  for i, v in ipairs{...} do
    local arg = {Type="",Value=""}
    if type(v) == "string" then
      arg.Type="S"
    end
    if type(v) == "number" then
      arg.Type="I"
    end
		if type(v) == "boolean" then
			arg.Type="B"
		end	
    arg.Value=v

    table.insert(args,arg)
  end
 	
	if #args > 0 then
		ret.Args = args
 	end

  return json.encode(ret)
end

print( safe_format("res.print","test",10,math.floor(20),true,1,2,3,"order","you"))
