# URL Functions

## urlParse
Parses string for URL and produces dict with URL parts

```
urlParse "http://admin:secret@server.com:8080/api?list=false#anchor"
```

The above returns a dict, containing URL object:
```yaml
scheme:   'http'
host:     'server.com:8080'
hostname: 'server.com'
port:     '8080'
path:     '/api'
query:    'list=false'
opaque:   nil
fragment: 'anchor'
userinfo: 'admin:secret'
```

IPv6 addresses are also supported:
```
urlParse "http://[2001:db8::1]:8080/api"
```

Returns:
```yaml
scheme:   'http'
host:     '[2001:db8::1]:8080'
hostname: '2001:db8::1'
port:     '8080'
path:     '/api'
query:    ''
opaque:   nil
fragment: ''
userinfo: ''
```

For more info, check https://golang.org/pkg/net/url/#URL

## urlJoin
Joins map (produced by `urlParse`) to produce URL string

```
urlJoin (dict "fragment" "fragment" "host" "host:80" "path" "/path" "query" "query" "scheme" "http")
```

The above returns the following string:
```
proto://host:80/path?query#fragment
```

Alternatively, you can use `hostname` and `port` instead of `host`:
```
urlJoin (dict "fragment" "fragment" "hostname" "host" "port" "80" "path" "/path" "query" "query" "scheme" "http")
```

This will also produce:
```
proto://host:80/path?query#fragment
```

IPv6 addresses are automatically detected and wrapped in square brackets when using `hostname` and `port`:
```
urlJoin (dict "hostname" "2001:db8::1" "port" "8080" "path" "/path" "scheme" "http")
```

This produces:
```
http://[2001:db8::1]:8080/path
```
