-- 定义一个回调函数，用于处理拦截的请求
local function handle_request(url, method, resource_type)
    print("Intercepted request:")
    print("URL: " .. url)
    print("Method: " .. method)
    print("Resource Type: " .. resource_type)

    -- 假设我们想将所有请求重定向到一个特定的 URL
    if method == "GET" and resource_type == "document" then
        return "https://www.163.com"
    end

    -- 如果没有返回新的 URL，请求将继续原样
    print("not deal with....")
    return nil
end

-- 调用 Go 函数设置请求拦截器
intercept_request(handle_request, "**/*.js")

-- 继续执行其他脚本逻辑
print("Request interceptor has been set up.")
