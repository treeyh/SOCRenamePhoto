package main

import (
	"fmt"
	"github.com/rwcarlsen/goexif/exif"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func getFileList(path string) []string {
	fileList := []string{}

	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		fileList = append(fileList, path)
		return nil
	})

	if err != nil {
		fmt.Printf("filepath.Walk() returned %v\n", err)
	}

	return fileList
}

func buildFileName(paths []string, fileName, sep string, i int) string {
	if i == 0 {
		paths[len(paths)-1] = fileName + ".jpg"
	} else {
		s := fmt.Sprintf("%d", i)
		paths[len(paths)-1] = fileName + "_" + s + ".jpg"
	}
	newPath := ""
	for _, n := range paths {
		newPath += n + sep
	}
	newPath = newPath[0 : len(newPath)-1]
	return newPath
}

func renameImg(path, fix, sep string, ch chan int) {
	f, err := os.Open(path)
	if err != nil {
		fmt.Printf("renameImg() error %v\n", err)
	}
	x, err := exif.Decode(f)
	if err != nil {
		fmt.Printf("renameImg() error %v\n", err)
	}
	date, _ := x.Get(exif.DateTimeOriginal)
	f.Close()

	datestr := strings.Replace(strings.Replace(date.StringVal(), ":", "", -1), " ", "", -1)
	fileName := datestr + "_" + fix
	ss := strings.Split(path, sep)

	newPath := buildFileName(ss, fileName, sep, 0)
	err = os.Rename(path, newPath)
	if err != nil {
		for i := 1; i <= 1000000; i++ {

			newPath = buildFileName(ss, fileName, sep, i)
			err = os.Rename(path, newPath)
			if err != nil {
				continue
			}
			break
		}
	}
	ch <- 1
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	//照片保存目录
	path := "D:\\1111"
	//照片中间名字最终定义为yyyyMMddHHmmss_{name}.jpg
	name := "春游踏青"
	//windows的文件夹层级分隔符
	sep := "\\"

	start := time.Now()
	fileList := getFileList(path)
	ch := make(chan int)
	index := 0
	for _, filePath := range fileList {
		s := strings.ToLower(filePath[len(filePath)-4 : len(filePath)])
		if s != ".jpg" {
			continue
		}
		index++
		go renameImg(filePath, name, sep, ch)
	}
	fmt.Println(index)

	i := 0
L:
	for {
		select {
		case <-ch:
			i++
			if i >= index {
				break L
			}
		}
	}

	fmt.Println(i)
	end := time.Now()
	hs := end.UnixNano() - start.UnixNano()
	fmt.Println(hs)
}
