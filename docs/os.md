# OS Functions

_WARNING:_ These functions can lead to information leakage if not used
appropriately.

_WARNING:_ Some notable implementations of Sprig (such as
[Kubernetes Helm](http://helm.sh) _do not provide these functions for security
reasons_.

## env

The `env` function reads an environment variable:

```
env "HOME"
```

## readDir

The `readDir` function reads the named directory, returning all its directory entries sorted by filename:

```
readDir "."
```

## expandenv

To substitute environment variables in a string, use `expandenv`:

```
expandenv "Your path is set to $PATH"
```
