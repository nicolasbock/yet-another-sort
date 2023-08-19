package main

func UniqContents(contents ContentType) ContentType {
	var result ContentType = append(ContentType{}, contents...)
	if uniq {
		for i := 0; ; {
			if i == len(result)-1 {
				break
			}
			if result[i].CompareLine == result[i+1].CompareLine {
				result = append(result[:i+1], result[i+2:]...)
			} else {
				i++
			}
		}
	}
	return result
}
