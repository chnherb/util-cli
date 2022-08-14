package usecase

import (
	"fmt"
	"path/filepath"
	"regexp"
	"testing"
)

func Test_ParsePicBase64File(t *testing.T) {
	ParseImgBase64File("/Users/bo/hb_blog/00-util/git/git常用操作.md", "test", false)
}

func Test_ParsePicBase64Content(t *testing.T) {
	content := "sfd![图片](24hk34kidfsdf)sdjfksdjfk  kjdk ![图片](135234) dsfsfs "
	// regexp, err := regexp.Compile("![图片].*")
	// regexp := regexp.MustCompile(`\!\[.*\]\(.+\)`)
	regexp := regexp.MustCompile(`\!\[图片\]\(.*?\)`)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	ss := regexp.FindAllString(content, -1)
	fmt.Println(ss)
	indexs := regexp.FindAllIndex([]byte(content), -1)
	newContent := ""
	preIndex := 0
	for i, index := range indexs {
		fmt.Println(fmt.Sprintf("%d, %d", index[0], index[1]))
		newContent += content[preIndex:index[0]]
		newContent += fmt.Sprintf("![tupian](%d)", i)
		fmt.Printf("c=%s\n", newContent)
		preIndex = index[1]
	}
	newContent += content[preIndex:]
	fmt.Printf("c=%s\n", newContent)
}

func Test_PareseImgSuffix(t *testing.T) {
	picInfo := "![图片](data:image/jpeg;base64,/9j/4AAQSkZJR"
	// picInfo := content[preIndex:index[0]]
	imgTag := "image/"
	base64Tag := ";base64"
	tag := fmt.Sprintf("%s.*?%s", imgTag, base64Tag)
	regexp := regexp.MustCompile(tag)
	ss := regexp.FindString(picInfo)
	fmt.Printf("type: %s", ss)
	fmt.Printf("type: %s", ss[len(imgTag):len(ss)-len(base64Tag)])
}

func Test_filepath(t *testing.T) {
	path := "./img/test.png"
	dir := filepath.Dir(path)

	fmt.Printf("dir: %s\n", dir)
	s := filepath.Base(path)
	fmt.Println(s)
}
