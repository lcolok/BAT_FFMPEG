package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

type Object struct {
	From     float64 `json:"from"`
	To       float64 `json:"to"`
	Sid      int     `json:"sid"`
	Location int     `json:"location"`
	Content  string  `json:"content"`
	Music    int     `json:"music"`
}

func main() {
	dir := "./json"
	txtDir := "./txt"
	// 读取目录中的文件
	_, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
		return
	}

	os.MkdirAll(txtDir, os.ModePerm)

	processContent := func(contentList []string) string {
		// 纠错器函数：将所有包含 "l o r a" 和 "l i a" 的字眼改为 "LoRA"
		for i := range contentList {
			contentList[i] = strings.Replace(contentList[i], "l o r a", "LoRA", -1)
			contentList[i] = strings.Replace(contentList[i], "l i a", "LoRA", -1)
			contentList[i] = strings.Replace(contentList[i], "l r a", "LoRA", -1)

		}
		return strings.Join(contentList, "，")
	}

	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(path) != ".json" {
			return nil
		}

		fileContent, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		var objects []Object
		err = json.Unmarshal(fileContent, &objects)
		if err != nil {
			return err
		}

		var contentList []string
		for _, obj := range objects {
			contentList = append(contentList, obj.Content)
		}

		article := processContent(contentList)

		txtFileName := strings.TrimSuffix(info.Name(), filepath.Ext(info.Name())) + ".txt"
		err = ioutil.WriteFile(filepath.Join(txtDir, txtFileName), []byte(article), os.ModePerm)
		if err != nil {
			return err
		}

		fmt.Printf("已生成文章 %s\n", txtFileName)

		return nil
	})

	if err != nil {
		fmt.Println(err)
	}
}

func isEndOfSentence(s string) bool {
	r, _ := utf8.DecodeLastRuneInString(s)
	return r == '。' || r == '？' || r == '！'
}
