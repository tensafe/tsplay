navigate("https://www.wps.cn/")
wait_for_network_idle()
download_file("//*[@id=\"content\"]/div[1]/div/div[1]/div[1]/div[4]/button[1]", "/Users/tensafe/code/tsplay/demo/wps.zip")
