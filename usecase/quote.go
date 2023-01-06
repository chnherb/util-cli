package usecase

import (
	"fmt"
	"github.com/gookit/slog"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"util-cli/consts"
	"util-cli/utils"
)

type QuoteArgs struct {
	Src string
}

func BatchQuote(args *QuoteArgs) error {
	if err := checkQuoteArgs(args); err != nil {
		return err
	}
	return RecursiveQuote(args, args.Src)
}

func RecursiveQuote(args *QuoteArgs, dirName string) error {
	fileInfos, err := ioutil.ReadDir(dirName)
	if err != nil {
		return err
	}
	for _, f := range fileInfos {
		curName := f.Name()
		subName := fmt.Sprintf("%s/%s", dirName, f.Name())
		if strings.HasSuffix(curName, ".md") {
			QuoteFile(args, subName)
		}
		if f.IsDir() {
			RecursiveQuote(args, subName)
		}
	}
	return nil
}

func QuoteFile(args *QuoteArgs, subPath string) error {
	path := subPath
	if !utils.ExistFile(path) {
		return fmt.Errorf("there is no file in %s", path)
	}
	contentByte, err := ioutil.ReadFile(path)
	content := string(contentByte)
	if err != nil {
		return err
	}
	if !HandleAllQuote(&content) {
		return nil
	}
	err = ioutil.WriteFile(path, []byte(content), 0666)
	if err != nil {
		return err
	}
	slog.Infof("Collapse CodeFile: %s", path)
	return nil
}

func checkQuoteArgs(args *QuoteArgs) error {
	if args.Src == "" {
		src, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("--src is empty and \"pwd\" error: %v", err)
		}
		args.Src = src
	}
	return nil
}

// check > Before the image in a line
func CheckQuoteInLastLine(preContent string) bool {
	if len(preContent) <= 0 {
		return false
	}
	if !strings.HasSuffix(preContent, consts.LINE_IDENTIFIER) {
		return false
	}
	rmTailPreContent := preContent[:len(preContent)-1]
	i := strings.LastIndex(rmTailPreContent, consts.LINE_IDENTIFIER)
	lastLine := rmTailPreContent
	if i < len(rmTailPreContent) && i >= 0 {
		lastLine = rmTailPreContent[i+1:]
	}
	if strings.HasPrefix(lastLine, consts.QUOTE_IDENTIFIER) {
		return true
	}
	return false
}

// handle > Before the image on a line
func HandleAllQuoteBeforePic(content *string) bool {
	//imgReStr := `\!\[*?\]\(*?\)`
	imgReStr := `\!\[.*\]\(.+\)`
	b, err := regexp.MatchString(imgReStr, *content)
	if !b || err != nil {
		return false
	}
	regexp, err := regexp.Compile(imgReStr)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	indexs := regexp.FindAllIndex([]byte(*content), -1)
	newContent := ""
	preIndex := 0
	if len(indexs) == 0 {
		return false
	}
	result := false
	for _, index := range indexs {
		//fmt.Println(fmt.Sprintf("%d, %d", index[0], index[1]))
		preContent := (*content)[preIndex:index[0]]
		newContent += preContent
		if CheckQuoteInLastLine(preContent) {
			newContent += consts.LINE_IDENTIFIER
			result = true
		}
		picContent := (*content)[index[0]:index[1]]
		newContent += picContent
		preIndex = index[1]
	}
	if preIndex < len(*content) {
		newContent += (*content)[preIndex:]
	}
	*content = newContent
	return result
}

func HandleAllQuote(content *string) bool {
	ss := strings.Split(*content, consts.LINE_IDENTIFIER)
	if len(ss) < 2 {
		return false
	}
	result := false
	newContent := ss[0]
	for i := 1; i < len(ss); i++ {
		if strings.HasPrefix(ss[i-1], consts.QUOTE_IDENTIFIER) &&
			strings.TrimSpace(ss[i]) != consts.BLANK_IDENTIFIER &&
			!strings.HasPrefix(ss[i], consts.QUOTE_IDENTIFIER) {
			newContent += consts.LINE_IDENTIFIER
			result = true
		}
		newContent += consts.LINE_IDENTIFIER + ss[i]
	}
	*content = newContent
	return result
}
