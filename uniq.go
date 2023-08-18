package main

func UniqContents(contents ContentType, key int) ContentType {
	var result ContentType = append(ContentType{}, contents...)
	for i := 0; ; {
		if i == len(result)-1 {
			break
		}
		if result[i].Fields[key-1] == result[i+1].Fields[key-1] {
			result = append(result[:i+1], result[i+2:]...)
		} else {
			i++
		}
	}
	return result
}
