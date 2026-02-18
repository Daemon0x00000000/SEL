# SEL (Simple Expression Language)

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

SEL is a simple expression language inspired by ServiceNow, designed to evaluate boolean expressions against structured data in Go. It compiles expressions to bytecode and runs them on a lightweight stack-based VM for fast, repeated evaluation.

## 📦 Installation

```bash
go get github.com/Daemon0x00000000/sel
```

## 🚀 Quick Start

```go
package main

import (
    "fmt"
    "github.com/Daemon0x00000000/sel"
)

func main() {
    expr := &sel.Expression{}

    err := expr.Parse("status=active^age>18")
    if err != nil {
        panic(err)
    }

    data := map[string]interface{}{
        "status": "active",
        "age":    25,
    }

    result, err := expr.Eval(data)
    if err != nil {
        panic(err)
    }
    fmt.Println(result) // true
}
```

## 📖 Syntax

### Comparison operators

| Operator | Description | Example |
|----------|-------------|---------|
| `=` | Equality | `name=John` |
| `>` | Greater than | `age>25` |
| `<` | Less than | `price<100` |
| `>=` | Greater than or equal | `score>=80` |
| `<=` | Less than or equal | `quantity<=50` |
| `STARTSWITH` | Starts with | `emailSTARTSWITHadmin` |
| `ENDSWITH` | Ends with | `fileENDSWITH.pdf` |
| `CONTAINS` | Contains a substring | `descriptionCONTAINSerror` |
| `IN` | Membership in a list | `statusINactive,pending,review` |

### Negation

Any operator can be negated with the `!` prefix:

```
!=         → not equal
!IN        → not in list
!CONTAINS  → does not contain
!(expr)    → logical NOT of a group
```

### Logical operators

| Operator | Description | Example |
|----------|-------------|---------|
| `^` | AND | `a=1^b=2` |
| `^OR` | OR | `a=1^ORb=2` |
| `^XOR` | Exclusive OR | `a=1^XORb=2` |

### Grouping

Use parentheses to control precedence:

```
(a=1^ORb=2)^c=3
!(status=active^roleINguest)
```

### Values

- **Unquoted:** `field=value`, `statusINactive,pending`
- **Quoted (single quotes):** `field='value with spaces'`, `tagsIN'a,b','c'`
- **Escape sequences in quotes:** `\'`, `\\`, `\n`, `\t`, `\r`

## 💡 Examples

### Simple filter

```go
expr := &sel.Expression{}
expr.Parse("status=active^age>18")

result, _ := expr.Eval(map[string]interface{}{
    "status": "active",
    "age":    25,
}) // true
```

### IN operator

```go
expr.Parse("statusINpending,active,review")

result, _ := expr.Eval(map[string]interface{}{
    "status": "pending",
}) // true
```

### Negation

```go
// Not equal
expr.Parse("status!=closed")

// Not in list
expr.Parse("role!INguest,anonymous")

// NOT group
expr.Parse("!(status=active^role=guest)")
```

### Complex expression

```go
expr.Parse("sys_id=123^OR(roleINadmin,moderator^status=active)")

result, _ := expr.Eval(map[string]interface{}{
    "sys_id": "456",
    "role":   "admin",
    "status": "active",
}) // true
```

### Reusing a parsed expression

`Expression` is designed to be parsed once and evaluated many times:

```go
expr := &sel.Expression{}
expr.Parse("status=active^score>=80")

for _, record := range records {
    match, _ := expr.Eval(record)
    // ...
}
```

## 📐 Architecture

SEL compiles expressions to bytecode and executes them on a stack-based VM.

```
Expression String
    └─> ast.Parse()         — recursive descent parser
        └─> AST.Compile()   — generates bytecode
            └─> VM.Execute() — stack-based execution
                └─> bool
```

```
sel/
├── sel.go                  # Public API — Expression struct
├── internal/
│   ├── ast/                # Parser + AST → bytecode compiler
│   │   ├── ast.go
│   │   ├── parser.go
│   │   ├── nodes.go
│   │   ├── operators.go
│   │   └── types.go
│   └── vm/                 # Stack-based bytecode VM
│       ├── vm.go
│       ├── opcodes.go
│       ├── handlers.go
│       ├── types.go
│       └── utils.go
└── cmd/main.go             # Usage example
```

## 🧪 Tests

```bash
# Run all tests
go test ./internal/... -v

# Benchmarks
go test ./internal/... -bench=. -benchmem

# Coverage
go test ./internal/... -cover
```

## 📋 Roadmap

- [ ] **JIT compilation** — cache and reuse compiled expressions at runtime
- [ ] **Advanced type system** — explicit types, validation at parse time, type inference
- [ ] **Transformations** — UPPER, LOWER, TRIM, arithmetic, date functions
- [ ] **Aggregations** — COUNT, SUM, AVG
- [ ] **Sub-expressions** — nested query support
- [ ] **AOT compilation** — ahead-of-time mode

## 📝 License

MIT — see [LICENSE](LICENSE).

## 🙏 Acknowledgements

Inspired by the ServiceNow query language.

---

Made with ❤️
