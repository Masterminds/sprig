# URL Functions

Sprig has a number of URL escaping functions.

## pathEscape

The `pathEscape` function escapes the string so it can be safely placed inside a URL path segment.

```
pathEscape "Hello, World + Sprig!"
```

The above returns `Hello%2C%20World%20+%20Sprig%21z`

## pathUnescape

The `pathUnescape` function does the inverse transformation of `pathEscape`, converting each 3-byte encoded substring of the form "%AB" into the hex-decoded byte 0xAB. It returns an error if any % is not followed by two hexadecimal digits.

```
pathUnscape "Hello%2C%20World%20+%20Sprig%21"
```

The above returns `Hello, World + Sprig!`

## queryEscape

The `queryEscape` function escapes the string so it can be safely placed inside a URL query.

```
queryEscape "Hello, World + Sprig!"
```

The above returns `Hello%2C+World+%2B+Sprig%21`

## queryUnescape

The `queryUnescape` function does the inverse transformation of `queryEscape`, converting each 3-byte encoded substring of the form "%AB" into the hex-decoded byte 0xAB. It returns an error if any % is not followed by two hexadecimal digits.

```
queryUnscape "Hello%2C+World+%2B+Sprig%21"
```

The above returns `Hello, World + Sprig!`

