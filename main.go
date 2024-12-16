package main

import (
    "bytes"
    "encoding/binary"
    "fmt"
    "io"
    "log"
    "net/http"
    "net/http/httputil"
    "net/url"
    "os"
    "slices"
    "strconv"
    "strings"
)

func getEnv(key string) string {
    if value, ok := os.LookupEnv(key); ok {
        log.Printf("ENV: %s=%s", key, value)
        if value != "" {
            return value
        }
    }
    panic("env not set or empty")
    return "" // unreachable
}

func getTargetUrl() string {
    return getEnv("TARGET")
}

func getListenPort() string {
    return ":" + getEnv("PORT")
}

func getSensitiveHeader() string {
    return getEnv("SENSITIVE_HEADER")
}

func getAllowedIdList() []int {
    a := getEnv("ALLOWED_ID")
    aList := strings.Split(a, ",")
    res := make([]int, len(aList), len(aList))
    for i, s := range aList {
        n, err := strconv.Atoi(s)
        if err != nil {
            panic(err)
        }
        res[i] = n
    }
    return res
}

type pHead struct {
    Id     [3]byte
    Length [4]byte
}

func getIntFromLittleEndianBytes(s []byte) int {
    var res int
    _s := make([]byte, len(s), len(s))
    _ = copy(_s, s)
    slices.Reverse(_s)
    for _, v := range _s {
        res <<= 8
        res |= int(v)
    }
    return res
}

func getLittleEndianBytesFromInt(n uint32) []byte {
    res := make([]byte, 4, 4)
    binary.LittleEndian.PutUint32(res, n)
    return res
}

func main() {
    p := getListenPort()
    t := getTargetUrl()
    h := getSensitiveHeader()
    a := getAllowedIdList()

    u, _ := url.Parse(t)
    log.Printf("Forwarding %s -> %s\n", p, t)

    proxy := httputil.NewSingleHostReverseProxy(u)
    proxy.ModifyResponse = func(r *http.Response) error {
        b, err := io.ReadAll(r.Body)
        if err != nil {
            return err
        }
        defer r.Body.Close()

        if r.Header.Get(h) != "" {
            fmt.Printf("\n----------------------------------------------------------------\n")
            log.Printf("Get response, wanna inspect this...\n")
            log.Printf("\n")
            log.Printf("Bytes length: %d\n", len(b))
            if len(b) >= 8 {
                log.Printf("Dump: (0-7) % X\n", b[:8])
            }
            if len(b) >= 16 {
                log.Printf("     (8-15) % X\n", b[8:16])
            }
            if len(b) >= 24 {
                log.Printf("    (16-23) % X\n", b[16:24])
            }
            if len(b) >= 32 {
                log.Printf("    (24-31) % X\n", b[24:32])
            }
            log.Printf("\n")

            shouldLift := false
            var before int
            var after int
            var startIdx int
            var endIdx int
            var ph pHead
            i := 0

            for loop := true; loop; {
                l, err := binary.Decode(b[i:], binary.LittleEndian, &ph)
                if err != nil {
                    log.Printf("Read EOF. finish.\n")
                    shouldLift = false
                    break
                }
                i += l
                pId := getIntFromLittleEndianBytes(ph.Id[:3])
                pLen := getIntFromLittleEndianBytes(ph.Length[:4])
                log.Printf("Read packet: id=%d, len=%d\n", pId, pLen)

                switch pId {
                case 75:
                    ver := make([]byte, pLen, pLen)
                    l, _ := binary.Decode(b[i:i + pLen], binary.LittleEndian, &ver)
                    i += l
                    log.Printf("-> (75) protocol: ver=%d\n", getIntFromLittleEndianBytes(ver))
                case 5:
                    id := make([]byte, pLen, pLen)
                    l, _ := binary.Decode(b[i:i + pLen], binary.LittleEndian, &id)
                    i += l
                    idNum := getIntFromLittleEndianBytes(id)
                    if idNum >= 0x8000_0000 { // negative number
                        idNum = idNum - 0xFFFF_FFFF - 1
                    }
                    log.Printf("-> (5) reply: id=%d\n", idNum)
                    if slices.Contains(a, idNum) {
                        shouldLift = true
                        log.Printf("-> YES, please lift me up!!!\n")
                    } else {
                        shouldLift = false
                        loop = false
                        log.Printf("-> NO, I can't lift you...\n")
                    }
                case 71:
                    mode := make([]byte, pLen, pLen)
                    l, _ := binary.Decode(b[i:i + pLen], binary.LittleEndian, &mode)
                    startIdx = i
                    i += l
                    endIdx = i
                    before = getIntFromLittleEndianBytes(mode)
                    log.Printf("-> (71) client: mode=%d\n", before)
                    loop = false
                default:
                    log.Printf("-> (%d) don't care what it is. stop here.\n", pId)
                    shouldLift = false
                    loop = false
                }
                log.Printf("\n")
            }

            if shouldLift {
                log.Printf("######## Lift U UP!!! ########\n")
                log.Printf("Before: mode=%d (% X)\n", before, getLittleEndianBytesFromInt(uint32(before)))
                if (before & 0x1) == 1 {
                    after = before | (1 << 2)
                    as := getLittleEndianBytesFromInt(uint32(after))
                    log.Printf("After:  mode=%d (% X)\n", after, as)
                    if len(as) == (endIdx - startIdx) { // final check
                        b = slices.Replace(b, startIdx, endIdx, as...)
                        log.Printf("-> here you go! Dump: (%d-%d) % X\n", startIdx, endIdx - 1, b[startIdx:endIdx])
                    } else {
                        log.Printf("-> something wrong here... I can't lift you.\n")
                    }
                } else {
                    log.Printf("-> wait aren't you a player? I can't lift you.\n")
                }
            }
        }

        // make sure we set the body, and the relevant headers for well-formed clients to respect
        r.Body = io.NopCloser(bytes.NewReader(b))
        r.ContentLength = int64(len(b))
        r.Header.Set("Content-Length", strconv.Itoa(len(b)))
        return nil
    }

    log.Fatal(http.ListenAndServe(p, proxy))
}
