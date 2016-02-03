package cli

const (
	argumentFlag = iota
	optionFlag
	longOptionFlag
)

const (
	optional = 1
	isArray  = 2
	required = 4
)

type Flag struct {
	kind        int8
	options     int8
	name        string
	description string
}

func (f Flag) isMatched(val string) {

}

func NewFlag() *Flag {
	return &Flag{}
}
