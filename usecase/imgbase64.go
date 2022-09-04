package usecase

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	_ "strings"
	"util-cli/consts"
	"util-cli/utils"

	"github.com/gookit/slog"
	"github.com/olekukonko/tablewriter"
)

type ImgBase64Args struct {
	Src     string
	Rewrite bool
	Chapter string
}

func ParseImgType(imgInfo string) (string, int) {
	// imgInfo := "![图片](data:image/jpeg;base64,/9j/4AAQSkZJR"
	// imgInfo := content[preIndex:index[0]]
	imgTag := "image/"
	base64Tag := ";base64,"
	tag := fmt.Sprintf("%s.*?%s", imgTag, base64Tag)
	regexp := regexp.MustCompile(tag)
	imgPrefix := regexp.FindString(imgInfo)
	slog.Debugf("imgPrefix: %s", imgPrefix)
	imgType := imgPrefix[len(imgTag) : len(imgPrefix)-len(base64Tag)]
	slog.Debugf("imgType: %s", imgType)
	return imgType, len(imgPrefix)
}

func SaveImg(dir, chapter string, index int, base64ImgContent string) (string, error) {
	imgReStr := `^\!\[图片\]\(data:\s*image\/(\w+);base64,`
	b, err := regexp.MatchString(imgReStr, base64ImgContent)
	if !b || err != nil {
		return "", err
	}
	re, err := regexp.Compile(imgReStr)
	if err != nil {
		return "", err
	}
	allData := re.FindAllSubmatch([]byte(base64ImgContent), 2)
	imgExt := string(allData[0][1])
	imgFilePath := fmt.Sprintf("./imgs/%s_%d.%s", chapter, index+1, imgExt)
	filePath := fmt.Sprintf("%s/%s", dir, imgFilePath)
	fileDir := filepath.Dir(filePath)
	utils.CreateDirIfNon(fileDir)
	base64Str := re.ReplaceAllString(base64ImgContent, "")
	byteImg, _ := base64.StdEncoding.DecodeString(base64Str)
	err = ioutil.WriteFile(filePath, byteImg, 0666)
	if err != nil {
		return "", err
	}
	return imgFilePath, nil
}

func ParseImgBase64File(path string, chapter string, rewrite bool) error {
	if !utils.ExistFile(path) {
		return fmt.Errorf("there is no file in %s", path)
	}
	contentByte, err := ioutil.ReadFile(path)
	content := string(contentByte)
	if err != nil {
		return err
	}
	newContent, err := ParseImgBase64Content(path, content, chapter)
	if err != nil {
		return err
	}
	if newContent == "" && hasHugoHeader(&content) {
		return nil
	}
	if newContent == "" {
		newContent = content
	}
	HandleHugoHeader(&newContent, path)
	newPath := path
	if !rewrite {
		newPath = strings.Replace(path, ".md", "_01.md", 1)
	}
	err = ioutil.WriteFile(newPath, []byte(newContent), 0666)
	if err != nil {
		return err
	}
	slog.Infof("ParseImgBase64File: %s", path)
	return nil
}

func ParseImgBase64Content(path, content string, chapter string) (string, error) {
	// chapter := "git_operation"
	imgReStr := `\!\[图片\]\(data:image.*?\)`
	b, err := regexp.MatchString(imgReStr, content)
	if !b || err != nil {
		return "", err
	}
	regexp, err := regexp.Compile(imgReStr)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	// ss := regexp.FindAllString(content, -1)
	// fmt.Println(ss)
	indexs := regexp.FindAllIndex([]byte(content), -1)
	newContent := ""
	preIndex := 0
	if len(indexs) == 0 {
		return "", nil
	}
	imgFilePaths := []string{}
	for i, index := range indexs {
		// fmt.Println(fmt.Sprintf("%d, %d", index[0], index[1]))
		base64ImgContent := content[index[0]:index[1]]
		imgFilePath, err := SaveImg(filepath.Dir(path), chapter, i, base64ImgContent)
		if err != nil {
			fmt.Errorf("SaveImg ERR, i: %d", i)
		}
		imgFilePaths = append(imgFilePaths, imgFilePath)
		preContent := content[preIndex:index[0]]
		newContent += preContent
		if CheckQuoteInLastLine(preContent) {
			newContent += consts.LINE_IDENTIFIER
		}
		newContent += fmt.Sprintf("![%s](%s)", filepath.Base(imgFilePath), imgFilePath)
		// slog.Debugf("c=%s\n", newContent)
		preIndex = index[1]
	}
	if preIndex < len(content) {
		newContent += content[preIndex:]
	}
	// slog.Debugf("c=%s\n", newContent)
	showImgFilePath(imgFilePaths)
	return newContent, nil
}

func HandleHugoHeader(content *string, path string) {
	if !hasHugoHeader(content) {
		base := filepath.Base(path)
		ext := filepath.Ext(path)
		filename := strings.TrimSuffix(base, ext)
		template := `---
categories: [""]
tags: [""]
title: "%s"
# linkTitle: ""
weight: 10
description: >

---

`
		header := fmt.Sprintf(template, filename)
		*content = header + *content
	}
}

func hasHugoHeader(content *string) bool {
	ss := strings.Split(*content, "\n")
	preLen := 5
	if len(ss) < preLen {
		preLen = len(ss)
	}
	lines := strings.Join(ss[:preLen], "\n")
	if strings.Contains(lines, "---") && strings.Contains(lines, "title:") {
		return true
	}
	return false
}

func ParseImgBase64Dir(args *ImgBase64Args) {

}

func ParseImgBase64ContentWithAltText(content string, altText string) error {
	regexStr := fmt.Sprintf("![%s](*)", altText)
	regexp, _ := regexp.Compile(regexStr)
	fmt.Println(regexp.FindString(content))
	return nil
}

func RecursiveHandleFiles(dirName string, chapter string, rewrite bool) error {
	fileInfos, err := ioutil.ReadDir(dirName)
	if err != nil {
		return err
	}
	for _, f := range fileInfos {
		curName := f.Name()
		subName := fmt.Sprintf("%s/%s", dirName, f.Name())
		if strings.HasSuffix(curName, ".md") {
			ParseImgBase64File(subName, chapter, rewrite)
		}
		if f.IsDir() {
			RecursiveHandleFiles(subName, chapter, rewrite)
		}
	}
	return nil
}

func checkImgBaser64Args(args *ImgBase64Args) error {
	if args.Src == "" {
		src, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("--src is empty and \"pwd\" error: %v", err)
		}
		args.Src = src
	}
	if args.Chapter == "" {
		return fmt.Errorf("--chapter is empty")
	}
	return nil
}

func Replace(args *ImgBase64Args) error {
	if err := checkImgBaser64Args(args); err != nil {
		return err
	}
	return RecursiveHandleFiles(args.Src, args.Chapter, args.Rewrite)
}

func showImgFilePath(imgFilePaths []string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Dir", "FileName"})
	for _, path := range imgFilePaths {
		table.Append([]string{filepath.Dir(path), filepath.Base(path)})
	}
	table.Render()
}
