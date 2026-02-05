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
			"View Revenue",
			"Report Broken PS",
			"Exit",
		},
	}
	_, result, _ := prompt.Run()
	return result
}
