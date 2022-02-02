package gtsasamiserver

type GuardianType int32

const (
	Warrior GuardianType = 0
	Ranged  GuardianType = 1
	Support GuardianType = 2
	Tanker  GuardianType = 3
)

type Guardian struct {
	ID      string       `json:"id"`
	Name    string       `json:"name"`
	Class   GuardianType `json:"class"`
	Element Element      `json:"element"`
}
