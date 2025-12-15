package state

type State string

const (
	Attached State = "attached"
	Detached State = "detached"
	Visible  State = "visible"
	Hidden   State = "hidden"
)
