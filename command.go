package cli

type Command struct {
	Name        string
	Version     string
	Description string
	Author      string
	Signature   string
	Action      func()
}
