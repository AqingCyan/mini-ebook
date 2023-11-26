-- 发送到的 key，也就是是 code:业务:手机号码
local key = KEYS[1]
-- 使用次数，也就是验证次数
local cntKey = key..":cnt"
-- 用户输入的验证码
local expectedCode = ARGV[1]

local cnt = tonumber(redis.call("get", cntKey))
local code = redis.call("get", key)

if cnt == nil or cnt <= 0 then
    -- 验证次数已经用完
    return -1
end

if code == expectedCode then
    -- 验证码正确，验证次数就设置为 0
    redis.call("set", cntKey, 0)
    return 0
else
    -- 验证码错误，可能是用户输入错误，验证次数减一
    redis.call("decr", cntKey)
    return -2
end