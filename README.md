本项目为webrtc的signal的服务器和js实现


1. 安装启动turnserver
启动示例
nohup turnserver -o -a -f -v --mobility -m 10 --max-bps=100000 --min-port=32355 --max-port=65535 --user=ling:ling1234 --user=ling2:ling1234 -r demo &
修改home.html configuration参数

2. make
    go get github.com/gorilla/websocket
    go run *.go

3. firefox打开http://ip:8080/, 其他浏览器未测试

