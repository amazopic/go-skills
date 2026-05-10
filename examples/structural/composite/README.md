# Composite — Example

## Pattern summary

Composite lets clients treat leaves (`File`) and branches (`Directory`) through a single `Component` interface. The tree is printed recursively; callers never need to distinguish between the two node types.

## Structure

| Type | Role |
|---|---|
| `Component` | Uniform interface for both leaves and branches |
| `Directory` | Composite — holds a `[]Component` slice; `Print` recurses into children |
| `File` | Leaf — `Add` is a no-op; `Child` returns an empty slice |

## Example tree

```
rootDir (Directory)
├── usrDir (Directory)
│   └── B (File)
└── A (File)
```

`rootDir.Print("")` produces:

```
/root
/root/usr
/root/usr/B
/root/A
```

## Key design choice

`File.Add` silently does nothing (consistent with GoF). In stricter APIs, returning an error or panicking makes the misuse visible earlier — see `skills/structural/composite.md` for the narrow-interface alternative.

## Run the tests

```bash
go test -race ./...
```

## Related skill

`skills/structural/composite.md`
