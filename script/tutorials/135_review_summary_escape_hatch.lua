print("lesson 135: lua escape hatch sample")

payload = {
    lesson = "135",
    mode = "lua",
    recommendation = "extract_to_flow",
    reason = "simple_orchestration_only",
}

write_json("artifacts/tutorials/135/review-summary-from-lua.json", payload)
print("wrote artifacts/tutorials/135/review-summary-from-lua.json")
