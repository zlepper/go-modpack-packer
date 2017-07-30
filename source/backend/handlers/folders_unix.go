// +build darwin dragonfly freebsd linux nacl netbsd openbsd solaris

package handlers

// Gets the available drives on the system
func GetDrives() (r []string, err error) {
	// You win this time unix!!
	return getFolderList("/")
}
