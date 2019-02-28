package main

import (
    "os"
    "fmt"
    "bytes"
    "encoding/binary"
)

func dword2uint32(bs []byte) uint32 {
    var n uint32
    binary.Read(bytes.NewReader(bs), binary.BigEndian, &n)
    return n
}

func main() {

    if len(os.Args) != 2 {
        fmt.Fprintf(os.Stderr, "Usage: eat-cookie <FILE.binarycookies>")
        os.Exit(0)
    }

    file, err := os.Open(os.Args[1])
    if err != nil {
        fmt.Fprintf(os.Stderr, "%s", err)
        os.Exit(1)
    }
    defer file.Close()

    // Check whether the magic number is "cook"
    magic_number := make([]byte, 4)
    file.Read(magic_number)
    if string(magic_number) != "cook" {
        fmt.Fprint(os.Stderr, "The magic number has unexpected", )
        os.Exit(1)
    }

    // Read number of pages 
    num_pages_buf := make([]byte, 4)
    file.Read(num_pages_buf)
    num_pages := dword2uint32(num_pages_buf)

    // Read the each page size
    page_sizes := []uint32{}
    for i := uint32(0); i < num_pages; i++ {
        page_size_buf := make([]byte, 4)
        file.Read(page_size_buf)
        page_sizes = append(page_sizes, dword2uint32(page_size_buf))
    }
    
    pages := []string{}
    for _, v := range page_sizes {
        page_buf := make([]byte, v)
        file.Read(page_buf)
        pages = append(pages, string(page_buf))
        fmt.Println(pages)
    }
}
