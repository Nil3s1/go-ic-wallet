package transit

type TerminalTapInCommand struct {
	CardNo       string
	StartStation string
}

type TerminalTapOutCommand struct {
	CardNo     string
	EndStation string
}
