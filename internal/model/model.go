package model

//ValidationResult contains the results of a validation
type ValidationResult struct {
	Valid   bool     `json:"valid"`
	Message string   `json:"message"`
	Errors  []string `json:"errors"`
}

//Checklist contains the name and file of a checklist
type Checklist struct {
	Name string
	File string
}

//Sample tracks samples JSON in the pipeline
type Sample struct {
	UUID     string
	Document string
}

//InterrogationResult contains the checklists sample is a candidate for
type InterrogationResult struct {
	Sample              Sample
	CandidateChecklists []Checklist
}

//Curation is a transformation of a sample document content
type Curation struct {
	Characteristic string
	NewValue       string
}

//CurationPlanResult is the result of executing a curation plan
type CurationPlanResult struct {
	CurationPlan CurationPlan
	Sample       Sample
}
