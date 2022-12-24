package main

import (
	"archive/tar"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

func getArchiveName(filename, path string) string {
	var name string
	var timestamp string

	if strings.HasSuffix(filename, ".log") {
		name = strings.TrimSuffix(filename, ".log")
	} else {
		log.Fatalln("Incorrect file extension: ", filename)
	}

	timestamp = strconv.Itoa(int(time.Now().Unix()))
	name += "_" + timestamp + ".tar.gz"

	if path != "" {
		name = path + "/" + filepath.Base(name)
	}
	return name
}

func createArchive(filename string, out io.Writer) error {
	gw := gzip.NewWriter(out)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()

	err := addToArchive(tw, filename)
	if err != nil {
		return err
	}

	return nil
}

func addToArchive(tw *tar.Writer, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	header, err := tar.FileInfoHeader(info, info.Name())
	if err != nil {
		return err
	}

	header.Name = filename

	err = tw.WriteHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(tw, file)
	if err != nil {
		return err
	}

	return nil
}

func archive(filename, path string, group *sync.WaitGroup) {
	archname := getArchiveName(filename, path)

	out, err := os.Create(archname)
	if err != nil {
		log.Fatalln("Error writing archive:", err)
	}
	defer out.Close()

	err = createArchive(filename, out)
	if err != nil {
		log.Fatalln("Error creating archive:", err)
	}

	fmt.Println("Archive created successfully")

	group.Done()
}

func main() {
	a := flag.String("a", "", "put archive into specified directory")
	flag.Parse()

	if *a == "" && len(os.Args) >= 2 {
		var group sync.WaitGroup
		for i := 1; i < len(os.Args); i++ {
			group.Add(1)
			go archive(os.Args[i], *a, &group)
		}
		group.Wait()
	} else if *a != "" && len(os.Args) >= 4 {
		var group sync.WaitGroup
		for i := 3; i < len(os.Args); i++ {
			group.Add(1)
			go archive(os.Args[i], *a, &group)
		}
		group.Wait()
	} else {
		fmt.Println("Pass one or more filenames to archive them")
		fmt.Println("You can also use flag '-a' to specify directory where archives will be stored")
	}
}
