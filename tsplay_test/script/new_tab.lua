navigate("http://www.baidu.com")
wait_for_network_idle()

new_tab("http://www.163.com")

wait_for_network_idle()

switch_to_tab(0)

close_tab()
