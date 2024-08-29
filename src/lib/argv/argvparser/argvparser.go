package argvparser

type ArgParserState int

const (
	INIT ArgParserState = iota
	CMD
	FLAG
	ARG
)

func (aps ArgParserState) String() string {
	switch aps {
	case INIT:
		return "INIT"
	case CMD:
		return "CMD"
	case FLAG:
		return "FLAG"
	case ARG:
		return "ARG"
	default:
		panic("invalid ArgParserState: " + string(aps))
	}
}
