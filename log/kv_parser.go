//go:build !release

// Only used to help tests assert logs. Quick and dirty.

package log

import "strings"

func KvParseAll(l string) []map[string]string {
	lines := strings.Split(l, "\n")
	lookup := make([]map[string]string, len(lines))
	for i, line := range lines {
		lookup[i] = KvParse(line)
	}
	return lookup
}

func KvParse(line string) map[string]string {
	if len(line) == 0 {
		return nil
	}
	s := strings.TrimRight(line, "\n")
	lookup := make(map[string]string)

	for len(s) > 0 {
		if s[0] == ' ' {
			s = s[1:]
		}
		keyEnd := strings.Index(s, "=")
		key := s[:keyEnd]
		s = s[keyEnd+1:]
		if len(s) == 0 {
			lookup[key] = ""
			break
		}

		i := 0
		quoted := s[0] == '"'
		if quoted {
			i = 1
		}
		escape := false
		for ; i < len(s); i++ {
			c := s[i]
			if !quoted {
				if c == ' ' {
					break
				}
				continue
			}

			if c == '\\' {
				escape = true
				continue
			}

			if c == '"' && !escape {
				i += 1 // skip closing quote
				break
			}
			escape = false
		}

		lookup[key] = s[:i]
		s = s[i:]
	}
	return lookup
}
