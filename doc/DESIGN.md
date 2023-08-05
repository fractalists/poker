# board 提供
- query 接口
query() *model.board

# ai/human 提供
- interact 回调
func(board *model.Board, type mode.InteractType) model.Action

- initAI 接口
InitInteract(selfIndex int, getBoardInfoFunc func() *Board) func(board *Board, interactType InteractType) Action


# 模型
model.InteractType
- notify
- ask

# Profile
`go tool pprof -http=:8000 poker.pprof`

# Log Level
Debug // for debug
Info // for gameplay and console UI
Warn or above // for error and train mode