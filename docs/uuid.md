# UUID Functions

Sprig can generate UUID v4 universally unique IDs.

```
uuidv4
```

The above returns a new UUID of the v4 (randomly generated) type.

Sprig can also generate UUID v7 unique IDs which are timebased.

```
uuidv7
```

The above returns a new UUID of the v7 (time-ordered) type, embedding the current timestamp (ms precision) with random bits, making it roughly sortable while remaining globally unique.
