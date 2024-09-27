package main

import (
	"fmt"
	"net/http"
)

func main() {

	fmt.Println("hello world")

	// Create the store and storeHandler
	store := recipes.NewMemStore()
	RecipesHandler := newRecipeHandler(store)
	// Create a new request multiplexer
	// Take incoming request and dispatch them to matching handlers
	mux := http.NewServeMux()

	// Register the routes and handlers
	mux.Handle("/", &homeHandler{})
	mux.Handle("/recipe", &RecipesHandler{})
	mux.Handle("/recipes/", &RecipesHandler{})

	// Run the server
	http.ListenAndServe(":8080", mux)
}

go get -u github.com/gosimple/slug


type recipeStore interface {
	Add(name string, recipe recipes.Recipe) error
	Get(name string) (recipes.Recipe) error
	Update(name string, recipe recipes.Recipe) error
	List() (map[string]recipes.Recipe, error)
	Remove(name string) error
}

type homeHandler struct{}


type RecipesHandler struct{
	store recipeStore
}

func newRecipeHandler(s recipeStore) *RecipesHandler {
	return &RecipesHandler{
		store: s,
	}
}	

func createRecipeHandler(w http.ResponseWriter, r *http.request) {
	// Recipe object that will be populated from JSON payload
	var recipe recipes.Recipe
	if err := json.NewDecoder(r.Body).Decode(&recipe); err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	// Convert the name of the recipe into URL friendly string
	resourceID := slug.Make(recipe.Name)
	// Call the store to add the recipe
	if err := h.store.Add(resourceID, recipe); err != nil {
		InternalServerErrorHandler(w, r)
		return 
	}

	// Set status code 200
	w.WriteHeader(http.statusOk)


}



func InternalServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInterlServerError)
	w.Write([]byte("500 Internal Server Error"))
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write("404 Not Found")

}

func (h *RecipesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPost && RecipeRe.MatchString(r.URL.Path):
		h.CreateRecipe(w,r)
		return
	case r.Method == http.MethodPost && RecipeRe.MatchString(r.URL.Path):
		h.ListRecipe(w,r)
		return
	case r.Method == http.MethodPost && RecipeRe.MatchString(r.URL.Path):
		h.GetRecipe(w,r)
		return
	case r.Method == http.MethodPost && RecipeRe.MatchString(r.URL.Path):
		h.UpdateRecipe(w,r)
		return
	case r.Method == http.MethodPost && RecipeRe.MatchString(r.URL.Path):
		h.DeleteeRecipe(w,r)
		return
	default:
		return
}



func (h *homeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is my home page"))
}

