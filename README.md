# Go Nullable with Generics

[![PkgGoDev](https://pkg.go.dev/badge/github.com/LukaGiorgadze/gonull)](https://pkg.go.dev/github.com/LukaGiorgadze/gonull) ![go-mod-verify](https://github.com/LukaGiorgadze/gonull/workflows/Go%20mod/badge.svg) ![go-vuln](https://github.com/LukaGiorgadze/gonull/workflows/Security/badge.svg) ![golangci-lint](https://github.com/LukaGiorgadze/gonull/workflows/Linter/badge.svg) [![codecov](https://codecov.io/gh/LukaGiorgadze/gonull/branch/main/graph/badge.svg?token=76089e7b-f137-4459-8eae-4b48007bd0d6)](https://codecov.io/gh/LukaGiorgadze/gonull)

## Go package simplifies nullable fields handling with Go Generics.

`gonull` is a Go package that provides type-safe handling of nullable values using generics. It's designed to work seamlessly with JSON and SQL operations, making it perfect for web services and database interactions.

## Features

- 🎯 Type-safe nullable values using Go generics
- 💡 Omitzero support
- 🔄 Built-in JSON marshaling/unmarshaling
- 📊 SQL database compatibility
- 🔢 Handles numeric values returned as strings by SQL drivers
- ✨ Zero dependencies

## Usage

```bash
go get github.com/LukaGiorgadze/gonull/v2
```

### Example

```go
package main

import (
    "encoding/json"
    "fmt"

    "github.com/LukaGiorgadze/gonull"
)

type MyCustomInt int
type MyCustomFloat32 float32

type Person struct {
    Name     string                           `json:"name"`
    Age      gonull.Nullable[MyCustomInt]     `json:"age"`
    Address  gonull.Nullable[string]          `json:"address"`
    Height   gonull.Nullable[MyCustomFloat32] `json:"height"`
    IsZero   gonull.Nullable[bool]            `json:"is_zero,omitzero"` // This property will be omitted from the output since it's not present in jsonData.
}

func main() {
    jsonData := []byte(`
    {
        "name": "Alice",
        "age": 15,
        "address": null,
        "height": null
    }`)

    var person Person
    json.Unmarshal(jsonData, &person)
    fmt.Printf("Unmarshalled Person: %+v\n", person)

    marshalledData, _ := json.Marshal(person)
    fmt.Printf("Marshalled JSON: %s\n", string(marshalledData))

    // Output:
    // Unmarshalled Person: {Name:Alice Age:15 Address: Height:0 IsZero:false}
    // Marshalled JSON: {"name":"Alice","age":15,"address":null,"height":null}
    // As you see, IsZero is not present in the output, because we used the omitzero tag introduced in go v1.24.
}
```

### Database example

```go
type User struct {
    Name     gonull.Nullable[string]
    Age      gonull.Nullable[int]
}

func main() {
    // ...
    rows, err := db.Query("SELECT id, name, age FROM users")
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    for rows.Next() {
        var user User
        err := rows.Scan(&user.Name, &user.Age)
        if err != nil {
            log.Fatal(err)
        }
        fmt.Printf("ID: %d, Name: %v, Age: %v\n", user.Name.Val, user.Age.Val)
    }
    // ...
}
```

## Contribution

<a href="https://github.com/LukaGiorgadze/gonull/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=LukaGiorgadze/gonull" />
</a>

Made with [contrib.rocks](https://contrib.rocks).