package types

type PermissionPolicy string

const (
	Open    PermissionPolicy = "Open"
	Notify  PermissionPolicy = "Notify"
	Request PermissionPolicy = "Request"
	Unknown PermissionPolicy = "Unknown"
	FTB     PermissionPolicy = "FTB"
	Closed  PermissionPolicy = "Closed"
)
