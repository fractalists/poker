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