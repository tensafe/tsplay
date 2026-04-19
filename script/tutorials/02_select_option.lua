demo_url = os.getenv("TSPLAY_DEMO_URL") or "http://127.0.0.1:8000/demo/demo.html"

navigate(demo_url)
wait_for_selector("#options", 5000)
select_option("#options", "5")

option_5_selected = is_selected("#options option[value='5']")
if not option_5_selected then
    error("expected option 5 to be selected")
end

write_json("artifacts/tutorials/02-select-option-lua.json", {
    lesson = "02",
    mode = "lua",
    demo_url = demo_url,
    selected = option_5_selected
})

print("selected option 5:", tostring(option_5_selected))
print("wrote artifacts/tutorials/02-select-option-lua.json")
print("press Ctrl+C when you finish inspecting the browser")
