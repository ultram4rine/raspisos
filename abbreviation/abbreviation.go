package abbreviation

import "strings"

//Faculty returns a shortened form of faculty name, if can't returns full name
func Faculty(name string) string {
	parts := strings.Split(name, " ")

	if strings.Contains(parts[0], "-") {
		p := strings.Split(parts[0], "-")
		r1 := []rune(p[0])
		r2 := []rune(p[1])
		return string(r1[0:3]) + "-" + string(r2[0:3])
	} else if parts[len(parts)-1] == "факультет" {
		r := []rune(parts[0])
		if string(r[0:3]) == "Гео" {
			switch string(r[3]) {
			case "г":
				return string(r[0:7])
			case "л":
				return string(r[0:6])
			}
		} else {
			switch strings.ContainsAny(string(r[4]), "аэиуыеёюя") {
			case true:
				return string(r[0:2]) + "фак"
			case false:
				return string(r[0:3]) + "фак"
			}
		}
	} else if parts[len(parts)-1] == "колледж" {
		r := []rune(parts[0])
		return string(r[0:3]) + "кол"
	} else if parts[0] == "Институт" {
		var (
			runes [][]rune
			abbr  string
		)

		for i := range parts {
			r := []rune(parts[i])
			runes = append(runes, r)
		}

		for _, p := range runes {
			if string(p) == "и" {
				abbr += "и"
			} else {
				abbr += strings.ToUpper(string(p[0]))
			}
		}

		return abbr
	}

	return name
}
