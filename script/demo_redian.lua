-- 导航到百度首页
navigate("http://www.baidu.com")

-- 等待热点内容加载完成
wait_for_selector("xpath=//*[@id='s-hotsearch-wrapper']")

-- 定位热点内容的父级容器
hotsearch = find_element("xpath=//*[@id='s-hotsearch-wrapper']")
if hotsearch then
    -- 获取热点容器中包含的热点条目
    hot_items = find_elements("xpath=//*[@id='s-hotsearch-wrapper']//a")

    -- 遍历所有热点条目并打印标题和链接
    print("最新热点如下：")
    for _, item in ipairs(hot_items) do
        -- 获取热点标题
        title = get_text(item)
        -- 获取热点链接
        link = get_attribute(item, "href")
        print("热点标题：" .. title)
        print("热点链接：" .. link)
    end
else
    print("未找到热点内容")
end
