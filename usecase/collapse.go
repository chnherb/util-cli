package usecase

import (
	"fmt"
	"github.com/gookit/slog"
	"io/ioutil"
	"os"
	"strings"
	"util-cli/utils"
)

const (
	CODE_LINE_LIMIT        = 20
	CODE_COLLAPSE_TITILE   = "Expand/Collapse Code Block"
	CODE_IDENTIFIER        = "```"
	CODE_COLLAPSE_TEMPLATE = "{{%% code ctitle=\"%s\" %%}}\n%s%s%s\n{{%% /code %%}}\n\n" // %% 转义成 %（不是\\%）
	LINE_IDENTIFIER        = "\n"
)

var (
	// https://www.jianshu.com/p/1f223eb78ad8
	MD_LANGs = []string{"text", "plain", "shell", "bash", "json", "yaml", "yml", "xml", "xhtml", "html",
		"java", "go", "py", "python", "cpp", "c#", "c-sharp", "scala", "objc", "obj-c", "swift", "pl", "perl",
		"sql",
		"css", "js", "jscript", "javascript", "php"}
)

type CollapseArgs struct {
	Src       string
	MinRow    int
	NeedTitle bool
}

func BatchCollapseCode(args *CollapseArgs) error {
	if err := checkCollapseCodeArgs(args); err != nil {
		return err
	}
	return RecursiveCollapseCode(args, args.Src)
}

func RecursiveCollapseCode(args *CollapseArgs, dirName string) error {
	fileInfos, err := ioutil.ReadDir(dirName)
	if err != nil {
		return err
	}
	for _, f := range fileInfos {
		curName := f.Name()
		subName := fmt.Sprintf("%s/%s", dirName, f.Name())
		if strings.HasSuffix(curName, ".md") {
			CollapseCodeFile(args, subName)
		}
		if f.IsDir() {
			RecursiveCollapseCode(args, subName)
		}
	}
	return nil
}

func CollapseCodeFile(args *CollapseArgs, subPath string) error {
	path := subPath
	codeLineLimit := args.MinRow
	needTitle := args.NeedTitle
	if !utils.ExistFile(path) {
		return fmt.Errorf("there is no file in %s", path)
	}
	contentByte, err := ioutil.ReadFile(path)
	content := string(contentByte)
	if err != nil {
		return err
	}
	if !CollapseCodeWithLine(&content, codeLineLimit, needTitle) {
		return nil
	}
	err = ioutil.WriteFile(path, []byte(content), 0666)
	if err != nil {
		return err
	}
	slog.Infof("Collapse CodeFile: %s", path)
	return nil
}

func checkCollapseCodeArgs(args *CollapseArgs) error {
	if args.Src == "" {
		src, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("--src is empty and \"pwd\" error: %v", err)
		}
		args.Src = src
	}
	return nil
}

func CollapseCode(content *string, needTitle bool) bool {
	return CollapseCodeWithLine(content, CODE_LINE_LIMIT, needTitle)
}

func CollapseCodeWithLine(content *string, codeLineLimit int, needTitle bool) bool {
	ss := strings.Split(*content, CODE_IDENTIFIER)
	newContent := ""
	matchCode := false
	preMatch := false
	for idx := 0; idx < len(ss)-1; idx++ {
		curCode := ss[idx]
		matchCode = false // reset
		if hasMDLang(curCode) && !hasMDLang(ss[idx+1]) {
			if strings.Count(curCode, LINE_IDENTIFIER) >= codeLineLimit {
				collapseTitle := parseCollapseTitle(curCode, needTitle)
				newContent += fmt.Sprintf(CODE_COLLAPSE_TEMPLATE, collapseTitle, CODE_IDENTIFIER, curCode, CODE_IDENTIFIER)
				matchCode = true
			}
		}
		if !matchCode {
			if idx == 0 || preMatch {
				newContent += curCode
			} else {
				newContent += fmt.Sprintf("%s%s", CODE_IDENTIFIER, curCode)
			}
		}
		preMatch = matchCode
	}
	if newContent == "" {
		newContent = *content
		return false
	}
	if !matchCode {
		newContent += fmt.Sprintf("%s%s", CODE_IDENTIFIER, ss[len(ss)-1])
	} else {
		newContent += ss[len(ss)-1]
	}

	*content = newContent
	return true
}

func hasMDLang(code string) bool {
	for _, lang := range MD_LANGs {
		if strings.HasPrefix(code, lang) {
			return true
		}
	}
	return false
}

func parseCollapseTitle(code string, needTitle bool) string {
	collapseTitle := CODE_COLLAPSE_TITILE
	if !needTitle {
		return collapseTitle
	}
	funcKeymap := map[string][]string{
		"java": {"public", "protected", "private"},
		"go":   {"func"},
	}

	langType := ""
	if strings.HasPrefix(code, "java") {
		langType = "java"
	} else if strings.HasPrefix(code, "go") {
		langType = "go"
	}
	if keys, ok := funcKeymap[langType]; ok {
		ss := strings.Split(code, LINE_IDENTIFIER)
		lineCount := 5
		for i := 1; i <= lineCount; i++ {
			if i >= len(ss) {
				break
			}
			for _, key := range keys {
				if strings.Contains(ss[i], key) {
					return strings.TrimSpace(ss[i])
				}
			}
		}
	}
	return collapseTitle
}
