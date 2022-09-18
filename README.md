# httpwut

CLI tool get HTTP status code information.

_This is a tool for me to learn Go, so you probably shouldn't use this._

```
λ httpwut -h        
Get HTTP status code information
usage: httpwut [100..511] [-v|--verbose] [-c|--cats] [-d|--dogs]

λ httpwut 502       
502 - Bad Gateway

λ httpwut 502 --verbose
502 - Bad Gateway
The server, while acting as a gateway or proxy, received an invalid response from an inbound server it accessed while attempting to fulfill the request.
https://www.rfc-editor.org/rfc/rfc9110.html#name-502-bad-gateway
```

Passing `--cats` and/or `--dogs` (you can do both at the time same) will open
the status codes on https://http.cats and https://httpstatusdogs.com,
respectively.