package constants

import "errors"

type Sex int

const (
	_ Sex = iota
	Male
	Female
	NonBinary
)

func (s Sex) String() string {
	names := [...]string{"Male", "Female", "NonBinary"}
	if s < 0 || int(s) >= len(names) {
		return "Unknown"
	}
	return names[s]
}

func ParseSex(s string) (Sex, error) {
	switch s {
	case "Male":
		return Male, nil
	case "Female":
		return Female, nil
	case "NonBinary":
		return NonBinary, nil

	default:
		return 0, errors.New("invalid sex value")
	}
}
