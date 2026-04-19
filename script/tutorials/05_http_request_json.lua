api_url = os.getenv("TSPLAY_HTTP_JSON_URL") or "http://127.0.0.1:8000/demo/data/order_summary.json"

api_result = http_request({
    url = api_url,
    response_as = "json"
})

http_status = json_extract(api_result, "$.status")
open_count = json_extract(api_result, "$.body.summary.open")
first_order_id = json_extract(api_result, "$.body.orders[0].id")

write_json("artifacts/tutorials/05-http-request-json-lua.json", {
    lesson = "05",
    mode = "lua",
    api_url = api_url,
    http_status = http_status,
    open_count = open_count,
    first_order_id = first_order_id
})

print("http status:", tostring(http_status))
print("open orders:", tostring(open_count))
print("first order id:", tostring(first_order_id))
print("wrote artifacts/tutorials/05-http-request-json-lua.json")
