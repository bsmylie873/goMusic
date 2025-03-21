package constants

import "errors"

type Title int

const (
	_ Title = iota
	Mr
	Mrs
	Ms
	Dr
	Prof
)

func (t Title) String() string {
	names := [...]string{"Mr", "Mrs", "Ms", "Dr", "Prof"}
	if t < 0 || int(t) >= len(names) {
		return "Unknown"
	}
	return names[t]
}

func ParseTitle(s string) (Title, error) {
	switch s {
	case "Mr":
		return Mr, nil
	case "Mrs":
		return Mrs, nil
	case "Ms":
		return Ms, nil
	case "Dr":
		return Dr, nil
	case "Prof":
		return Prof, nil

	default:
		return 0, errors.New("invalid title value")
	}
}
