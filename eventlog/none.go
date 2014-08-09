package eventlog

func NewMockEventLog() *mockEventLog {
	return &mockEventLog{}
}

type mockEventLog struct{}

func (log *mockEventLog) NewGame(gameId string)  {}
func (log *mockEventLog) Moved(gameId string)    {}
func (log *mockEventLog) Finished(gameId string) {}
