page_url = os.getenv("TSPLAY_MULTI_UPLOAD_URL") or "http://127.0.0.1:8000/demo/multi_upfile.html"
file_a = os.getenv("TSPLAY_UPLOAD_FILE_A") or "demo/data/upload_receipt.pdf"
file_b = os.getenv("TSPLAY_UPLOAD_FILE_B") or "demo/data/upload_avatar.png"
expected_a = os.getenv("TSPLAY_UPLOAD_FILENAME_A") or "upload_receipt.pdf"
expected_b = os.getenv("TSPLAY_UPLOAD_FILENAME_B") or "upload_avatar.png"

navigate(page_url)
wait_for_selector("#fileInput", 5000)
upload_multiple_files("#fileInput", file_a, file_b)
assert_text("#fileInfo", expected_a, 5000)
assert_text("#fileInfo", expected_b, 5000)

file_info = extract_text("#fileInfo", 5000)

write_json("artifacts/tutorials/19-upload-multiple-files-lua.json", {
    lesson = "19",
    mode = "lua",
    page_url = page_url,
    files = {file_a, file_b},
    expected_files = {expected_a, expected_b},
    file_info = file_info
})

print("uploaded:", expected_a .. ", " .. expected_b)
print("wrote artifacts/tutorials/19-upload-multiple-files-lua.json")
print("press Ctrl+C when you finish inspecting the browser")
