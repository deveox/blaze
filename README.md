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

Encoder always has a `Read` scope. Decoder can have `Write`, `Update` or `Create` scopes. Scope specific API can be obtained by defining singletons (each scope use it's own pull, so Blaze requires explicit definition to save your memory if you don't use some) as follows:

```go
// Admin scope encoder (default)
blaze.Marshal(v)

// You can also create a singleton
var AdminEncoder = encoder.Config{
    Scope: scopes.CONTEXT_ADMIN,
}
AdminEncoder.Marshal(v)

// Client scope encoder, you need to create a singleton
var ClientEncoder = encoder.Config{
    Scope: scopes.CONTEXT_CLIENT,
}
ClientEncoder.Marshal(v)
// DB scope encoder, you need to create a singleton
var DBEncoder = encoder.Config{
    Scope: scopes.CONTEXT_DB,
}
DBEncoder.Marshal(v)

// Admin scope WRITE (update and create) decoder (default)
blaze.Unmarshal(data, &v)
// Admin scope UPDATE decoder
blaze.UnmarshalScoped(data, &v, scopes.DECODE_UPDATE)
// Admin scope CREATE decoder
blaze.UnmarshalScoped(data, &v, scopes.DECODE_CREATE)

// Client scope WRITE decoder
var ClientDecoder = decoder.Config{
    Scope: scopes.CONTEXT_CLIENT,
}
ClientDecoder.Unmarshal(data, &v)
// Client scope UPDATE decoder
ClientDecoder.UnmarshalScoped(data, &v, scopes.DECODE_UPDATE)
// Client scope CREATE decoder
ClientDecoder.UnmarshalScoped(data, &v, scopes.DECODE_CREATE)

// DB scope WRITE decoder
var DBDecoder = decoder.Config{
    Scope: scopes.CONTEXT_DB,
}
DBDecoder.Unmarshal(data, &v)
// DB scope UPDATE decoder
DBDecoder.UnmarshalScoped(data, &v, scopes.DECODE_UPDATE)
// DB scope CREATE decoder
DBDecoder.UnmarshalScoped(data, &v, scopes.DECODE_CREATE)
```

To use scopes, you need to define them in your structs using `blaze` tag. Blaze scope tags have a format of `scope:operation1.operation2`. If you omit scope, it will be applied to all scopes. If you omit operation, it will be applied to all operations.

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

### Decoding of 'null'

Standard library deserializes `null` to `nil` for pointers, and ignore the value for non-pointers. Blaze deserializes `null` to the zero value of the type.

Structs are exceptions, because of scopes. Blaze will set zero values only for field that are available in current scope.

### Custom (de)serialization

You can implement `encoder.Marshaler` and `decoder.Unmarshaler` interfaces to provide custom (de)serialization for your types. You can also register custom (de)serializers for built-in types and 3rd-party types using `blaze.RegisterEncoder` and `blaze.RegisterDecoder` (unstable feature).

Blaze interfaces will take precedence over standard library `json.Marshaler` and `json.Unmarshaler`. The order of importance: `blaze` -> `json` -> `encoding.TextMarshaler/Unmarshaler` -> `builtin decoding`.

```go
// Implement custom marshaler
func (*MyStruct) MarshalBlaze(d *encoder.Encoder) error {
    // Optionally you can do different things based on context
    switch d.Context() {
    case scopes.CONTEXT_ADMIN:
        // ...
    case scopes.CONTEXT_CLIENT:
        // ...
    case scopes.CONTEXT_DB:
        // ...
    }
    // Always use 'd' to encode, it will preserve the scope
    // Use this to encode fields reusing the same encoder
    return d.Encode(1)
    // If you need access to bytes, you can use this, but it will create a new encoder
    b, err := d.Marshal(1)
    // do something with b
    e.Write(b)
    return err
    // ...
}

// Implement custom unmarshaler
func (*MyStruct) UnmarshalBlaze(d *decoder.Decoder, data []byte) error {
    // Optionally you can do different things based on context
    switch d.Context() {
    case scopes.CONTEXT_ADMIN:
        // ...
    case scopes.CONTEXT_CLIENT:
    // ...
    case scopes.CONTEXT_DB:
        // ...
    }
    // Optionally you can do different things based on operation
    switch d.Operation() {
    case scopes.DECODE_CREATE:
        // ...
    case scopes.DECODE_UPDATE:
        // ...
        case scopes.DECODE_WRITE:
        // ...
    }
    // If you implement your own unmarshaler, you need to handle changes yourself
    // Always use 'd' to decode, it will preserve the scope
    type OtherStruct struct {
        // ...
    }
    var other OtherStruct
    err := d.Unmarshal(data, &other)
    if err != nil {
        return err
    }
    // ...
    *MyStruct = res
    return nil
}
```

### Partial marshaling
In Blaze you can marshal only a part of the struct. This can be useful when you want to send only a part of the struct to the client. You can implement GraphQL-like queries using this feature. 

There are two types of partial marshaling, which can be used together:
- **Short** - you can specify fields you want to include to the short output using `blaze:"short"`. It's useful when you need to define short version of the struct statically.
- **Fields** - you can provide an array of fields you want to include in the output *(e.g. `[]string{"name","nested.email"}`)*. It's useful when you need to define the fields dynamically. Nested fields are supported and can be accessed using dot notation.

```go

type Nested struct {
    Age int  `blaze:"short"`
    Email string
}

type User struct {
    ID int `blaze:"short"`
    Name string `blaze:"short"`
    Role string
    Nested Nested `blaze:"short"`
}

blaze.MarshalPartial(v, []string{"name", "nested.email"}, false)
// results in {"name":"John", "nested":{"email":"email@gmail.com"}}

blaze.MarshalPartial(v, nil, true)
// results in {"id":1, "name":"John", "nested":{"age":25}}

blaze.MarshalPartial(v, []string{"name", "nested.age"}, true)
// results in {"name":"John", "nested":{"age":25, "email":"email@gmail.com"}}
```

## Non-standard behavior

### Deserialization

- `input` bytes are considered mutable and may be modified during deserialization. This is not the case with the standard library.
- `Complex128` and `Complex64` is not supported.
- `json.Number` is not supported.
- Map keys are not sorted.
- `encoding.TextUnmarshaler` is partially supported.
- Streaming is not yet supported.

### Serialization

- Map keys are not sorted.
- `encoding.TextMarshaler` is partially supported.

## Performance

Deserialization is much faster for payloads without deep nesting. String deserialization is much faster. Worst case scenario performance is around 25% better.

Serialization is 50%-100% faster.

If you use scopes, the more fields you ignore in particular scope, the faster it gets.
