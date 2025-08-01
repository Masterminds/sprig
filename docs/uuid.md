# UUID Functions

## UUID v4

Sprig can generate UUID v4 universally unique IDs.

```
uuidv4
```

The above returns a new UUID of the v4 (randomly generated) type.

## UUID v5

Sprig can generate UUID v5 which can output deterministic IDs for a given namespace and name.

```
uuidv5 "2be4f575-0625-4376-bfca-fc237ac4fd8a" "Hello World"
```

The above returns a new deterministic UUID that will always be the same as long as the parameters stay identical.
The output here would be `c493a152-08ba-5679-a958-de98ebcc8160`
