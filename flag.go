package cli

const (
	valueFlag = iota
	optionFlag
	longOptionFlag
)

type Flag struct {
	kind uint
	name string
}

func (f Flag) isMatched(val string) {

}

func NewFlag() *Flag {
	return &Flag{}
}
