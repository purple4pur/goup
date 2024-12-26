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
    "slices"
    "strconv"

    envutil "github.com/purple4pur/goup/envutil"
)

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
    env := new(envutil.Env)
    err := env.Init()
    if err != nil {
        panic(err)
    }

    u, _ := url.Parse(env.TargetUrl)
    log.Printf("Forwarding %s -> %s\n", env.ListenPort, env.TargetUrl)

    proxy := httputil.NewSingleHostReverseProxy(u)
    proxy.ModifyResponse = func(r *http.Response) error {
        b, err := io.ReadAll(r.Body)
        if err != nil {
            return err
        }
        defer r.Body.Close()

        if r.Header.Get(env.SensitiveHeader) != "" {
            msg := ""
            msg += fmt.Sprintf(">>> Get response, wanna inspect this...\n")
            msg += fmt.Sprintf("Bytes length: %d\n", len(b))
            if len(b) >= 8 {
                msg += fmt.Sprintf("Dump: (0-7) % X\n", b[:8])
            }
            if len(b) >= 16 {
                msg += fmt.Sprintf("     (8-15) % X\n", b[8:16])
            }
            if len(b) >= 24 {
                msg += fmt.Sprintf("    (16-23) % X\n", b[16:24])
            }
            if len(b) >= 32 {
                msg += fmt.Sprintf("    (24-31) % X\n", b[24:32])
            }
            msg += fmt.Sprintf("\n")

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
                    msg += fmt.Sprintf("Read EOF. finish.\n")
                    shouldLift = false
                    break
                }
                i += l
                pId := getIntFromLittleEndianBytes(ph.Id[:3])
                pLen := getIntFromLittleEndianBytes(ph.Length[:4])
                msg += fmt.Sprintf("Read packet: id=%d, len=%d\n", pId, pLen)

                switch pId {
                case 75:
                    ver := make([]byte, pLen, pLen)
                    l, _ := binary.Decode(b[i:i + pLen], binary.LittleEndian, &ver)
                    i += l
                    msg += fmt.Sprintf("-> (75) protocol: ver=%d\n", getIntFromLittleEndianBytes(ver))
                case 5:
                    id := make([]byte, pLen, pLen)
                    l, _ := binary.Decode(b[i:i + pLen], binary.LittleEndian, &id)
                    i += l
                    idNum := getIntFromLittleEndianBytes(id)
                    if idNum >= 0x8000_0000 { // negative number
                        idNum = idNum - 0xFFFF_FFFF - 1
                    }
                    msg += fmt.Sprintf("-> (5) reply: id=%d\n", idNum)
                    if slices.Contains(env.AllowedIdList, idNum) {
                        shouldLift = true
                        msg += fmt.Sprintf("-> YES, please lift me up!!!\n")
                    } else {
                        shouldLift = false
                        loop = false
                        msg += fmt.Sprintf("-> NO, I can't lift you...\n")
                    }
                case 71:
                    mode := make([]byte, pLen, pLen)
                    l, _ := binary.Decode(b[i:i + pLen], binary.LittleEndian, &mode)
                    startIdx = i
                    i += l
                    endIdx = i
                    before = getIntFromLittleEndianBytes(mode)
                    msg += fmt.Sprintf("-> (71) client: mode=%d\n", before)
                    loop = false
                default:
                    msg += fmt.Sprintf("-> (%d) don't care what it is. stop here.\n", pId)
                    shouldLift = false
                    loop = false
                }
                msg += fmt.Sprintf("\n")
            }

            if shouldLift {
                msg += fmt.Sprintf("######## Lift U UP!!! ########\n")
                msg += fmt.Sprintf("Before: mode=%d (% X)\n", before, getLittleEndianBytesFromInt(uint32(before)))
                if (before & 0x1) == 1 {
                    after = before | (1 << 2)
                    as := getLittleEndianBytesFromInt(uint32(after))
                    msg += fmt.Sprintf("After:  mode=%d (% X)\n", after, as)
                    if len(as) == (endIdx - startIdx) { // final check
                        b = slices.Replace(b, startIdx, endIdx, as...)
                        msg += fmt.Sprintf("-> here you go! Dump: (%d-%d) % X\n", startIdx, endIdx - 1, b[startIdx:endIdx])
                    } else {
                        msg += fmt.Sprintf("-> something wrong here... I can't lift you.\n")
                    }
                } else {
                    msg += fmt.Sprintf("-> wait aren't you a player? I can't lift you.\n")
                }
            }

            msg += fmt.Sprintf("<<< DONE\n\n\n\n")
            log.Print("\n" + msg)
        }

        // make sure we set the body, and the relevant headers for well-formed clients to respect
        r.Body = io.NopCloser(bytes.NewReader(b))
        r.ContentLength = int64(len(b))
        r.Header.Set("Content-Length", strconv.Itoa(len(b)))
        return nil
    }

    log.Fatal(http.ListenAndServe(env.ListenPort, proxy))
}
