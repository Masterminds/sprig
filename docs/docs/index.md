# Sprig Function Documentation

The Sprig library provides over 70 template functions for Go's template language.

## Basic Functions

- [String Functions](./functions/strings.md): `trim`, `wrap`, `randAlpha`, `plural`, etc.
    - [String List Functions](./functions/string-slice.md): `splitList`, `sortAlpha`, etc.
- [Integer Math Functions](./functions/math.md): `add`, `max`, `mul`, etc.
    - [Integer Slice Functions](./functions/integer-slice.md): `until`, `untilStep`
- [Float Math Functions](./functions/mathf.md): `addf`, `maxf`, `mulf`, etc.
- [Date Functions](./functions/date.md): `now`, `date`, etc.
- [Defaults Functions](./functions/defaults.md): `default`, `empty`, `coalesce`, `fromJson`, `toJson`, `toPrettyJson`, `toRawJson`, `ternary`
- [Encoding Functions](./functions/encoding.md): `b64enc`, `b64dec`, etc.
- [Lists and List Functions](./functions/lists.md): `list`, `first`, `uniq`, etc.
- [Dictionaries and Dict Functions](./functions/dicts.md): `get`, `set`, `dict`, `hasKey`, `pluck`, `dig`, `deepCopy`, etc.
- [Type Conversion Functions](./functions/conversion.md): `atoi`, `int64`, `toString`, etc.
- [Path and Filepath Functions](./functions/paths.md): `base`, `dir`, `ext`, `clean`, `isAbs`, `osBase`, `osDir`, `osExt`, `osClean`, `osIsAbs`
- [Flow Control Functions](./functions/flow-control.md): `fail`

## Advanced Functions

- [UUID Functions](./advanced/uuid.md): `uuidv4`
- [OS Functions](./advanced/os.md): `env`, `expandenv`
- [Version Comparison Functions](./advanced/semver.md): `semver`, `semverCompare`
- [Reflection](./advanced/reflection.md): `typeOf`, `kindIs`, `typeIsLike`, etc.
- [Cryptographic and Security Functions](./advanced/crypto.md): `derivePassword`, `sha256sum`, `genPrivateKey`, etc.
- [Network](./advanced/network.md): `getHostByName`
