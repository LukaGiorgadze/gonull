package examples

import (
	"encoding/json"
	"fmt"

	"github.com/LukaGiorgadze/gonull"
)

type MyCustomInt int
type MyCustomFloat32 float32

type Person struct {
	Name    string                           `json:"name"`
	Age     gonull.Nullable[MyCustomInt]     `json:"age"`
	Address gonull.Nullable[string]          `json:"address"`
	Height  gonull.Nullable[MyCustomFloat32] `json:"height"`
}

func Example() {
	jsonData := []byte(`{"name":"Alice","age":15,"address":null,"height":null}`)

	var person Person
	if err := json.Unmarshal(jsonData, &person); err != nil {
		panic(err)
	}

	// Age is present and valid.
	fmt.Printf("Person.Age is valid: %t, present: %t\n", person.Age.Valid, person.Age.Present)

	// Address is present but invalid (explicit null).
	fmt.Printf("Person.Address is valid: %t, present: %t\n", person.Address.Valid, person.Address.Present)

	// Same for the height.
	fmt.Printf("Person.Height is valid: %t, present: %t\n", person.Height.Valid, person.Height.Present)

	marshalledData, err := json.Marshal(person)
	if err != nil {
		panic(err)
	}
	// Null values will be kept when marshalling to JSON.
	fmt.Println(string(marshalledData))

	// Output:
	// Person.Age is valid: true, present: true
	// Person.Address is valid: false, present: true
	// Person.Height is valid: false, present: true
	// {"name":"Alice","age":15,"address":null,"height":null}
}
