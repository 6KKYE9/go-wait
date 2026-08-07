package main

import (
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestWaitFileAppears(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "ready.txt")
	go func() {
		time.Sleep(200 * time.Millisecond)
		os.WriteFile(p, []byte("ok"), 0644)
	}()
	if err := waitFile(p, 3*time.Second, 50*time.Millisecond); err != nil {
		t.Fatalf("应等到文件出现，实际 %v", err)
	}
}

func TestWaitFileTimeout(t *testing.T) {
	p := filepath.Join(t.TempDir(), "never.txt")
	if err := waitFile(p, 400*time.Millisecond, 100*time.Millisecond); err == nil {
		t.Fatalf("文件不存在应该超时报错")
	}
}

func TestWaitPortConnects(t *testing.T) {
	// 起一个本地监听，等它可连
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()
	go func() {
		for {
			c, e := ln.Accept()
			if e != nil {
				return
			}
			c.Close()
		}
	}()
	addr := ln.Addr().String()
	if err := waitPort(addr, 3*time.Second, 50*time.Millisecond); err != nil {
		t.Fatalf("应等到端口可连，实际 %v", err)
	}
}

func TestWaitPortTimeout(t *testing.T) {
	// 连一个没人监听的端口，应当超时
	if err := waitPort("127.0.0.1:1", 400*time.Millisecond, 100*time.Millisecond); err == nil {
		t.Fatalf("端口不应可连，应该超时")
	}
}

func TestWaitCmdSuccess(t *testing.T) {
	// exit 0 永远成功（Windows cmd 和 sh 都支持 exit）
	if err := waitCmd("exit 0", 2*time.Second, 50*time.Millisecond); err != nil {
		t.Fatalf("exit 0 应立刻成功，实际 %v", err)
	}
}

func TestWaitCmdTimeout(t *testing.T) {
	// exit 1 永远失败，应当超时
	if err := waitCmd("exit 1", 400*time.Millisecond, 100*time.Millisecond); err == nil {
		t.Fatalf("exit 1 应超时")
	}
}
