package onepassword

// Credentials defines git credentials.
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type response struct {
	Details struct {
		Fields []field `json:"fields"`
	} `json:"details"`
}

type response2 struct {
	Fields []field2 `json:"fields"`
}

type login struct {
	Fields []field `json:"fields"`
	Notes  string  `json:"notesPlain"`
}

type field struct {
	Name        string `json:"name"`
	Value       string `json:"value"`
	Type        string `json:"type"`
	Designation string `json:"designation"`
}

type field2 struct {
	Name  string `json:"label"`
	Value string `json:"value"`
}
