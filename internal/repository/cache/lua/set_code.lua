-- 发送到的 key，也就是是 code:业务:手机号码
local key = KEYS[1]
-- 使用次数，也就是验证次数
local cntKey = key..":cnt"
-- 你准备存储的验证码
local val = ARGV[1]
-- 验证码的有效时间是十分钟，600秒
local ttl = tonumber(redis.call("ttl", key))

-- 在 Redis 中，ttl 为 -2 表示 key 不存在，-1 表示 key 存在，但是没有过期时间
-- 其余情况，ttl 为剩余的过期时间

-- 在下面的逻辑中，如果  ttl 为 -1，表示 key 存在，但是没有过期时间，因此不能发送验证码，这是一个系统错误，可能是某些对 Redis 不规范的使用造成的
-- 如果 ttl 为 -2，表示 key 不存在，可以发送验证码
-- 或者 ttl 剩余时间小于 540 秒，也可以发送验证码，这意味着上一次发送的验证码已经超过 1 分钟了

if ttl == -1 then
    -- key 存在，但是没有过期时间
    return -2
elseif ttl == -2 or ttl < 540 then
    -- 可以发送验证码
    redis.call("set", key, val)
    -- 设置过期时间为 600 秒，这是一条验证码的有效期，超过这个时间，验证码就失效了
    redis.call("expire", key, 600)
    -- 设置验证次数为 3，这是一条验证码允许的最大被验证次数，超过这个次数，验证码就失效了
    redis.call("set", cntKey, 3)
    redis.call("expire", cntKey, 600)
    return 0
else
    -- 发送太频繁
    return -1
end
