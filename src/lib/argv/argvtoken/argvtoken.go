package argvtoken

type ArgvTokenKind int

const (
	CMD = iota
	FLAG
	ARG
	END
)

func (a ArgvTokenKind) String() string {
	switch a {
	case CMD:
		return "CMD"
	case FLAG:
		return "FLAG"
	case ARG:
		return "ARG"
	case END:
		return "END"
	default:
		panic("invalid ArgvToken: " + string(rune(a)))
	}
}

type ArgvToken struct {
	Kind ArgvTokenKind
	Val  string
}
