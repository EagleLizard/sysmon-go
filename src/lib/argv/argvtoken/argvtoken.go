package argvtoken

type ArgvTokenKind int

const (
	CMD = iota
	FLAG
	ARG
	END
)

type ArgvToken struct {
	Kind ArgvTokenKind
	Val  string
}
