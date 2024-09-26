package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"
)

var (
	origDir   string
	trash_dir = ".wait_clean"
)

func readDirRecursion(dirName string) (files []interface{}, err error) {
	fs, err := os.ReadDir(dirName)
	if err != nil {
		return nil, err
	}
	for _, f := range fs {
		if f.IsDir() {
			childrenFiles, err := readDirRecursion(path.Join(dirName, f.Name()))
			if err != nil {
				return files, err
			}
			files = append(files, childrenFiles...)
		} else {
			info, _ := f.Info()
			if (strings.HasPrefix(f.Name(), "._") && info.Size() == 4096) ||
				f.Name() == ".DS_Store" {
				files = append(files, path.Join(dirName, f.Name()))
			}
		}
	}
	return
}

func moveFileRecursion(files []interface{}) error {
	for _, f := range files {
		switch f := f.(type) {
		case string:
			err := os.Rename(f, path.Join(trash_dir, path.Base(f)))
			fmt.Println("已归档", f)
			if err != nil {
				return err
			}
		case []interface{}:
			err := moveFileRecursion(f)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func main() {
	// command line tool
	var cmdlDir string
	flag.StringVar(&cmdlDir, "t", "", "Clean the target directory, the default is the current directory")
	flag.Parse()
	origDir = cmdlDir

	if origDir == "" {
		origDir, _ = os.Getwd()
	}
	trash_dir = path.Join(origDir, trash_dir)

	os.Mkdir(trash_dir, 0755)
	files, err := readDirRecursion(origDir)
	if err != nil {
		fmt.Println(err)
	}
	if files == nil {
		os.Remove(trash_dir)
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
		err = os.RemoveAll(trash_dir)
		if err != nil {
			fmt.Println(err)
		}
	}
}
