package cli

type Command struct {
	Name        string
	Version     string
	Description string
	Author      string
	Signature   string
	Flags       []*Flag
	Action      func()
}

func NewCommand() {
}
