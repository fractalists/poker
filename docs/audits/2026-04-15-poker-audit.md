# Repository Audit Report

## Executive Summary

- Repository condition: 可构建、规模不大，但真正决定对局行为的流程代码存在多个高风险错误，而且这些路径几乎没有测试覆盖。
- Highest-risk themes:
  - 训练模式泄露隐藏信息，导致 AI 训练与 profile 结果不可信
  - 双人局行动顺序与加注状态维护存在规则级错误
  - 训练统计并发不安全，可能随机 panic 或污染结果
- Top 3 next actions:
  - 修复 `TrainMode` 下的隐藏信息泄露
  - 修复双人局 post-flop 行动起点与 `LastRaiseAmount` 计算，并补回归测试
  - 去掉训练统计的共享 map 并发写

## Audit Scope and Mode

- Language: 中文
- Mode: `guided-audit`
- Selected dimensions: `correctness`, `verification`, `architecture`, `maintainability`, `ops-config`
- Exclusions: `.idea/`、根目录历史二进制、`generated/log/`、`generated/pprof/`、旧审计目录 `.audit-work/2026-04-15-poker-guided/`
- Deepening used: `yes`，但仅限源码与验证基线层面的定向加深

## Repository Map

- Stack: Go 1.20, Go modules, Fyne, logrus, ants, testify
- Major directories:
  - `config/`
  - `game/`
  - `interact/`
  - `model/`
  - `process/`
  - `util/`
  - `doc/`
  - `generated/`
- Entry points:
  - `main.go`
  - `game/unlimited.PlayPoker`
  - `game/unlimited.Train`
  - `game/colosseum.PlayPoker`
- Test paths:
  - `process/score_test.go`
  - `process/settle_test.go`
- Build and deployment paths:
  - `go test ./...`
  - `go build ./...`
  - 未发现 CI、部署或迁移文件

## Verification Baseline

- Commands run:
  - `git status --short`
  - `git -c safe.directory=D:/Git/go/src/poker status --short`
  - `go test ./...`
  - `New-Item -ItemType Directory -Force '.audit-work\\go-build-cache' | Out-Null; $env:GOCACHE=(Resolve-Path '.audit-work\\go-build-cache').Path; go test ./...`
  - `New-Item -ItemType Directory -Force '.audit-work\\go-build-cache' | Out-Null; $env:GOCACHE=(Resolve-Path '.audit-work\\go-build-cache').Path; go build ./...`
  - `New-Item -ItemType Directory -Force '.audit-work\\go-build-cache' | Out-Null; $env:GOCACHE=(Resolve-Path '.audit-work\\go-build-cache').Path; go test ./... -cover`
  - `New-Item -ItemType Directory -Force '.audit-work\\go-build-cache' | Out-Null; $env:GOCACHE=(Resolve-Path '.audit-work\\go-build-cache').Path; go test ./... '-coverprofile=.audit-work/coverage2.out'`
  - `go tool cover '-func=.audit-work/coverage2.out'`
- Passed:
  - 带仓库内 `GOCACHE` 的 `go test ./...`
  - 带仓库内 `GOCACHE` 的 `go build ./...`
  - 覆盖率命令
- Failed:
  - 初次 `git status --short` 受 `safe.directory` 策略阻断
  - 初次 `go test ./...` 受默认 `GOCACHE` 写权限阻断
- Missing:
  - 交互式对局验证
  - 训练并发验证
  - GUI 运行验证
- Confidence impact:
  - 可构建性结论较强
  - 流程正确性结论主要依赖源码证据，因为现有测试没有跑到这些路径

## Key Findings

### [high] 训练模式把完整牌桌直接暴露给 AI

- `title`: 训练模式把完整牌桌直接暴露给 AI
- `angle`: correctness
- `severity`: high
- `confidence`: high
- `evidence`: `model/board.go:99-101` 在 `config.TrainMode` 下直接返回真实 `board`；`process/process.go:421` 会把这份对象直接交给玩家交互实现。
- `impact`: 训练模式下 AI 可以读取其他玩家手牌和未翻开的公共牌，训练结果、胜率估计与 CPU profile 都失去参考价值。
- `scope`: `model/board.go`, `process/process.go`, `game/unlimited/training.go`
- `recommendation`: 训练模式只应跳过渲染，不应跳过信息隔离；仍然要为每个玩家返回去敏后的牌桌视图。
- `status`: confirmed

### [high] 双人局翻牌后行动顺序错误

- `title`: 双人局翻牌后行动顺序错误
- `angle`: correctness
- `severity`: high
- `confidence`: high
- `evidence`: `process/process.go:125-128` 让双人局 `UTG` 与 `SB` 重合；`process/process.go:247` 让非翻前轮默认从 `actualSmallBlindIndex` 开始，只在 `process/process.go:264` 把翻前起手位改成 `UTG`。
- `impact`: heads-up 进入 flop / turn / river 后仍由小盲先行动，和标准德州规则相反，会直接改变可用动作与结果。
- `scope`: `process/process.go`
- `recommendation`: 双人局应显式按阶段区分起手位，保证翻牌后由大盲先行动。
- `status`: confirmed

### [high] `LastRaiseAmount` 被写成整次投入而不是真实 raise delta

- `title`: `LastRaiseAmount` 被写成整次投入而不是真实 raise delta
- `angle`: correctness
- `severity`: high
- `confidence`: high
- `evidence`: `process/process.go:502-503` 先更新 `game.CurrentAmount` 再计算 `game.LastRaiseAmount`；`process/process.go:519-524` 的 `ALL_IN` 分支重复了同样的顺序错误。
- `impact`: 后续最小加注门槛会被错误抬高，合法的 raise / reopen 可能被拒绝。
- `scope`: `process/process.go`
- `recommendation`: 先保存旧的 `CurrentAmount`，再按 `newCurrentAmount - oldCurrentAmount` 计算 raise delta，并补直接回归测试。
- `status`: confirmed

### [high] 训练统计存在并发 map 写入风险

- `title`: 训练统计存在并发 map 写入风险
- `angle`: correctness
- `severity`: high
- `confidence`: high
- `evidence`: `game/unlimited/training.go:27` 建立共享 map，`game/unlimited/training.go:30-32` 启动 10 个 worker，`game/unlimited/training.go:78` 在无锁情况下执行 `(*memory)[finalWinnerIndex]++`。
- `impact`: 训练过程可能随机 panic，或者产生被竞争写坏的统计结果。
- `scope`: `game/unlimited/training.go`
- `recommendation`: 改为 worker 本地累计后统一合并，或使用互斥锁保护共享状态。
- `status`: confirmed

### [medium] 连续三次非法动作后会直接 panic

- `title`: 连续三次非法动作后会直接 panic
- `angle`: correctness
- `severity`: medium
- `confidence`: high
- `evidence`: `process/process.go:417-440` 在三次非法动作后仍无条件调用 `performAction`；`process/process.go:528-529` 会对零值 `ActionType` 触发 panic。
- `impact`: 错误 AI、边界输入或未来规则扩展中的异常动作会直接终止整局，而不是安全降级。
- `scope`: `process/process.go`
- `recommendation`: 达到重试上限后返回安全动作，例如 `FOLD`，或显式向上返回错误。
- `status`: confirmed

### [medium] 绿色测试没有覆盖真正高风险的流程层

- `title`: 绿色测试没有覆盖真正高风险的流程层
- `angle`: verification
- `severity`: medium
- `confidence`: high
- `evidence`: `go test ./... -cover` 只有 `poker/process` 有覆盖率，且仅 `35.5%`；`go tool cover -func=.audit-work/coverage2.out` 显示 `process/process.go` 的 `PlayGame`、`interactWithPlayers`、`checkAction`、`performAction`、`Start` 全部是 `0.0%`。
- `impact`: 行动顺序、加注规则、训练模式、入口启动等回归都可以在“测试全绿”的情况下漏出。
- `scope`: 全仓库验证策略
- `recommendation`: 至少补 `PlayGame`、`checkAction`、`performAction` 和一个完整 hand flow 的回归测试。
- `status`: confirmed

### [medium] 入口程序依赖开发机路径且无法运行时配置

- `title`: 入口程序依赖开发机路径且无法运行时配置
- `angle`: ops-config
- `severity`: medium
- `confidence`: high
- `evidence`: `main.go:17-28` 写死 `switch 1`；`main.go:51`, `main.go:62`, `main.go:72` 写死 `D:/Git/go/src/poker/...` 输出路径。
- `impact`: 二进制默认只能按作者当前机器布局工作，切换模式或路径都要改源码。
- `scope`: `main.go`
- `recommendation`: 改成 flag / config 驱动，默认输出到相对路径或临时目录。
- `status`: confirmed

## Angle-Based Assessment

### Correctness

- 记分与结算辅助逻辑有一定测试，但高风险问题集中在流程编排而不是纯函数。
- 当前最严重的问题分别落在训练模式信息泄露、双人局起手顺序、加注状态维护和训练统计并发安全。

### Verification

- 测试结果是绿的，但覆盖非常窄。
- 现有用例几乎只证明 `score.go` 和部分 `settle.go` 没明显问题，无法证明对局流程正确。

### Architecture

- 包边界对小项目来说还算清楚。
- 主要架构弱点是 `config/` 中的全局可变状态把渲染、训练、并发池和流程控制绑得很紧。

### Maintainability

- 代码规模不大，但运行约束缺少文档。
- `README.md` 几乎没有信息；仓库中混有历史日志、profile 和二进制，增加了判断噪音。

### Ops and Config

- 通过标准 Go 命令可以构建。
- 但入口默认假设固定 Windows 路径，且宿主机策略问题需要手工绕开，说明运行方式缺乏可移植性。

## Coverage and Unknowns

- Reviewed:
  - 所有 Go 源文件
  - 现有测试与覆盖率
  - 入口、训练、交互、结算、工具代码
- Not reviewed:
  - 真实人机交互牌局
  - Fyne GUI 实际运行
  - 仓库外脚本或 CI
- Why not reviewed:
  - 仓库内没有交互/GUI 自动化 harness
  - 本次为 guided audit，不额外搭建新执行脚本
- Limits on confidence:
  - 源码证据支持的发现置信度高
  - 真实交互体验与 GUI 结论置信度较低

## Recommended Next Steps

### Immediate

- 修复训练模式信息泄露。
- 修复双人局 post-flop 起手位。
- 修复 `LastRaiseAmount` 的更新顺序。
- 去掉训练统计的共享 map 并发写。

### This Week

- 为 `checkAction` / `performAction` / `PlayGame` 增加回归测试。
- 加一组 heads-up 行动顺序测试。
- 把入口模式和输出路径迁移到 flag 或配置。

### Later Governance

- 建立最小 CI：`go test ./...` + 覆盖率检查。
- 补 README，明确运行模式、训练模式与生成产物策略。
