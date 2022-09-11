package main

import (
	"bytes"
	"fmt"
	"io"
	"log"

	"gopkg.in/yaml.v3"
)

var sampling = `
kind: User
data:
  username: foo
  password: foo
`

type (
	Root struct {
		Kind string `yaml:"kind"`
	}

	UserData struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	}
	UserKind struct {
		Data UserData `yaml:"data"`
	}
)

func main() {
	var buf bytes.Buffer

	a := Root{}
	buf.WriteString(sampling)

	err := yaml.Unmarshal(buf.Bytes(), &a)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	userKind(&buf)
}

func userKind(r io.Reader) {
	user := UserKind{}
	var buf bytes.Buffer

	buf.ReadFrom(r)

	yaml.Unmarshal(buf.Bytes(), &user)

	fmt.Println(user)
}
