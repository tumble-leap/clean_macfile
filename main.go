package main

import (
	"fmt"
	"os"
	"path"
	"strings"
)

var (
	origDir   = "./"
	resultDir = ".wait_clean"
)

func readDirRecursion(dirName string) (files []interface{}, err error) {
	f, err := os.ReadDir(dirName)
	if err != nil {
		return nil, err
	}
	for _, v := range f {
		if v.IsDir() {
			files1, err := readDirRecursion(path.Join(dirName, v.Name()))
			if err != nil {
				return files, err
			}
			files = append(files, files1...)
		} else {
			info, _ := v.Info()
			if (strings.HasPrefix(v.Name(), "._") && info.Size() == 4096) ||
				v.Name() == ".DS_Store" {
				files = append(files, path.Join(dirName, v.Name()))
			}
		}
	}
	return
}

func moveFileRecursion(files []interface{}) error {
	for _, v := range files {
		switch v := v.(type) {
		case string:
			err := os.Rename(v, path.Join(resultDir, path.Base(v)))
			fmt.Println("已归档", v)
			if err != nil {
				return err
			}
		case []interface{}:
			err := moveFileRecursion(v)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func main() {
	os.Mkdir(resultDir, 0755)
	files, err := readDirRecursion(origDir)
	if err != nil {
		fmt.Println(err)
	}
	if files == nil {
		os.Remove(resultDir)
		return
	}
	err = moveFileRecursion(files)
	if err != nil {
		fmt.Println(err)
	}
	var i = "y"
	fmt.Print("是否要删除(Y/n)：")
	fmt.Scanln(&i)
	if i == "y" || i == "Y" {
		err = os.RemoveAll(resultDir)
		if err != nil {
			fmt.Println(err)
		}
	}
}
