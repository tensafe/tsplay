navigate("http://localhost:63342/tsplay/demo/demo.html?_ijt=seneplhjrg1m5kt2m981cricto&_ij_reload=RELOAD_ON_SAVE")
wait_for_network_idle()

aa = is_selected("#options option[value='5']")

print(aa)
