package base

import "bytes"

func ToEventString(s string, maxLen int) string {
	if maxLen == 0 {
		const defaultLen = 100
		maxLen = defaultLen
	}
	type stageT int
	const (
		stageFirstChar stageT = iota
		stageSecondChar
		stageOthers
	)
	stage := stageFirstChar
	upperCase := true
	var buf bytes.Buffer
	for i, ll := 0, len(s); i < ll; i++ {
		if buf.Len() >= maxLen {
			break
		}
		ch := s[i]
		if !isChar(ch) {
			stage = stageFirstChar
			continue
		}
		switch stage {
		case stageFirstChar:
			buf.WriteByte(byteToUpper(ch))
			stage = stageSecondChar
		case stageSecondChar:
			upperCase = ch >= 'A' && ch <= 'Z'
			buf.WriteByte(byteToLower(ch))
			stage = stageOthers
		default:
			if upperCase {
				buf.WriteByte(byteToLower(ch))
			} else {
				buf.WriteByte(ch)
			}
		}

	}
	return buf.String()
}

func byteToUpper(ch byte) byte {
	if ch >= 'a' && ch <= 'z' {
		return ch - 'a' + 'A'
	} else {
		return ch
	}
}

func byteToLower(ch byte) byte {
	if ch >= 'A' && ch <= 'Z' {
		return ch - 'A' + 'a'
	} else {
		return ch
	}
}

func isChar(ch byte) bool {
	if ch >= 'A' && ch <= 'Z' {
		return true
	}
	if ch >= 'a' && ch <= 'z' {
		return true
	}
	return false
}
