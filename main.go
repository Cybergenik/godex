package main

import (
    "os"
    "fmt"
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
        os.Exit(1)
    }
    var total int64
    c := make(chan int64)
    threads := 0
    for _, f := range files {
        if f.IsDir() {
            go traverse(filepath.Join(file, f.Name()), c)
            threads++
        } else {
            finfo, _ := f.Info()
            total += finfo.Size()
        }
    }
    for i:=0;i<threads;i++{
        total += <-c
    }
    recChan<-total
}

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
        var total_size int64
        if file.IsDir() {
            c := make(chan int64)
            go traverse(fname, c)
            total_size = <-c
        } else {
            total_size = file.Size()
        }
        if total_size < 1000.0 { // bytes
            fmt.Printf("%v size: %d\n", file.Name(), total_size)
        } else if total_size < 1_000_000 { // kilobyte
            fmt.Printf("%v size: %.1fK\n", file.Name(), float64(total_size)/1_000.0)
        } else if total_size < 1_000_000_000 { // megabyte
            fmt.Printf("%v size: %.1fMb\n", file.Name(), float64(total_size)/1_000_000.0)
        } else {
            fmt.Printf("%v size: %.1fGb\n", file.Name(), float64(total_size)/1_000_000_000.0)
        }
    }
}
