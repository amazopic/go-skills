## Registry

The Registry pattern provides a goroutine-safe, typed lookup table for named values.
It is commonly used to register plugins, codecs, drivers, or shared singletons by name,
decoupling producers from consumers without a direct import dependency.

### Files

| File | Purpose |
|---|---|
| `registry.go` | Generic `Registry[T]` with `Register`, `Get`, `MustRegister`, `MustGet`, `Names`, `Len` |
| `registry_test.go` | Table-driven tests; concurrent Register/Get exercised under `-race` |

### Usage

```go
var codecs Registry[*Codec]

codecs.MustRegister("json", &Codec{Name: "json"})
codecs.MustRegister("msgpack", &Codec{Name: "msgpack"})

c, ok := codecs.Get("json")
// c = &Codec{Name: "json"}, ok = true
```

### Running the tests

```bash
go test -race ./behavioral/registry/
```

## -~- THE END -~-
