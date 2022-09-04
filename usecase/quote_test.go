package usecase

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func Test_HandleAllQuoteBeforePic(t *testing.T) {
	content, err := os.ReadFile("./collapse_demo.md")
	if err != nil {
		log.Fatal(err)
	}
	c := string(content)
	HandleAllQuoteBeforePic(&c)
	fmt.Println(c)
}
