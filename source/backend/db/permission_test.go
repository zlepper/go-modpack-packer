package db

import (
	"fmt"
	"testing"
)

func TestGetPermissionsDb(t *testing.T) {
	fmt.Println("TestGetPermissionDb")
	db := GetPermissionsDb()
	if db == nil {
		t.Error("DB was not created successfully")
	}
}

func TestDbHasPermissions(t *testing.T) {
	fmt.Println("TestDbHasPermissions")
	db := GetPermissionsDb()
	if len(db.Permissions) < 100 {
		t.Error("Expected at least 100 permissions in permissionStore")
	}
}

func TestPermissionsDB_GetPermission(t *testing.T) {
	fmt.Println("TestPermissionsDb_GEtPermission")
	db := GetPermissionsDb()
	permissions := db.GetPermission(db.Permissions[0].Modids[0])
	if permissions == nil {
		t.Error("Didn't find permission")
	}

	permissions = db.GetPermission("NotTHereModEYEYYEYEYEYEYYE")
	if permissions != nil {
		t.Error("Expected not to find a permissions, somwhow did. ")
		t.Error(*permissions)
	}
}
