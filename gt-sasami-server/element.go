package gtsasamiserver

type ElementType int32

const (
	Basic ElementType = 0
	Dark  ElementType = 1
	Earth ElementType = 2
	Fire  ElementType = 3
	Light ElementType = 4
	Water ElementType = 5
)

type Element struct {
	ID            string      `json:"id"`
	Name          string      `json:"name"`
	Type          ElementType `json:"type"`
	StrongAgainst ElementType `json:"strongAgainst"`
	WeakAgainst   ElementType `json:"weakAgainst"`
}
