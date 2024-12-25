package main

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

func main() {
	args := os.Args[1:]
	refreshPath := ""
	for {
		if len(args) <= 0 {
			break
		}
		arg := args[0]
		if arg == "-h" || arg == "--help" {
			PrintHelp()
		} else {
			// Figure if is path
			_, e := os.Stat(arg)
			if e == nil {
				refreshPath = arg
			}
		}
		if len(args) <= 1 {
			break
		}
		args = args[1:]
	}
	if refreshPath == "" {
		ex, err := os.Executable()
		if err != nil {
			panic(err)
		}
		exPath := filepath.Dir(ex)
		RecursiveRefresh(exPath)
	} else {
		refreshPath, e := filepath.Abs(refreshPath)
		if e != nil {
			panic(e)
		}
		slog.Debug("Refreshing path: " + refreshPath)
		RecursiveRefresh(refreshPath)
	}
}

func RecursiveRefresh(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// path does not exist or other issue
		panic(err)
	}
	file, e := os.Stat(path)
	if e != nil {
		panic(e)
	}
	if file.IsDir() {
		items, e := os.ReadDir(path)
		if e != nil {
			panic(e)
		}
		for _, i := range items {
			RecursiveRefresh(path + "/" + i.Name())
		}
	} else {
		slog.Debug("Refreshing item: " + path)
		e := os.Chtimes(path, time.Now(), time.Now())
		if e != nil {
			slog.Error(e.Error())
		}
	}
}

func PrintHelp() {
	fmt.Println("Refresher program. Used to refresh access and modify datetime. Useful if something deletes unused files.")
	fmt.Println("Use: refresher [args] <path>")
	fmt.Println("Possible arguments:")
	fmt.Println("\t-h\t--help\tPrints out this help")
	fmt.Println("<path>\tPath to be refreshed. Includes subdirectories. If omitted, program will use current folder.")
}
