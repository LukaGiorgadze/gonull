
<div align="right">
  <details>
    <summary >🌐 Language</summary>
    <div>
      <div align="center">
        <a href="https://openaitx.github.io/view.html?user=LukaGiorgadze&project=gonull&lang=en">English</a>
        | <a href="https://openaitx.github.io/view.html?user=LukaGiorgadze&project=gonull&lang=zh-CN">简体中文</a>
        | <a href="https://openaitx.github.io/view.html?user=LukaGiorgadze&project=gonull&lang=zh-TW">繁體中文</a>
        | <a href="https://openaitx.github.io/view.html?user=LukaGiorgadze&project=gonull&lang=ja">日本語</a>
        | <a href="https://openaitx.github.io/view.html?user=LukaGiorgadze&project=gonull&lang=ko">한국어</a>
        | <a href="https://openaitx.github.io/view.html?user=LukaGiorgadze&project=gonull&lang=hi">हिन्दी</a>
        | <a href="https://openaitx.github.io/view.html?user=LukaGiorgadze&project=gonull&lang=th">ไทย</a>
        | <a href="https://openaitx.github.io/view.html?user=LukaGiorgadze&project=gonull&lang=fr">Français</a>
        | <a href="https://openaitx.github.io/view.html?user=LukaGiorgadze&project=gonull&lang=de">Deutsch</a>
        | <a href="https://openaitx.github.io/view.html?user=LukaGiorgadze&project=gonull&lang=es">Español</a>
        | <a href="https://openaitx.github.io/view.html?user=LukaGiorgadze&project=gonull&lang=it">Italiano</a>
        | <a href="https://openaitx.github.io/view.html?user=LukaGiorgadze&project=gonull&lang=ru">Русский</a>
        | <a href="https://openaitx.github.io/view.html?user=LukaGiorgadze&project=gonull&lang=pt">Português</a>
        | <a href="https://openaitx.github.io/view.html?user=LukaGiorgadze&project=gonull&lang=nl">Nederlands</a>
        | <a href="https://openaitx.github.io/view.html?user=LukaGiorgadze&project=gonull&lang=pl">Polski</a>
        | <a href="https://openaitx.github.io/view.html?user=LukaGiorgadze&project=gonull&lang=ar">العربية</a>
        | <a href="https://openaitx.github.io/view.html?user=LukaGiorgadze&project=gonull&lang=fa">فارسی</a>
        | <a href="https://openaitx.github.io/view.html?user=LukaGiorgadze&project=gonull&lang=tr">Türkçe</a>
        | <a href="https://openaitx.github.io/view.html?user=LukaGiorgadze&project=gonull&lang=vi">Tiếng Việt</a>
        | <a href="https://openaitx.github.io/view.html?user=LukaGiorgadze&project=gonull&lang=id">Bahasa Indonesia</a>
      </div>
    </div>
  </details>
</div>

# Go Nullable with Generics

[![PkgGoDev](https://pkg.go.dev/badge/github.com/LukaGiorgadze/gonull)](https://pkg.go.dev/github.com/LukaGiorgadze/gonull) ![go-mod-verify](https://github.com/LukaGiorgadze/gonull/workflows/Go%20mod/badge.svg) ![go-vuln](https://github.com/LukaGiorgadze/gonull/workflows/Security/badge.svg) ![golangci-lint](https://github.com/LukaGiorgadze/gonull/workflows/Linter/badge.svg) [![codecov](https://codecov.io/gh/LukaGiorgadze/gonull/branch/main/graph/badge.svg?token=76089e7b-f137-4459-8eae-4b48007bd0d6)](https://codecov.io/gh/LukaGiorgadze/gonull) [![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/LukaGiorgadze/gonull)

## Go package simplifies nullable fields handling with Go Generics.

`gonull` is a Go package that provides type-safe handling of nullable values using generics. It's designed to work seamlessly with JSON and SQL operations, making it perfect for web services and database interactions.

## Features

- 🎯 One generic `Nullable[T]` works with **any** type
- 💡 Omitzero support
- 🔄 Built-in JSON marshaling/unmarshaling
- 📊 SQL database compatibility
- 🔢 Handles numeric values returned as strings by SQL drivers
- 🧩 Works seamlessly with your own alias or enum types
- ✨ Zero dependencies

### Why gonull?

`Nullable[T]` keeps your code concise by using Go generics for any type. You don't need separate wrappers for strings, ints or custom enumerations. Built-in `sql.Scanner` and `json` support make it easy to integrate with databases and APIs.

```go
type Status string

type Task struct {
    ID    int
    State gonull.Nullable[Status]
}
```

---

## Usage

```bash
go get github.com/LukaGiorgadze/gonull/v2
```

### Example #1

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

### Example #2

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

### Explore More Examples
See [./examples](./examples) directory.

## Contribution

<a href="https://github.com/LukaGiorgadze/gonull/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=LukaGiorgadze/gonull" />
</a>

Made with [contrib.rocks](https://contrib.rocks).
