package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/cespare/xxhash"
)

func main() {
	flag.Parse()
	path := flag.Arg(0)
	if path == "" {
		fmt.Fprintln(os.Stderr, "Please provide a path")
		os.Exit(1)
	}
	var wg sync.WaitGroup
	err := filepath.Walk(path, func(p string, i os.FileInfo, err error) error {
		if err != nil {
			panic(err)
		}
		if !i.IsDir() {
			wg.Add(1)
			go func() {
				f, err := ioutil.ReadFile(p)
				if err != nil {
					panic(err)
				}
				h := xxhash.Sum64(f)
				os.Rename(p, filepath.Dir(p)+"/"+fmt.Sprintf("%016x", h)+filepath.Ext(p))
				wg.Done()
			}()
		}
		return nil
	})
	wg.Wait()
	if err != nil {
		panic(err)
	}
}
