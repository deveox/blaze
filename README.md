# Blaze - JSON serializer and deserializer

The idea behind Blaze is to provide a JSON (de)serializer which can reduce amount of `structs` needed for different API endpoints and contexts. It's designed to be used in REST API and databases that support JSON data types.

It's not a drop-in replacement for standard library, though migration shouldn't be hard. Perfomance is generally better around 20%-100% for deserialization and 10%-25% for serialization. Blaze is also more memory efficient.

## Features

### Scopes

Blaze has 3 scopes:

- **Admin** - for admin API endpoints
- **Client** - for client API endpoints
- **DB** - for database

`Admin` and `Client` scopes can be further configured for specific operations:

- **Read** - for reading data. Tag: `read`;
- **Write** - for writing data. Tag: `write`;
  - **Update** - for updating data (e.g. PATCH handlers). Tag: `update`;
  - **Create** - for creating data (e.g. POST handlers). Tag: `create`;
- **Ignore** - to exclude data from serialization and deserialization. Tag: `-`;

You can combine operations in a single tag using `.` as a separator, e.g. `blaze:"admin:read.create"`.

Encoder always has a `Read` scope. Decoder can have `Write`, `Update` or `Create` scopes. Scope specific scopes can be obtained as follows:

```go
// Admin scope encoder
blaze.Marshal(v)
blaze.AdminEncoder.Marshal(v)
blaze.MarshalScoped(v, scopes.CONTEXT_ADMIN)
// Client scope encoder
blaze.ClientEncoder.Marshal(v)
blaze.MarshalScoped(v, scopes.CONTEXT_CLIENT)
// DB scope encoder
blaze.DBEncoder.Marshal(v)
blaze.MarshalScoped(v, scopes.CONTEXT_DB)

// Admin scope WRITE (update and create) decoder
blaze.Unmarshal(data, &v)
blaze.AdminDecoder.Unmarshal(data, &v)
// Admin scope UPDATE decoder
blaze.AdminDecoder.UnmarshalScoped(data, &v, scopes.DECODE_UPDATE)
// Admin scope CREATE decoder
blaze.AdminDecoder.UnmarshalScoped(data, &v, scopes.DECODE_CREATE)

// Client scope WRITE decoder
blaze.ClientDecoder.Unmarshal(data, &v)
// Client scope UPDATE decoder
blaze.ClientDecoder.UnmarshalScoped(data, &v, scopes.DECODE_UPDATE)
// Client scope CREATE decoder
blaze.ClientDecoder.UnmarshalScoped(data, &v, scopes.DECODE_CREATE)

// DB scope WRITE decoder
blaze.DBDecoder.Unmarshal(data, &v)
// DB scope UPDATE decoder
blaze.DBDecoder.UnmarshalScoped(data, &v, scopes.DECODE_UPDATE)
// DB scope CREATE decoder
blaze.DBDecoder.UnmarshalScoped(data, &v, scopes.DECODE_CREATE)
```

To use scopes, you need to define them in your structs using `blaze` tag. Blaze tag is a comma separated list of scopes and operations. If you omit scope, it will be applied to all scopes. If you omit operation, it will be applied to all operations.

```go
// Define field scopes using blaze tag
type User struct {
    // Admin scope can only read (encode) this field.
    // For client scope this field is not available at all.
    ID int `blaze:"admin:read,client:-"`
    // Admin scope doesn't have any restrictions (can read and write).
    // For client scope this field is available for reading and creating.
    Name string `blaze:"client:read.create"`
    // Both admin and client scopes can only read this field.
    Role string `blaze:"read"`
    // This field is ignored for DB scope. Admin and client scopes can read and write.
    MySecret string `blaze:"no-db"`
    // This field is ignored for all scopes.
    Ignored string `blaze:"-"` // `json:"-"` will also work
}
```

### Unmarshal with changes

Standard library deserialization will overwrite existing struct values only if the field is present in the input. Blaze does the same, but also can optionally provide you with `[]string` of changed fields. This can be useful for implementing `PATCH` requests, where you want to update only the fields that are present in the input.

```go
type UserRole struct {
    Name string
    Role string
}

type User struct {
    ID int
    Name string
    Role UserRole
    Field2 string
    Field3 string
}
data := []byte(`{"name":"John","role":{"name":"John"}, "field2":"value2"}`)
v := &User{}
changes, err := blaze.UnmarshalWithChanges(data, &v)
// v will be {Name: "John", Role: {Name: "John"}, Field2: "value2"}
// changes will be ["name", "role", "role.name", "field2"]
```

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
