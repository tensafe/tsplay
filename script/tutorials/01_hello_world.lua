print("lesson 01: hello world from lua")

greeting = set_var("greeting", "hello from lua")
payload = {
    lesson = "01",
    mode = "lua",
    greeting = greeting
}

write_json("artifacts/tutorials/01-hello-world-lua.json", payload)
print("wrote artifacts/tutorials/01-hello-world-lua.json")
