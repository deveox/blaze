# Comparison with standard library

## Breaking changes

- `null` deserialization is different. The standard library deserializes `null` to `nil` for pointers, and ignore the value for non-pointers. Blaze deserializes `null` to the zero value of the type.
- `input` bytes are considered mutable and may be modified during deserialization. This is not the case with the standard library.
- `Complex128` and `Complex64` is not supported.
- `json.Number` is not supported.
- Map keys are not sorted during deserialization.
- `encoding.TextUnmarshaler` is not supported. In std lib, it has lower priority than `json.Unmarshaler`, the only exception is map keys, where `encoding.TextUnmarshaler` has higher priority.
- Streaming is not yet supported.

## Changes

- JSON input is not strictly validated at the beginning of the deserialization. Due to lazy deserialization, the input is validated during the deserialization process. In most cases this should not be an issue.

## New features

- Lazy deserialization. The input is not fully deserialized at the beginning, instead the top level value is scanned.
- Iterator for arrays and objects. This allows for deserializing large arrays and objects without loading the entire input into memory.

## Recomendations

- `json.Unmarshaler` is not recomended for use. While it's supported, it misses some of the Blaze features.
