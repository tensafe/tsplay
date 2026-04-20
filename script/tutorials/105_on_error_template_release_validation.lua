page_url = os.getenv("TSPLAY_TEMPLATE_RELEASE_URL") or "http://127.0.0.1:8000/demo/template_release_lab.html"
invalid_ticket = os.getenv("TSPLAY_INVALID_RELEASE_TICKET") or "BAD-TICKET"

navigate(page_url)
wait_for_selector("#release-ticket", 5000)

type_text("#release-ticket", invalid_ticket)
click("#finalize-release-button")

recovery_path = "unexpected_success"
failure_status = ""
error_message = ""

ok, err = pcall(function()
    status_text = get_text("#finalize-status")
    if string.find(status_text or "", "Release finalized", 1, true) == nil then
        error("expected finalize status to contain 'Release finalized', actual=" .. tostring(status_text))
    end
end)

if not ok then
    recovery_path = "handled_error"
    failure_status = get_text("#finalize-status")
    error_message = tostring(err)
    click("#reset-finalize")
end

final_status = get_text("#finalize-status")

write_json("artifacts/tutorials/105-on-error-template-release-validation-lua.json", {
    lesson = "105",
    mode = "lua",
    page_url = page_url,
    invalid_ticket = invalid_ticket,
    recovery_path = recovery_path,
    failure_status = failure_status,
    final_status = final_status,
    error_message = error_message
})

print("recovery path:", tostring(recovery_path))
print("wrote artifacts/tutorials/105-on-error-template-release-validation-lua.json")
