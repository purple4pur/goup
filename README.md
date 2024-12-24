# Go-up!

**DISCLAIMER: This project is for study only, and does nothing with to-server requests. Use with caution and at your own risk.**

A http proxy, and LIFT U UP in some cases.

### Run

Run with environment variables (example):

- `PORT=8080` - http proxy listen port
- `TARGET=http://example.com:80` - proxy target
- `SENSITIVE_HEADER=xxx` - find this header then do something
- `ALLOWED_ID=123456,8888888` - who can be lifted

### Thanks

- [[stn1slv/http-proxy-logger]](https://github.com/stn1slv/http-proxy-logger)
- [[wtsi-hgi/http-proxy-logger]](https://github.com/wtsi-hgi/http-proxy-logger)
- [[Modifying the response body of an httputil.ReverseProxy response]](https://www.jvt.me/posts/2024/06/25/modify-go-reverseproxy-response/)
- A private server implement of the game
