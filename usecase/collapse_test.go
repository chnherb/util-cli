package usecase

import (
	"fmt"
	"log"
	"os"
	"testing"
)
func Test_CollapseCode(t *testing.T) {
	content, err := os.ReadFile("./collapse_demo.md")
	if err != nil {
		log.Fatal(err)
	}
	c := string(content)
	CollapseCodeWithLine(&c, 3,false)
	fmt.Println(c)
}

