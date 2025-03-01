function printTable(t, indent)
    indent = indent or "" -- 缩进，用于格式化输出
    for k, v in pairs(t) do
        if type(v) == "table" then
            print(indent .. k .. ":")
            printTable(v, indent .. "  ") -- 递归打印嵌套表
        else
            print(indent .. k .. ": " .. tostring(v))
        end
    end
end

print("hello!")
navigate("https://www.baidu.com")
wait_for_network_idle()
type_text("#kw", "山东大学")
click("#su")
wait_for_network_idle()
wait_for_selector("xpath=//div[@class='FYB_RD']")

--data = get_html("xpath=//div[@class='FYB_RD']")
--print(data)

links = get_all_links("xpath=//div[@class='FYB_RD']")
printTable(links, 2)

