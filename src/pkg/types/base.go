package types

// Base is the base view model that others extend
type Base struct {
	Title string
}

// Get is the base for all GET requests
type Get struct {
	ID string
}
