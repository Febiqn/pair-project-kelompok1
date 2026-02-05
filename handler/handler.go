package handler

import "github.com/manifoldco/promptui"

func RoleMenu() string {
	prompt := promptui.Select{
		Label: "Choose Role",
		Items: []string{"User", "Admin", "Exit"},
	}
	_, result, _ := prompt.Run()
	return result
}

func ShowUserMenu() string {
	prompt := promptui.Select{
		Label: "User Menu",
		Items: []string{
			"Register Membership",
			"Check Membership",
			"Rent PlayStation",
			"Check Time Left",
			"Pay Bill",
			"Exit",
		},
	}
	_, result, _ := prompt.Run()
	return result
}

func ShowAdminMenu() string {
	prompt := promptui.Select{
		Label: "Admin Menu",
		Items: []string{
			"Update User Membership",
			"View Revenue",
			"Report Broken PS",
			"View PS Condition",
			"Exit",
		},
	}
	_, result, _ := prompt.Run()
	return result
}
