page_url = os.getenv("TSPLAY_UPLOAD_URL") or "http://127.0.0.1:8000/demo/upload.html"
input_file = os.getenv("TSPLAY_UPLOAD_FILE") or "demo/data/upload_receipt.pdf"
expected_filename = os.getenv("TSPLAY_UPLOAD_FILENAME") or "upload_receipt.pdf"

navigate(page_url)
wait_for_selector("#fileInput", 5000)
upload_file("#fileInput", input_file)
assert_text("#fileInfo", expected_filename, 5000)

file_info = extract_text("#fileInfo", 5000)

write_json("artifacts/tutorials/18-upload-single-file-lua.json", {
    lesson = "18",
    mode = "lua",
    page_url = page_url,
    input_file = input_file,
    expected_filename = expected_filename,
    file_info = file_info
})

print("uploaded:", expected_filename)
print("wrote artifacts/tutorials/18-upload-single-file-lua.json")
print("press Ctrl+C when you finish inspecting the browser")
