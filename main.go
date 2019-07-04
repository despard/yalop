package main

import (
	"errors"
	"flag"
	"fmt"
	"math"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/BixData/gluasocket"
	"github.com/yalop/requester"
	lua "github.com/yuin/gopher-lua"
)

const (
	heyUA = "yalop.1.0.0"
)

var (
	s      = flag.String("script", "", "")
	output = flag.String("o", "", "")
	c      = flag.Int("c", 50, "")
	n      = flag.Int("n", 400, "")
	q      = flag.Float64("q", 0, "")
	t      = flag.Int("t", 20, "")
	z      = flag.Duration("z", 0, "")

	cpus = flag.Int("cpus", runtime.GOMAXPROCS(-1), "")
)

var usage = `Usage: press [options...] <url>

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
                        (default for current machine is %d cores)
  -script  Lua script file path. Load lua script for test.
`

type HttpClient struct {
	domain string
	port   int
	uri    string
}

func (hc *HttpClient) DoRequest() (error, int, int64) {
	if hc.domain == "" {
		return errors.New("domain is empty."), 20, 0
	}
	if hc.port == 0 {
		hc.port = 80
	}
	url := fmt.Sprintf("http://%s:%d/%s", hc.domain, hc.port, hc.uri)
	if resp, err := http.Get(url); err != nil {
		return err, -1, 0
	} else {
		return err, resp.StatusCode, resp.ContentLength
	}
}

func (uc *HttpClient) DoClose() error {
	return nil
}

type LuaClient struct {
	script string
	host   string
	port   string
	args   string
}

//每次都会启动一个lua解释器
func (lc *LuaClient) DoRequest() (error, int, int64) {
	L := lua.NewState()
	// Preload LuaSocket modules
	gluasocket.Preload(L)
	defer L.Close()
	// 加载lua脚本
	if err := L.DoFile(lc.script); err != nil {
		panic(err)
	}
	// 调用request函数
	err := L.CallByParam(lua.P{
		Fn:      L.GetGlobal("request"), // 获取fib函数引用
		NRet:    2,                      // 指定返回值数量
		Protect: true,                   // 如果出现异常，是panic还是返回err
	}, lua.LString(lc.host), lua.LString(lc.port), lua.LString(lc.args)) // 传递输入参数
	if err != nil {
		panic(err)
	}
	// 获取返回结果
	code := L.ToInt(-2)
	size := L.ToInt64(-1)
	// 从堆栈中扔掉返回结果
	L.Pop(1)
	L.Pop(1)
	return nil, code, size
}

func (lc *LuaClient) DoClose() error {
	return nil
}
func main() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, fmt.Sprintf(usage, runtime.NumCPU()))
	}

	flag.Parse()
	if flag.NArg() < 1 {
		usageAndExit("input addr.")
	}

	runtime.GOMAXPROCS(*cpus)
	num := *n
	conc := *c
	q := *q
	dur := *z
	script := *s

	if dur > 0 {
		num = math.MaxInt32
		if conc <= 0 {
			usageAndExit("-c cannot be smaller than 1.")
		}
	} else {
		if num <= 0 || conc <= 0 {
			usageAndExit("-n and -c cannot be smaller than 1.")
		}

		if num < conc {
			usageAndExit("-n cannot be less than -c.")
		}
	}
	saddr := strings.Split(flag.Args()[0], ":")[0]
	port  := strings.Split(flag.Args()[0], ":")[1]
	args := ""
	if len(flag.Args()) > 1 {
		args = flag.Args()[1]
	}
	var w *requester.Work
	if script != "" {
		lc := &LuaClient{
			script: script,
			host:   saddr,
			port:   port,
			args:   args,
		}
		w = &requester.Work{
			Client:   lc,
			N:        num,
			C:        conc,
			QPS:      q,
			Timeout:  *t,
			Output:   *output,
			Duration: dur,
		}
	} else {
		iPort, _ := strconv.Atoi(port)
		uc := &HttpClient{
			domain: saddr,
			port:   iPort,
			uri:    args,
		}
		w = &requester.Work{
			Client:   uc,
			N:        num,
			C:        conc,
			QPS:      q,
			Timeout:  *t,
			Output:   *output,
			Duration: dur,
		}
	}

	w.Init()
	w.Run()
}

func errAndExit(msg string) {
	fmt.Fprintf(os.Stderr, msg)
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(1)
}

func usageAndExit(msg string) {
	if msg != "" {
		fmt.Fprintf(os.Stderr, msg)
		fmt.Fprintf(os.Stderr, "\n\n")
	}
	flag.Usage()
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(1)
}

func parseInputWithRegexp(input, regx string) ([]string, error) {
	re := regexp.MustCompile(regx)
	matches := re.FindStringSubmatch(input)
	if len(matches) < 1 {
		return nil, fmt.Errorf("could not parse the provided input; input = %v", input)
	}
	return matches, nil
}
