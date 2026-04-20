# Network Functions

Sprig network manipulation functions.

## getHostByName

The `getHostByName` receives a domain name and returns the ip address.

```
getHostByName "www.google.com" would return the corresponding ip address of www.google.com
```

## joinHostPort

`joinHostPort` combines a host and port into a network address of the form `host:port`. For IPv6 addresses the host is enclosed in square brackets: `[host]:port`. This is the inverse of `splitHostPort`.

```
joinHostPort "192.0.2.1" "80"    -> "192.0.2.1:80"
joinHostPort "2001:db8::1" "80"  -> "[2001:db8::1]:80"
joinHostPort "example.com" "443" -> "example.com:443"
```

## splitHostPort

`splitHostPort` splits a network address of the form `host:port` or `[host]:port` into its host and port components, returning a map with keys `host` and `port`. This is the inverse of `joinHostPort`.

```
{{- with splitHostPort "192.0.2.1:80" }}
host: {{ .host }}
port: {{ .port }}
{{- end }}
```
