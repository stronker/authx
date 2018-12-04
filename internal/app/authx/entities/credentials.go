package entities

// BasicCredentialsData is the struct that is store in the database.
type BasicCredentialsData struct {
	// Username is the credential id.
	Username       string
	// Password is the user defined password.
	Password       [] byte
	// RoleID is the assigned role.
	RoleID         string
	// OrganizationID is the assigned organization.
	OrganizationID string
}

// NewBasicCredentialsData creates an instance of BasicCredentialsData.
func NewBasicCredentialsData(username string, password [] byte, roleID string, organizationID string) *BasicCredentialsData {
	return &BasicCredentialsData{
		Username:       username,
		Password:       password,
		RoleID:         roleID,
		OrganizationID: organizationID,
	}
}

// EditBasicCredentialsData is an object that allows to edit the credetentials record.
type EditBasicCredentialsData struct {
	Password *[] byte
	RoleID   *string
}

// WithPassword allows to change the password.
func (d *EditBasicCredentialsData) WithPassword(password [] byte) *EditBasicCredentialsData {
	d.Password = &password
	return d
}

// WithRoleID allows to change the roleID
func (d *EditBasicCredentialsData) WithRoleID(roleID string) *EditBasicCredentialsData {
	d.RoleID = &roleID
	return d
}

// NewEditBasicCredentialsData create a new instance of EditBasicCredentialsData.
func NewEditBasicCredentialsData() *EditBasicCredentialsData {
	return &EditBasicCredentialsData{}
}
