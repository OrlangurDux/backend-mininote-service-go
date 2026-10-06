package models

// RenderTemplate -> template for rendering
type RenderTemplate struct {
	Subject string
	Body    string
}

// ContactData -> contact data
type ContactData struct {
	Name    string
	Phone   string
	Message string
}

// RecoveryPasswordData -> recovery password data
type RecoveryPasswordData struct {
	Name        string
	RecoveryURL string
}
