# Network Functions

Sprig network manipulation functions.

## getHostByName

The `getHostByName` receives a domain name, performs a forward DNS lookup using the local resolver, and returns an ip address. If multiple addresses are returned, one is picked at random.

```
getHostByName "www.google.com" would return the corresponding ip address of www.google.com
```

## getHostByAddr

The `getHostByAddr` receives an IP address (as a string), performs a reverse DNS lookup using the local resolver, and returns a hostname. Note the response will be fully-qualified (i.e. includes a final `.`).

```
getHostByAddr "8.8.8.8" would return `dns.google.com.`.
```