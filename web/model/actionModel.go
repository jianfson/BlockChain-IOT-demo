package model

type Action struct {
	ActionID         int
	UserID           int
	ActionType       string
	ActionTargetName string
	ActionTargetID   string
	ActionTime       string
}
