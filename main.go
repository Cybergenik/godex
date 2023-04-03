package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func usage() {
	fmt.Println(
		`Godex (Go-indexer): Simple multi-threaded Filesystem traverser
USAGE: 
    godex <path to dir/file>
    Ex: godex /home/user1
        godex /home/user1 /home/user2
OPTIONS:
    --help: Prints this message`)
}

func traverse(file string, recChan chan int64) {
	files, err := os.ReadDir(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening \"%v\": %v\n", file, err)
        recChan <- 0
        return
	}
    var total int64
	c := make(chan int64)
	threads := 0
	for _, f := range files {
        t := f.Type()
        if t.IsDir() || t.IsRegular() {
            finfo, err := f.Info()
            if err != nil {
                fmt.Fprintf(os.Stderr, "Error opening \"%v\": %v\n", f, err)
                recChan <- 0
                return
            }
            total += finfo.Size()
            if t.IsDir() {
                threads++
                go traverse(filepath.Join(file, f.Name()), c)
            }
        }
	}
	for i := 0; i < threads; i++ {
		total += <-c
	}
	recChan <- total
}

const (
	KILOBYTE = 1024
	MEGABYTE = 1_048_576
	GIGABYTE = 1_073_741_824
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Godex: No arg provided, for usage do: --help")
		os.Exit(1)
	}
	if os.Args[1] == "--help" || os.Args[1] == "-help" {
		usage()
		os.Exit(0)
	}
	for _, fname := range os.Args[1:] {
		file, err := os.Lstat(fname)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening \"%v\": %v\n", fname, err)
			os.Exit(1)
		}
        total_size := file.Size()
		if file.IsDir() {
			c := make(chan int64)
			go traverse(fname, c)
			total_size += <-c
		}
		if total_size < KILOBYTE { // bytes
			fmt.Printf("%v size: %d\n", file.Name(), total_size)
		} else if total_size < MEGABYTE { // kilobyte
			fmt.Printf("%v size: %.1fK\n", file.Name(), float64(total_size)/KILOBYTE)
		} else if total_size < GIGABYTE { // megabyte
			fmt.Printf("%v size: %.1fMb\n", file.Name(), float64(total_size)/MEGABYTE)
		} else {
			fmt.Printf("%v size: %.1fGb\n", file.Name(), float64(total_size)/GIGABYTE)
		}
	}
}
