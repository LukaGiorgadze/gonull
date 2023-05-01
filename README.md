# Go Nullable with Generics

## Go package simplifies nullable fields handling with Go Generics.

Package gonull provides a generic `Nullable` type for handling nullable values in a convenient way.
This is useful when working with databases and JSON, where nullable values are common.
Unlike other nullable libraries, gonull leverages Go's generics feature, enabling it to work seamlessly with any data type, making it more versatile and efficient.

## Advantages
- Use of Go's generics allows for a single implementation that works with any data type.
- Seamless integration with `database/sql` and JSON marshalling/unmarshalling.
- Reduces boilerplate code and improves code readability.


## Usage

```bash
go get https://github.com/lomsa-dev/gonull
```

```go
type User struct {
	Name     null.Nullable[string]
	Age      null.Nullable[int]
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
		err := rows.Scan( &user.Name, &user.Age)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("ID: %d, Name: %v, Age: %v\n", user.Name.Val, user.Age.Val)
	}
    // ...
}
```

Another example

```go
type Person struct {
	Name    string
	Age     int
	Address gonull.Nullable[string]
}

func main() {
	jsonData := []byte(`{"Name":"Alice","Age":30,"Address":null}`)

	var person Person
	err := json.Unmarshal(jsonData, &person)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Unmarshalled Person: %+v\n", person)

	marshalledData, err := json.Marshal(person)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Marshalled JSON: %s\n", string(marshalledData))
}
```
