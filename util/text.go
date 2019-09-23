package util

import (
	"regexp"
)

// RegSplit Split a paragraph into sentences.
// Note: unlike the standard regexp.Split, the delimeters are retained
func RegSplit(text string, delimeter string) []string {
	reg := regexp.MustCompile(delimeter)
	indexes := reg.FindAllStringIndex(text, -1)
	haslast := len(indexes) > 0 && indexes[len(indexes)-1][1] == len(text)
	laststart := 0
	var result []string
	if haslast {
		result = make([]string, len(indexes))
	} else {
		result = make([]string, len(indexes)+1)
	}
	for i, element := range indexes {
		// retain the delimiter
		result[i] = text[laststart:element[1]]
		laststart = element[1]
	}
	if !haslast {
		result[len(indexes)] = text[laststart:len(text)]
	}
	return result
}

// splitSents Split sentences with Chinese punctuations
func SplitSents(para string) []string {
	arr := []string{}
	split := RegSplit(para, "([。！：？]”?)+( )+(\\[#....\\])?|[\\.\\!\\:\\?]( )+(\\[#....\\])?")
	for i := range split {
		arr = append(arr, split[i])
	}
	return arr
}
