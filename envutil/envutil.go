package envutil

import (
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

func (e *Env) _getEnv(envEntry string) string {
    res, _ := os.LookupEnv(envEntry)
    return res
}

func (e *Env) _getListenPort(envEntry string) string {
    if res := e._getEnv(envEntry); res != "" {
        return ":" + res
    } else {
        return ""
    }
}

func (e *Env) _getAllowedIdList(envEntry string) []int {
    a := e._getEnv(envEntry)
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

func (e Env) _isValid() bool {
    if e.ListenPort == "" || e.TargetUrl == "" {
        return false
    }
    if e.SensitiveHeader == "" {
        return false
    }
    return true
}

func (e *Env) Init() {
    e.ListenPort = e._getListenPort("PORT")
    e.TargetUrl = e._getEnv("TARGET")
    e.SensitiveHeader = e._getEnv("SENSITIVE_HEADER")
    e.AllowedIdList = e._getAllowedIdList("ALLOWED_ID")

    msg := fmt.Sprintf("[Env/Init] --------------------------------\n")
    msg += fmt.Sprintf("Read and parsed from env:\n")
    msg += fmt.Sprintf("  PORT             -> ListenPort      = %s\n", e.ListenPort)
    msg += fmt.Sprintf("  TARGET           -> TargetUrl       = %s\n", e.TargetUrl)
    msg += fmt.Sprintf("  SENSITIVE_HEADER -> SensitiveHeader = %s\n", e.SensitiveHeader)
    msg += fmt.Sprintf("  ALLOWED_ID       -> AllowedIdList   = %+v\n", e.AllowedIdList)
    if e._isValid() {
        log.Print(msg)
    } else {
        msg += fmt.Sprintf("error: Config not valid!\n")
        log.Fatal(msg)
    }
}
