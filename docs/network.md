# Network Functions

Sprig network manipulation functions.

## getHostByName

The `getHostByName` receives a domain name and returns the ip address.

```
getHostByName "www.google.com" would return the corresponding ip address of www.google.com
```

## cidrNetmask

The `cidrNetmask` takes in a cidr and returns the netmask.

```
cidrNetmask "1.2.3.4/32" would return 255.255.255.255
```
