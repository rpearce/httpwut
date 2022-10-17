# httpwut

CLI tool get HTTP status code information.

_This is a tool for me to learn Go, so you probably shouldn't use this._

```
λ httpwut is -h    
Lookup HTTP status codes

Usage:
  httpwut is [flags]

Flags:
  -c, --cats      Open HTTP Status Cats webpage for status
  -d, --dogs      Open HTTP Status Dogs webpage for status
  -h, --help      help for is
  -v, --verbose   Print status code description and URL

λ httpwut is 502       
502 - Bad Gateway

λ httpwut is 502 --verbose
502 - Bad Gateway
The server, while acting as a gateway or proxy, received an invalid response from an inbound server it accessed while attempting to fulfill the request.
https://www.rfc-editor.org/rfc/rfc9110.html#name-502-bad-gateway
```

Passing `--cats` and/or `--dogs` (you can do both at the time same) will open
the status codes on https://http.cats and https://httpstatusdogs.com,
respectively.