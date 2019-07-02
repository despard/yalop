# yalop 通用压测程序框架

###用法

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
