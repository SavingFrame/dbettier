package messages

type NvimFinishedMsg struct {
	Error    error
	FileName string
}
