# Slim-Sprig Function Documentation

The Slim-Sprig library provides over 70 template functions for Go's template language.

- [String Functions](strings.md): `trim`, `randAlpha`, `plural`, etc.
  - [String List Functions](string_slice.md): `splitList`, `sortAlpha`, etc.
- [Integer Math Functions](math.md): `add`, `max`, `mul`, etc.
  - [Integer Slice Functions](integer_slice.md): `until`, `untilStep`
- [Date Functions](date.md): `now`, `date`, etc.
- [Defaults Functions](defaults.md): `default`, `empty`, `coalesce`, `fromJson`, `toJson`, `toPrettyJson`, `toRawJson`, `ternary`
- [Encoding Functions](encoding.md): `b64enc`, `b64dec`, etc.
- [Lists and List Functions](lists.md): `list`, `first`, `uniq`, etc.
- [Dictionaries and Dict Functions](dicts.md): `get`, `set`, `dict`, `hasKey`, `pluck`, `dig`, etc.
- [Type Conversion Functions](conversion.md): `atoi`, `int64`, `toString`, etc.
- [Path and Filepath Functions](paths.md): `base`, `dir`, `ext`, `clean`, `isAbs`, `osBase`, `osDir`, `osExt`, `osClean`, `osIsAbs`
- [Flow Control Functions](flow_control.md): `fail`
- Advanced Functions
  - [OS Functions](os.md): `env`, `expandenv`
  - [Reflection](reflection.md): `typeOf`, `kindIs`, `typeIsLike`, etc.
  - [Cryptographic and Security Functions](crypto.md): `sha256sum`, etc.
