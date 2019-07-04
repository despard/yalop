# yalop 通用压测程序框架

## 用法

* 编写客户端类实现如下接口

```
type IClient interface {
	DoRequest() (error, int, int64)
	DoClose() error
}
```

* 生成压测工具类，传入客户端类。、

```
w := &requester.Work{
		Client:   $IClient,
		N:        num,
		C:        conc,
		QPS:      q,
		Timeout:  *t,
		Output:   *output,
		Duration: dur,
	}
```

* 运行

```
w.Init()
w.Run()
```
## yalop压测工具
```
Usage: press [options...] <url>

Options:
  -n  Number of requests to run. Default is 200.
  -c  Number of requests to run concurrently. Total number of requests cannot
      be smaller than the concurrency level. Default is 50.
  -q  Rate limit, in queries per second (QPS). Default is no rate limit.
  -z  Duration of application to send requests. When duration is reached,
      application stops and exits. If duration is specified, n is ignored.
      Examples: -z 10s -z 3m.
  -o  Output type. If none provided, a summary is printed.
      "csv" is the only supported alternative. Dumps the response
      metrics in comma-separated values format.

  -t  Timeout for each request in seconds. Default is 20, use 0 for infinite.

  -cpus                 Number of used cpu cores.
                        (default for current machine is 4 cores)
  -script  Lua script file path. Load lua script for test.
```

* 默认http请求
```
./yalop -n 10 -c 10 www.baidu.com:80
```
* 加载lua脚本，脚本实现请求逻辑
```
./yalop -n 10 -c 10 -script ./script/callHttp.lua www.baidu.com:80
```
