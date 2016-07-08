package types

type UserPermission struct {
	ModId          string           `json:"modId,omitempty"`
	LicenseLink    string           `json:"licenseLink"`
	ModLink        string           `json:"modLink"`
	PermissionLink string           `json:"permissionLink"`
	Policy         PermissionPolicy `json:"policy"`
}
