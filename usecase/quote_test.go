package usecase

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"util-cli/consts"
)

func Test_HandleAllQuoteBeforePic(t *testing.T) {
	content, err := os.ReadFile("./collapse_demo.md")
	if err != nil {
		log.Fatal(err)
	}
	c := string(content)
	HandleAllQuote(&c)
	fmt.Println(c)
}

func Test_QuoteTable(t *testing.T) {
	s := "|参数|解释|"
	flag := strings.HasPrefix(s, consts.QUOTE_IDENTIFIER)
	fmt.Println(flag)
}
