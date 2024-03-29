package envar

import (
	"bytes"
	"fmt"
	"io"
	"unicode"
)

const (
	exportPrefix  = "export"
	inlineComment = " #"
)

func parse(buf *bytes.Buffer) (envMap map[string]string, err error) {
	lines, err := getValidLines(buf)
	if err != nil {
		return
	}

	envMap = make(map[string]string)

	for _, line := range lines {
		key, rest, err := extractKey(line)
		if err != nil {
			return nil, err
		}
		value, err := extractValue(rest, &envMap)
		if err != nil {
			return nil, err
		}

		envMap[key] = value
	}

	return
}

func extractKey(src []byte) (key string, rest []byte, err error) {
	endOfKey := bytes.IndexByte(src, '=')
	if endOfKey == -1 {
		err = fmt.Errorf("invalid env variable format in line %q\n", string(src))
		return
	}

	key = string(sanitizeKey(src[:endOfKey]))
	if len(key) == 0 {
		err = fmt.Errorf("empty key in line %q\n", string(src))
	}

	rest = bytes.TrimFunc(src[endOfKey+1:], unicode.IsSpace)
	return
}

func extractValue(src []byte, envMap *map[string]string) (value string, err error) {
	if len(src) == 0 {
		return
	}

	isQuoted, endOfQuote, quote, err := isQuoted(src)
	if err != nil {
		return
	}

	if isQuoted {
		// no need in removing inline comment in the "quoted" case
		// as far as it was omitted in slicing up to `endOfQuote` value
		src = bytes.ReplaceAll(src[1:endOfQuote], []byte{'\\', quote}, []byte{quote})
	} else {
		if bytes.Contains(src, []byte(inlineComment)) {
			src = bytes.Split(src, []byte(inlineComment))[0]
			src = bytes.TrimRightFunc(src, unicode.IsSpace)
		}
	}

	src, err = substituteValue(src, envMap)

	value = string(src)
	return
}

func substituteValue(src []byte, envMap *map[string]string) (substituted []byte, err error) {
	if len(*envMap) == 0 {
		return src, nil
	}

	substituted = src

	for key, value := range *envMap {
		target := []byte(key)
		targetValue := []byte(value)

		target = append([]byte{'$', '{'}, target...)

		targetIdx := bytes.Index(substituted, target)
		for targetIdx != -1 {
			targetEnd := substituted[targetIdx+len(target)]
			if targetEnd != '}' {
				err = fmt.Errorf("invalid substitution format in %q", string(src))
				return
			}

			substituted = bytes.Replace(substituted, append(target, targetEnd), targetValue, 1)

			targetIdx = bytes.Index(substituted, target)
		}

		target = append([]byte{'$'}, target[2:]...)
		targetIdx = bytes.Index(substituted, target)
		for targetIdx != -1 {
			substituted = bytes.Replace(substituted, target, targetValue, 1)
			targetIdx = bytes.Index(substituted, target)
		}
	}

	return
}

func isRuneFunc(target rune) func(r rune) bool {
	return func(r rune) bool { return target == r }
}

func isQuoted(src []byte) (v bool, endOfQuote int, quote byte, err error) {
	first := rune(src[0])
	switch first {
	case '"', '\'':
		offset := 1
		for {
			cutset := src[offset:]
			endOfQuote = bytes.IndexFunc(cutset, isRuneFunc(first))

			if endOfQuote == -1 {
				err = fmt.Errorf("unterminated quoted value near %q", string(src))
				return
			}
			if len(cutset) == 1 {
				if cutset[endOfQuote] != '\\' {
					break
				}
			} else {
				if cutset[endOfQuote-1] != '\\' {
					break
				}
			}

			// + 1 -> cuts found quote symbol
			offset = endOfQuote + offset + 1
		}
		quote = byte(first)
		endOfQuote = endOfQuote + offset
		v = true
		return
	default:
		return
	}
}

func sanitizeKey(key []byte) (sanitized []byte) {
	sanitized, _ = bytes.CutPrefix(key, []byte(exportPrefix))
	sanitized = bytes.TrimFunc(sanitized, unicode.IsSpace)
	return
}

func getValidLines(buf *bytes.Buffer) (lines [][]byte, err error) {
	for {
		line, err := buf.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		trimmedLine := bytes.TrimFunc(line, unicode.IsSpace)
		if len(trimmedLine) != 0 && trimmedLine[0] != '#' {
			lines = append(lines, trimmedLine)
		}
	}

	return
}
