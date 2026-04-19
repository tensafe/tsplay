demo_url = os.getenv("TSPLAY_DEMO_URL") or "http://127.0.0.1:8000/demo/demo.html"
target_value = os.getenv("TSPLAY_OPTION_VALUE") or "7"

navigate(demo_url)
wait_for_selector("#options", 5000)
select_option("#options", target_value)

selected_ok = is_selected("#options option[value='" .. target_value .. "']")
if not selected_ok then
    error("expected option " .. target_value .. " to be selected")
end

selected_label = extract_text("#options option:checked", 5000)

write_json("artifacts/tutorials/11-select-another-option-lua.json", {
    lesson = "11",
    mode = "lua",
    demo_url = demo_url,
    target_value = target_value,
    selected_ok = selected_ok,
    selected_label = selected_label
})

print("selected value:", tostring(target_value))
print("selected label:", tostring(selected_label))
print("wrote artifacts/tutorials/11-select-another-option-lua.json")
print("press Ctrl+C when you finish inspecting the browser")
