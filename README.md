# http-from-tcp

A simple HTTP server implementation in go.
The server implements HTTP/1.1 and is just for educational purposes.

## HTTP message

From [RFC 9112 Section 2.1](https://datatracker.ietf.org/doc/html/rfc9112#name-message-format)

```
start-line CRLF
*( field-line CRLF )
CRLF
[ message-body ]
```

The `start-line` is either a `request-line` or a `status-line`.
The `request-line` has the format `method SP request-target SP HTTP-version`

The `*( field-line CRLF )` can contain zero or more lines of HTTP headers. These are key-value pairs.