package model

type Interact interface {
	InitInteract(selfIndex int, getBoardInfoFunc func() *Board) func(board *Board, interactType InteractType) Action
}

type InteractType string

const InteractTypeAsk InteractType = "ASK"
const InteractTypeNotify InteractType = "NOTIFY"
