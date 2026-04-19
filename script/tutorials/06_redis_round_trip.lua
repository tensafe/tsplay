redis_del("tutorial:greeting")
redis_del("tutorial:payload")
redis_del("tutorial:counter")

redis_set("tutorial:greeting", "hello from lua", 600)
greeting = redis_get("tutorial:greeting")

redis_set("tutorial:payload", {
    lesson = "06",
    mode = "lua",
    status = "ready"
}, 600)
payload_text = redis_get("tutorial:payload")
payload_status = json_extract(payload_text, "$.status")

counter_one = redis_incr("tutorial:counter")
counter_three = redis_incr("tutorial:counter", 2)

write_json("artifacts/tutorials/06-redis-round-trip-lua.json", {
    lesson = "06",
    mode = "lua",
    greeting = greeting,
    payload_status = payload_status,
    counter_one = counter_one,
    counter_three = counter_three
})

print("greeting:", tostring(greeting))
print("payload status:", tostring(payload_status))
print("counter after +1:", tostring(counter_one))
print("counter after +2:", tostring(counter_three))
print("wrote artifacts/tutorials/06-redis-round-trip-lua.json")
