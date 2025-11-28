package utils

import (
 	"strconv"
 	"strings"
)

func ParseInt(s string, def int) int {
 	i, err := strconv.Atoi(s)
 	if err != nil {
 		return def
 	}
 	return i
}

func ParseFloat(s string, def float64) float64 {
 	f, err := strconv.ParseFloat(s, 64)
 	if err != nil {
 		return def
 	}
 	return f
}

func ParseBool(s string, def bool) bool {
 	s = strings.ToLower(strings.TrimSpace(s))
 	if s == "" {
 		return def
 	}
 	if s == "1" || s == "true" || s == "yes" {
 		return true
 	}
 	if s == "0" || s == "false" || s == "no" {
 		return false
 	}
 	return def
}
