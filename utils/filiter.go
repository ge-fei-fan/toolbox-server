package utils

import "regexp"

// FilterFilename  文件名过滤特殊字符
func FilterFilename(src string) string {
	reg, _ := regexp.Compile(`[/ : * ? " < > | \\]`)
	dst := reg.ReplaceAllString(src, " ")
	if dst == "" {
		dst = "无名"
	}
	return dst
}

// FilterString 提取第一次匹配的字符串
func FilterString(role, src string) string {
	reg, _ := regexp.Compile(role)
	return reg.FindString(src)
}
