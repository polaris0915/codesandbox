package _string

import "strings"

// GetCorrectString 这个是以防数据库的问题的测试用例中的输入和输出是`1 2\\n1 2 3 4 5\\n`这种将`\n`当作了`\\n`
func GetCorrectString(str *string) {
	*str = strings.ReplaceAll(*str, "\\n", "\n")
	EndWithBr(str)
}

func EndWithBr(str *string) {
	if !strings.HasSuffix(*str, "\n") {
		*str += "\n"
	}
}

func EndWithoutBr(str *string) {
	*str, _ = strings.CutSuffix(*str, "\n")
}
