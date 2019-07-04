#!/usr/bin/env lua
function request(host, port, args)
  local http=require("socket.http")

  local request_body = {}
  if args then
      request_body = args
  end
  response_body = {}
  local res, header, code = http.request(string.format("http://%s",host))

  -- if type(response_headers) == "table" then
  --   for k, v in pairs(response_headers) do
  --     print(k, v)
  --   end
  -- end

  -- print("Response body:")
  -- if type(response_body) == "table" then
  --   print(table.concat(response_body))
  -- else
  --   print("Not a table:", type(response_body))
  -- end
  
  local respLen = string.len(res)
  -- print(respLen, code)
  return code , respLen

end

-- print(request("www.baidu.com"))
