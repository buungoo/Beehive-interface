package recipes

// Represent a recipe
type Recipe struct {
	Name			string			`json:"name"`
	Ingredients		[]Ingredients	`json:"ingredients"`
}

// Represent induvidual ingredients
type Ingredient struct{
	Name string `json:"name"`
}