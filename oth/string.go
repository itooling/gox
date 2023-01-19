package oth

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/duke-git/lancet/v2/strutil"
	"github.com/google/uuid"
	"github.com/matoous/go-nanoid/v2"
)

func UUID() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}

func NanoID(len ...int) (string, error) {
	return gonanoid.New(len...)
}

func UnicodeToString(u string) string {
	al := strings.Split(u, "\\u")
	s := ""
	for _, v := range al {
		i, _ := strconv.ParseInt(v, 16, 0)
		s += fmt.Sprintf("%c", i)
	}
	return s
}

func ProjectPath() string {
	output, err := exec.Command("go", "env", "GOMOD").Output()
	path := "."
	if err == nil {
		path = filepath.Dir(string(output))
	}
	return path
}

func ReverseString(s string) string {
	return strutil.Reverse(s)
}
