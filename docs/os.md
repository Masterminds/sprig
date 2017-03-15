# OS Functions

_WARNING:_ These functions can lead to information leakage if not used
appropriately.

## `env`

The `env` function reads an environment variable:

```
env "HOME"
```

## `expandenv`

To substitute environment variables in a string, use `expandenv`:

```
expandenv "Your path is set to $PATH"
```
