// go-wait 等某个条件满足再退出，常写在脚本或 CI 里。
// 支持等端口可连、等文件出现、等一条命令跑成功，可设超时和轮询间隔。
package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"time"
)

func waitPort(addr string, timeout, interval time.Duration) error {
	deadline := time.Now().Add(timeout)
	for {
		conn, err := net.DialTimeout("tcp", addr, interval)
		if err == nil {
			conn.Close()
			return nil
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("等待 %s 超时：%v", addr, err)
		}
		time.Sleep(interval)
	}
}

func waitFile(path string, timeout, interval time.Duration) error {
	deadline := time.Now().Add(timeout)
	for {
		if _, err := os.Stat(path); err == nil {
			return nil
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("等待文件 %s 超时", path)
		}
		time.Sleep(interval)
	}
}

func waitCmd(command string, timeout, interval time.Duration) error {
	deadline := time.Now().Add(timeout)
	for {
		// Windows 用 cmd /C，其它系统用 sh -c，这样命令里能写管道
		var cmd *exec.Cmd
		if runtime.GOOS == "windows" {
			cmd = exec.Command("cmd.exe", "/C", command)
		} else {
			cmd = exec.Command("sh", "-c", command)
		}
		if err := cmd.Run(); err == nil {
			return nil
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("等待命令成功超时：%s", command)
		}
		time.Sleep(interval)
	}
}

func main() {
	port := flag.String("port", "", "等待 host:port 可连接，例如 127.0.0.1:8080")
	file := flag.String("file", "", "等待某个文件出现")
	cmd := flag.String("cmd", "", "等待某条命令（sh -c）成功退出")
	timeout := flag.Duration("timeout", 30*time.Second, "最长等待时间，例如 60s")
	interval := flag.Duration("interval", 500*time.Millisecond, "轮询间隔，例如 1s")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "用法: go-wait [-port 地址] [-file 路径] [-cmd 命令] [-timeout 时长] [-interval 间隔]")
		fmt.Fprintln(os.Stderr, "  三种条件选一种，条件满足即退出码 0，超时退出码 1")
	}
	flag.Parse()

	chosen := 0
	for _, v := range []string{*port, *file, *cmd} {
		if v != "" {
			chosen++
		}
	}
	if chosen != 1 {
		fmt.Fprintln(os.Stderr, "必须且只能指定 -port / -file / -cmd 中的一个")
		flag.Usage()
		os.Exit(2)
	}

	var err error
	switch {
	case *port != "":
		fmt.Printf("等待 %s 可连接...\n", *port)
		err = waitPort(*port, *timeout, *interval)
	case *file != "":
		fmt.Printf("等待文件 %s 出现...\n", *file)
		err = waitFile(*file, *timeout, *interval)
	case *cmd != "":
		fmt.Printf("等待命令成功: %s\n", *cmd)
		err = waitCmd(*cmd, *timeout, *interval)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, "失败:", err)
		os.Exit(1)
	}
	fmt.Println("条件已满足")
}
