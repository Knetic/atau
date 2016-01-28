package atau

type Resource struct {

	Name string
	Methods map[string]Method `json:"methods"`
}

// methods are not HTTP methods, they can be considered prefixes to the resource name.
// a method "insert" on the resource "coin" would be rendered in generated code as "insertCoin".
