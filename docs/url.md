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
path:     '/api'
query:    'list=false'
opaque:   nil
fragment: 'anchor'
userinfo: 'admin:secret'
```

For more info, check https://golang.org/pkg/net/url/#URL

When the path has a non-default encoding, such as an encoded slash (`%2F`),
the dict also includes `rawPath`. This preserves that encoding when passed to
`urlJoin`; `path` still contains the decoded path.

## urlJoin
Joins map (produced by `urlParse`) to produce URL string

An optional `rawPath` supplies the encoded form of `path`. If it no longer
matches `path` (for example, after editing `path`), the path is escaped normally.

```
urlJoin (dict "fragment" "fragment" "host" "host:80" "path" "/path" "query" "query" "scheme" "http")
```

The above returns the following string:
```
proto://host:80/path?query#fragment
```
