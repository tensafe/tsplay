--navigate("https://www.wps.cn/")
--wait_for_network_idle()
--download_file("//*[@id=\"content\"]/div[1]/div/div[1]/div[1]/div[4]/button[1]", "/Users/tensafe/code/tsplay/demo/wps.zip")
navigate("https://www.todesk.com/download.html")
wait_for_network_idle()
-- download_file("//*[@id=\"__layout\"]/div/div[1]/div[2]/div[2]/div[1]/a[2]", "./af.exe")
download_url("https://dl.todesk.com/irrigation/ToDesk_4.7.6.3.exe", "./demo.exe")
