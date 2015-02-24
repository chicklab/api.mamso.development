package main

import (
	"flag"
	"fmt"
	// "github.com/msgpack/msgpack-go"
	// "log"
	"io/ioutil"
	"os"
	"strings"
)

var (
	HOME   = os.Getenv("HOME")
	USER   = os.Getenv("USER")
	GOROOT = os.Getenv("GOROOT")
)

func init() {
	// fmt.Println("aaaaaaaaa")
	// if USER == "" {
	// 	log.Fatal("$USER not set")
	// }
	// if HOME == "" {
	// 	HOME = "/usr/" + USER
	// }
	if GOROOT == "" {
		GOROOT = HOME + "/go"
	}
	// GOROOTはコマンドラインから--gorrotフラグを指定することで上書き可能
	flag.StringVar(&GOROOT, "goroot", GOROOT, "Go root directory")
}

func main() {
	list, err := ioutil.ReadDir("/opt/hoge/")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}

	for _, finfo := range list {
		if finfo.IsDir() || -1 == strings.Index(finfo.Name(), ".txt") {
			continue
		}
		fmt.Printf("%q\n", finfo.Name())
	}

	f := flag.Int("flag1", 0, "flag 1")
	flag.Parse()

	fmt.Println("flag1: ", os.Stderr)
	fmt.Println("os.Args: ", os.Args)
	fmt.Println("compare")
	fmt.Println("flag.Args: ", flag.Args())
	if *f == 100 {
		fmt.Println("Hello")
	}

	x := []int{1, 2, 3}
	y := []int{4, 5, 6}
	x = append(x, y...)
	fmt.Println(x)

	fmt.Println(HOME)
	fmt.Println(USER)
	fmt.Println(GOROOT)

	// a := Sum([...]float64{1.2, 12.1, 3.5})
	// array := [...]float64{7.0, 10, 9.1}
	z := Sum(&[...]float64{7.0, 10, 9.1}) // 注：明示的にアドレス演算子を使用
	fmt.Println(z)
}

func Sum(a *[3]float64) (sum float64) {
	for _, v := range *a {
		sum += v
	}
	return
}
