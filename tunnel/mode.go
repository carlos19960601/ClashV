package tunnel

type TunnelMode int

var ModeMapping = map[string]TunnelMode{
	Global.String(): Global,
	Rule.String():   Rule,
	Direct.String(): Direct,
}

const (
	Global TunnelMode = iota
	Rule
	Direct
)

func (m TunnelMode) String() string {
	switch m {
	case Global:
		return "global"
	case Rule:
		return "rule"
	case Direct:
		return "direct"
	default:
		return "unknown"
	}
}
