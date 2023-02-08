local ans = {}
local videoIds = redis.call('lrange', KEYS[1], ARGV[1], ARGV[2])
for i, id in ipairs(videoIds) do
    local m = {}
    m[1] = redis.call('get', KEYS[2] .. id)
    m[2] = redis.call('get', KEYS[3] .. id)
    m[3] = redis.call('llen', KEYS[4] .. id)
    m[4] = redis.call('scard', KEYS[5] .. id)
    m[5] = redis.call('sismember', KEYS[5] .. id, ARGV[3])
    ans[i] = m
end
return ans