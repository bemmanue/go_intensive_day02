package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func showSymlinks(path string) {
	filepath.WalkDir(path, func(p string, info os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !info.Type().IsDir() && !info.Type().IsRegular() {
			point, err := filepath.EvalSymlinks(p)
			if err != nil {
				fmt.Printf("%s -> [broken]\n", p)
			} else {
				fmt.Printf("%s -> %s\n", p, point)
			}
		}
		return nil
	})
}

func showFiles(path string) {
	filepath.WalkDir(path, func(p string, info os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if info.Type().IsRegular() {
			fmt.Println(p)
		}
		return nil
	})
}

func showExtFiles(path string, ext string) {
	ext = "." + ext

	filepath.WalkDir(path, func(p string, info os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(info.Name()) == ext {
			fmt.Println(p)
		}
		return nil
	})
}

func showDirectories(path string) {
	filepath.WalkDir(path, func(p string, info os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			fmt.Println(p)
			fmt.Println(p)
		}
		return nil
	})
}

func showAll(path string) {
	showDirectories(path)
	showFiles(path)
	showSymlinks(path)
}

func main() {
	sl := flag.Bool("sl", false, "show symlinks")
	d := flag.Bool("d", false, "show directories")
	f := flag.Bool("f", false, "show files")
	ext := flag.String("ext", "", "show files with a certain extension")
	flag.Parse()

	path := flag.Arg(0)
	if !*sl && !*d && !*f {
		showAll(path)
	}
	if *sl {
		showSymlinks(path)
	}
	if *f {
		if *ext != "" {
			showExtFiles(path, *ext)
		} else {
			showFiles(path)
		}
	}
	if *d {
		showDirectories(path)
	}
}
