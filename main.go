package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"slices"
	"strconv"

	"github.com/purple4pur/goup/envutil"
	"github.com/purple4pur/goup/packets"
)

var env = envutil.NewEnv()

func ModifyResponse(r *http.Response) error {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	if r.Header.Get(env.SensitiveHeader) != "" {
		msg := "[main/ModifyResponse] --------------------------------\n"

		unpacker := packets.NewUnpackerFromBytes(b)
		unpacker.UnpackAll()
		packer := packets.NewPacker()

		shouldLift := false

		for _, p := range unpacker.GetData() {
			pt, err := p.Decode()

			// for unknown packettype, directly append to packer
			if err != nil {
				packer.Append(p)
				continue
			}

			switch pt.GetPacketType() {
			case 5:
				pt5 := pt.(*packets.PacketType5)
				msg += pt5.SprintLn()
				if slices.Contains(env.AllowedIdList, pt5.Id) {
					shouldLift = true
					msg += "-> allowed ID\n"
				} else {
					msg += "-> not allowed ID\n"
				}
			case 71:
				pt71 := pt.(*packets.PacketType71)
				if shouldLift {
					pt71.GoUp()
					msg += pt71.SprintLn()
					msg += "-> lifted\n"
				}
			}
			packer.Append(packets.NewPacketFromPacketTyper(pt))
		}

		log.Print(msg)

		packer.PackAll()
		b = packer.GetData().GetData()
	}

	// make sure we set the body, and the relevant headers for well-formed clients to respect
	r.Body = io.NopCloser(bytes.NewReader(b))
	r.ContentLength = int64(len(b))
	r.Header.Set("Content-Length", strconv.Itoa(len(b)))
	return nil
}

func main() {
	log.Printf("[Env] --------------------------------\n%s", env.SprintLn())

	u, _ := url.Parse(env.TargetUrl)
	log.Printf("[main] --------------------------------\nForwarding %s -> %s\n", env.ListenPort, env.TargetUrl)

	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.ModifyResponse = ModifyResponse

	log.Fatal(http.ListenAndServe(env.ListenPort, proxy))
}
