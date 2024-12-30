package envutil

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Env struct {
	ListenPort      string
	TargetUrl       string
	SensitiveHeader string
	AllowedIdList   []int
}

func (e *Env) readEnv(envEntry string) string {
	res, _ := os.LookupEnv(envEntry)
	return res
}

func (e *Env) readListenPort(envEntry string) string {
	if res := e.readEnv(envEntry); res != "" {
		return ":" + res
	} else {
		return ""
	}
}

func (e *Env) readAllowedIdList(envEntry string) []int {
	a := e.readEnv(envEntry)
	if a == "" {
		return []int{}
	}
	aList := strings.Split(a, ",")
	res := make([]int, len(aList), len(aList))
	for i, s := range aList {
		n, err := strconv.Atoi(s)
		if err != nil {
			log.Fatal(err)
		}
		res[i] = n
	}
	return res
}

func (e Env) validate() error {
	if e.ListenPort == "" || e.TargetUrl == "" {
		return errors.New("envutil.Env: no ListenPort or TargetUrl")
	}
	if e.SensitiveHeader == "" {
		return errors.New("envutil.Env: no SensitiveHeader")
	}
	return nil
}

func (e Env) SprintLn() string {
	msg := "{\n"
	msg += fmt.Sprintf("  PORT             -> ListenPort      = %s\n", e.ListenPort)
	msg += fmt.Sprintf("  TARGET           -> TargetUrl       = %s\n", e.TargetUrl)
	msg += fmt.Sprintf("  SENSITIVE_HEADER -> SensitiveHeader = %s\n", e.SensitiveHeader)
	msg += fmt.Sprintf("  ALLOWED_ID       -> AllowedIdList   = %+v\n", e.AllowedIdList)
	msg += "}\n"
	return msg
}

func (e *Env) Init() {
	e.ListenPort = e.readListenPort("PORT")
	e.TargetUrl = e.readEnv("TARGET")
	e.SensitiveHeader = e.readEnv("SENSITIVE_HEADER")
	e.AllowedIdList = e.readAllowedIdList("ALLOWED_ID")
	err := e.validate()
	if err != nil {
		log.Fatalf("error: %s\n%s", err, e.SprintLn())
	}
}

func NewEnv() *Env {
	e := new(Env)
	e.Init()
	return e
}
