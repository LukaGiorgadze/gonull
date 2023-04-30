# Go Nullable

## Go package simplifies nullable fields handling.

This package provides a generic nullable type implementation for use with Go's `database/sql` package.
It simplifies handling nullable fields in SQL databases by wrapping any data type with the `Nullable` type.
The Nullable type works with both basic and custom data types and implements the `sql.Scanner` and `driver.Valuer` interfaces, making it easy to use with the `database/sql` package.

## Use case
In a web application, you may have a user profile with optional fields like name, age, or whatever. These fields can be left empty by the user, and your database stores them as `NULL` values. Using the `Nullable` type from this library, you can easily handle these optional fields when scanning data from the database or inserting new records. By wrapping the data types of these fields with the `Nullable` type, you can handle `NULL` values without additional logic, making your code cleaner and more maintainable.


## Usage

```bash
go get https://github.com/lomsa-dev/gonull
```

```go
type User struct {
	ID       int
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
		err := rows.Scan(&user.ID, &user.Name, &user.Age)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("ID: %d, Name: %v, Age: %v\n", user.ID, user.Name.Val, user.Age.Val)
	}
    // ...
}

```
