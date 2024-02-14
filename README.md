# Blaze - JSON serializer and deserializer

The idea behind Blaze is to provide a JSON (de)serializer which can reduce amount of `structs` needed for different API endpoints and contexts. It's designed to be used in REST API and databases that support JSON data types.

It's not a drop-in replacement for standard library, though migration shouldn't be hard. Perfomance is generally better around 20%-100% for deserialization and 10%-25% for serialization. Blaze is also more memory efficient.

## Features

### Scopes

Blaze has 3 scopes:

- **Context** - `HTTP` or `DB`. To define different behavior for your handlers/router (HTTP) and database (e.g. PostgreSQL, MongoDB)
- **User** - `Admin` or `Client`. To define different behavior for different user types. Handy if you want to have same models for your /admin and /public API endpoints.
- **Operation** - `Read`, `Write`, `Update`, `Create`. To define different behavior for different operations. E.g. you may want to allow clients to read some fields, but not to update them. It can help you to dramatically decrease the amount of code in your handlers.

Example:

```go
// Define field scopes using blaze tag
type User struct {
    // This field is only available in admin context, in HTTP scope it's read only, in DB scope - no restrictions
    ID int `blaze:"http:read,client:-"`
    // This field is available as read/create only for clients. No context specific restrictions. No admin restrictions.
    Name string `blaze:"client:read.create"`
    // This field is available as read-only for any user or context scope.
    Role string `blaze:"read"`
}
```

Tag format: `blaze:"scope:operation1.operation2"`, where `scope` is one of `http`, `db`, `admin`, `client` and operation is one of `read`, `write`, `update`, `create`, `-`. Optionally you can omit `scope`: `blaze:"operation1.operation2"` in this case it will be applied to all context and user scopes.

### Unmarshal with changes

Standard library deserialization will overwrite existing struct values only if the field is present in the input. Blaze does the same, but also can optionally provide you with `[]string` of changed fields. This can be useful for implementing `PATCH` requests, where you want to update only the fields that are present in the input.

### Auto Camel Case

Blaze will automatically convert field names to camelCase. If you want to specify a custom name, you can use `json` tag as usual.

Example:

```go
type User struct {
    ID int // will be (de)serialized as "id"
    MyName string  `json:"name"` // will be (de)serialized as "name"
}
```

### Omit empty by default

Blaze will omit empty fields by default to reduce response size. If you want to include empty fields, you can use `keep` tag.

Example:

```go
type User struct {
    ID int // will be omitted if zero
    Name string  `blaze:"keep"` // will be included even if "" (zero value)
}

```

### Custom (de)serialization

Instead of relying on interface implementation, like `json.Marshaller` and `json.Unmarshaller`, Blaze provides a way to register custom (de)serializers for specific types. Which provides better performance, and gives you the ability to set behavior for 3rd-party types, built-ins and even interfaces.

```go
// Register custom marshaller on built-in slice type
blaze.RegisterMarshaler([]string{}, func(e *blaze.Encoder, v []string) ([]byte, error) {
    return []byte(strings.Join(v)), nil
}
// Register custom marshaller for 3rd-party type
blaze.RegisterMarshaler(&time.Time{}, func(e *blaze.Encoder, v time.Time) ([]byte, error) {
    return []byte(v.Format(time.RFC3339)), nil
}
// Register custom unmarshaler for interface
type MyInterface interface {
    Type() string
}

type Type1 struct {
    // ...
}

func (t *Type1) Type() string {
    return "type1"
}

type Type2 struct {
    // ...
}

func (t *Type2) Type() string {
    return "type2"
}

var mi MyInterface
blaze.RegisterUnmarshaler(mi, func(d *blaze.Decoder, data []byte) (MyInterface, error) {
    var res MyInterface
    type MyInterfaceGenericImpl struct {
        Type string
    }
    v := MyInterfaceImpl{}
    err := d.Decode(&v)
    switch v.Type {
    case "type1":
        v2 := &Type1{}
        err = d.Decode(v2)
        // ...
        res = v2
        // ...
    }
    return res
})
```

## Non-standard behavior

### Deserialization

- `null` deserialization is different. The standard library deserializes `null` to `nil` for pointers, and ignore the value for non-pointers. Blaze deserializes `null` to the zero value of the type.
- `input` bytes are considered mutable and may be modified during deserialization. This is not the case with the standard library.
- `Complex128` and `Complex64` is not supported.
- `json.Number` is not supported.
- Map keys are not sorted.
- `encoding.TextUnmarshaler` is not supported.
- `json.Unmarshaler` is not supported by default. But you can add support for them using `blaze.RegisterUnmarshaler` on per type basis, e.g. `blaze.RegisterUnmarshaler(&time.Time{})`
- Streaming is not yet supported.

### Serialization

- `json.Marshaler` is not supported by default. You can add support for it using `blaze.RegisterMarshaler` on per type basis, e.g. `blaze.RegisterMarshaler(&time.Time{})`
- Map keys are not sorted.
- `encoding.TextMarshaler` is not supported.

## Performance

Deserialization is much faster for payloads without deep nesting. String deserialization is much faster. Worst case scenario performance is around 10% better.

Serialization is 10%-25% faster.

And of course, if you use scopes, the more fields you ignore in particular scope, the faster it gets.
