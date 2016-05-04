本项目为webrtc的signal的服务器和js实现  
项目实现了聊天室和1对1视频的功能  
示例地址 http://45.32.245.83:8080/  

1. 安装启动turnserver  
1.1 项目地址  
https://github.com/coturn/coturn  
1.2 启动示例  
 turnserver -o -a -f -v --mobility -m 10 --max-bps=100000 --min-port=32355 --max-port=65535 --user=ling:ling1234 --user=ling2:ling1234 -r demo    
1.3 修改home.html   
configuration参数为自己配置的turnserver  

2. make  
    go get github.com/gorilla/websocket  
    go run *.go  

3. firefox打开http://ip:8080/, 其他浏览器暂不兼容  

