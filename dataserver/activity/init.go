package activity

// init 初始化
func init() {
	RegisterFlowCreator("to", NewFlowTo)
	RegisterFlowCreator("ifto", NewFlowIfTo)
	RegisterFlowCreator("loop", NewFlowLoop)

	RegisterAcitvityCreator("stdout", NewStdOutActivity)
}
