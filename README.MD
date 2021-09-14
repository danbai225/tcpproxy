# 简介
一个tcp通道 加密复用的程序

服务端运行
```shell
./socks5proxy -s -port 18889 -ps 123456 -addr=127.0.0.1:7891
```
客户端
```shell
./socks5proxy -port 18889 -ps 123456 -addr serverhost:18889
```