

### gsnova 0.34 修改版，用于openshift v4 docker部署  


#### 修改项  

目标：规避openshift代理检查等  
基准：官方 0.34 commit  c6d0717 aug 11,2018  
- 版本号v34_200516_tls13 docker 200516  
- 使用go1.14.3编译服务端和客户端,upx服务端  
- Docker增加procps ca-certificates  
- CipherSuites tls ver 客户端和服务端安全约束  
- stat 增加显示uptime, req的Header等信息,客户端启动显示server stat  
- 升级外部依赖包至20200410,yinqiwen xtaci marten-seemann lucas-clemente名下包暂不升  
- 修改本机socks5时,关闭sni sniff,适配tls为ip的应用场景  
- 客户端增加json配置项ForceTls13，适用tls http2 wss，服务端自适应  
- 选项"tls13"或"tls12"，“auto"或不配置表示 tls12-tls13  
- 客户端增加json配置项SecureVerify，用于tls相关校验域名证书   
- json文件增加行尾注释，格式"\t//"或" //"  
- PAC规则加入blocked，例{"Host":["ad.12306.cn"],"Remote":"blocked"}，用于ban ad域名  
- 修改pid文件名  
- 随机证书改为2048bit(1024)  
- client模式时，显示心跳包延迟时间  
- 仅在client模式时显示ASCIILogo，server模式时不显示  
- server端默认Mux.IdleTimeout改为1800(300)，即使无数据也保持长连接30分钟  
- 加大默认的Mux的MaxStreamWindow和StreamMinRefresh为原有4倍，即2048k和128k  
- 增加key的环境变量AVONSG_CIPHER_KEY，也保留原有GSNOVA_CIPHER_KEY，规避检查需要  
- 增加环境变量AVONSG_CIPHER_USER，仅用于服务端鉴权，优先级高于命令行和json  
- remote.indexCallback改为http.StatusOK,即去掉原有https访问时的版本提示，规避检查需要  
- 加入logger.Printf，修改所有包log.Printf为logger调用  
- 修改logger包，加入none及null选项，便于server端命令行模式时，使用-log none关闭所有提示  
- 修正loadGFWList长时间不释放https连接  
- 修正AllowUsers鉴权失效问题
- 同步官方cd936c6,增加HibernateAfterSecs参数，客户端默认30分钟无数据时关闭muxSession    
- 同步官方95be5a5,服务端内置10秒无活动传输断开时间调整为90秒,客户端需要json中修改      
- ServerConf.Mux.StreamIdleTimeout = 90 // 10  
- ServerConf.Mux.SessionIdleTimeout = 1800 //300  
- 增加了server.json(docker未使用)，增加了client.json(仅参考)  
- 同步至官方cee73a4  
- 同步至官方8571b04  
- 同步至官方cba06fa  
- 同步至官方c6d0717 0.34 阶段版  

#### docker  
<https://hub.docker.com/r/devwak/avonsg_openshift>  


#### Thanks : yinqiwen  
<https://github.com/yinqiwen/gsnova>  
  
  
  

