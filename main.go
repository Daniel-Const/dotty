package main

import "fmt"

/*
 * Main CLI App
 *
 * Prompt to select and deploy a profile
 */

func main() {
    p := NewProfile("daniel-pc")
    p.Load()
    for i := range p.dots {
        fmt.Printf("File: %s, Destination: %s\n", p.dots[i].path, p.dots[i].deployPath)
    }
}
