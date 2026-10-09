# Network Functions

Sprig network manipulation functions.

## getHostByName

The `getHostByName` receives a domain name and returns the ip address.

```
getHostByName "www.google.com" would return the corresponding ip address of www.google.com
```

## joinHostPort

The `joinHostPort` combines a host and port into a network address format. When the host contains a colon (such as an IPv6 literal), it encloses the host in square brackets.

```
joinHostPort "127.0.0.1" 80
```

The above returns:
```
127.0.0.1:80
```

For IPv6 addresses:
```
joinHostPort "::1" 8080
```

The above returns:
```
[::1]:8080
```

