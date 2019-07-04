--Lua--十进制转二进制
function dec_to_binary (data)
    local dst = ""
    local remainder, quotient

    --异常处理
    if not data then return dst end                 --源数据为空
    if not tonumber(data) then return dst end       --源数据无法转换为数字

    --如果源数据是字符串转换为数字
    if "string" == type(data) then
        data = tonumber(data)
    end

    while true do
        quotient = math.floor(data / 2)
        remainder = data % 2
        dst = dst..remainder
        data = quotient
        if 0 == quotient then
            break
        end
    end

    --翻转
    dst = string.reverse(dst)

    --补齐8位
    if 8 > #dst then
        for i = 1, 8 - #dst, 1 do
            dst = '0'..dst
        end
    end
    return dst
end

--Lua--二进制转十进制
function binary_to_dec (data)
    local dst = 0
    local tmp = 0

    --异常处理
    if not data then return dst end                 --源数据为空
    if not tonumber(data) then return dst end       --源数据无法转换为数字

    --如果源数据是字符串去除前面多余的0
    if "string" == type(data) then
        data = tostring(tonumber(data))
    end

    --如果源数据是数字转换为字符串
    if "number" == type(data) then
        data = tostring(data)
    end

    --转换
    for i = #data, 1, -1 do
        tmp = tonumber(data:sub(-i, -i))
        if 0 ~= tmp then
            for j = 1, i - 1, 1 do
                tmp = 2 * tmp
            end
        end
        dst = dst + tmp
    end
    return dst
end

function base64decode(data)
    local basecode = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
    local dst = ""
    local code = ""
    local tmp, index

    --异常处理
    if not data then return dst end                   --源数据为空

    data = data:gsub("\n", "")                        --去除换行符
    data = data:gsub("=", "")                         --去除'='

    for i = 1, #data, 1 do
        tmp = data:sub(i, i)
        index = basecode:find(tmp)
        if nil == index then
            return dst
        end
        index = index - 1
        tmp = dec_to_binary(index)
        code = code..tmp:sub(3)                       --去除前面多余的两个'00'
    end

    --开始解码
    for i = 1, #code, 8 do
        tmp = string.char(binary_to_dec(code:sub(i, i + 7)))
        if nil ~= tmp then
            dst = dst..tmp
        end
    end
    return dst
end

function request (host, port)
    local socket = require("socket.core")
    local udp = socket.udp()
    local lhost = host or '127.0.0.1'
    local lport = port or '10020'
    local sendPkg = base64decode("ATQVmQAAAe4AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAATiEAAATQAB7PLgEBAAAAAAEAAAAAAAAAAF0dTyYAAABk")
    local recvSize = 0
    udp:settimeout(5)
    function rec_msg()
        local recvmsg = udp:receive()
        if(recvmsg) then
            -- print('recudp data:'..recvmsg)
            return #recvmsg
        else
            -- print('recudp data nil')
        end
    end

    while 1 do
        udp:setpeername(lhost, lport)
        local udpsend = udp:send(sendPkg)
        if(udpsend) then
            -- print('udpsend ok')
            recvSize = rec_msg()
            break
        else
            -- print('udpsend err')
        end
    end
    udp:close()
    return 0, recvSize
end

print(request())
