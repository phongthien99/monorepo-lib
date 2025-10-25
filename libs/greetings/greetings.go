package greetings

import "fmt"

// Hello returns a greeting message for the given name
func Hello(name string) string {
	if name == "" {
		name = "World"
	}
	return fmt.Sprintf("Hello, %s!", name)
}

// Goodbye returns a goodbye message for the given name
func Goodbye(name string) string {
	if name == "" {
		name = "World"
	}
	return fmt.Sprintf("Goodbye, %s!", name)
}

// Welcome returns a welcome message for multiple names
func Welcome(names ...string) string {
	if len(names) == 0 {
		return "Welcome, everyone!"
	}

	message := "Welcome, "
	for i, name := range names {
		if i > 0 {
			if i == len(names)-1 {
				message += " and "
			} else {
				message += ", "
			}
		}
		message += name
	}
	return message + "!"
}
