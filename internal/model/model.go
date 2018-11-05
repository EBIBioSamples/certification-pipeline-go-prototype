package model

//ValidationResult contains the results of a validation
type ValidationResult struct {
	Valid   bool     `json:"valid"`
	Message string   `json:"message"`
	Errors  []string `json:"errors"`
}

//Sample tracks samples JSON in the pipleline
type Sample struct {
	UUID     string
	Document string
}
