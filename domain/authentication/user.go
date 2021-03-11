package authentication

type User struct {
	Sub               string `json:"subject"`           // mandatory
	Email             string `json:"email"`             // mandatory
	FamilyName        string `json:"familyName"`        // mandatory
	GivenName         string `json:"givenName"`         // mandatory
	PreferredUsername string `json:"preferredUsername"` // optional
}
