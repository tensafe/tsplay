page_url = os.getenv("TSPLAY_WAIT_URL") or "http://127.0.0.1:8000/demo/retry_wait_until.html"

navigate(page_url)
wait_for_selector("#start-job", 5000)
click("#start-job")

job_ready = false
poll_attempts = 0

for attempt = 1, 25 do
    poll_attempts = attempt
    if is_visible("#job-ready") then
        job_ready = true
        break
    end
    sleep(0.2)
end

if not job_ready then
    error("expected async job to become ready within timeout")
end

job_status = extract_text("#job-status", 5000)

write_json("artifacts/tutorials/17-wait-until-ready-lua.json", {
    lesson = "17",
    mode = "lua",
    page_url = page_url,
    poll_attempts = poll_attempts,
    job_ready = job_ready,
    job_status = job_status
})

print("job status:", tostring(job_status))
print("wrote artifacts/tutorials/17-wait-until-ready-lua.json")
print("press Ctrl+C when you finish inspecting the browser")
