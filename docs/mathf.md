# Float Math Functions

All math functions operate on `float64` values.

## maxf

Return the largest of a series of floats:

This will return `3`:

```
maxf 1 2.5 3
```

## minf

Return the smallest of a series of floats.

This will return `1.5`:

```
min 1.5 2 3
```

## floor

Returns the greatest float value less than or equal to input value

`floor 123.9999` will return `123.0`

## ceil

Returns the greatest float value greater than or equal to input value

`ceil 123.001` will return `124.0`

## round

Returns a float value with the remainder rounded to the given number to digits after the decimal point.

`round 123.555555` will return `123.556`
