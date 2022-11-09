package cli

type Command struct {
	Name          string
	Description   string
	Flags         []Flag
	SubCommands   []*Command
	Before        ActionFunc
	Action        ActionFunc
	ActionCommand ActionFunc
}
