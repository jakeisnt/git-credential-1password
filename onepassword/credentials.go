package onepassword

// Credentials defines git credentials.
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type response struct {
	Fields []field `json:"fields"`
}

type field struct {
	Name  string `json:"label"`
	Value string `json:"value"`
}
