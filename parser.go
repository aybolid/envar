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
		value, err := extractValue(rest)
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
	rest = bytes.TrimFunc(src[endOfKey+1:], unicode.IsSpace)
	return
}

func extractValue(src []byte) (value string, err error) {
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

	value = string(src)
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
			if cutset[endOfQuote-1] != '\\' {
				break
			}

			// + 1 -> cuts found quote symbol?
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
	sanitized, found := bytes.CutPrefix(key, []byte(exportPrefix))
	if found {
		sanitized = bytes.TrimLeftFunc(sanitized, unicode.IsSpace)
	}
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
