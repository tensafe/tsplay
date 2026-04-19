demo_url = os.getenv("TSPLAY_DEMO_URL") or "http://127.0.0.1:8000/demo/demo.html"
target_value = os.getenv("TSPLAY_OPTION_VALUE") or "5"

navigate(demo_url)
wait_for_selector("#options", 5000)
select_option("#options", target_value)

selected_ok = is_selected("#options option[value='" .. target_value .. "']")
if not selected_ok then
    error("expected option " .. target_value .. " to be selected")
end

selected_label = extract_text("#options option:checked", 5000)

write_json("artifacts/tutorials/12-custom-json-output-lua.json", {
    lesson = "12",
    mode = "lua",
    source = {
        page = "local-demo",
        demo_url = demo_url
    },
    selection = {
        value = target_value,
        label = selected_label,
        verified = selected_ok
    },
    summary = "Selected " .. selected_label
})

print("selected label:", tostring(selected_label))
print("wrote artifacts/tutorials/12-custom-json-output-lua.json")
print("press Ctrl+C when you finish inspecting the browser")
