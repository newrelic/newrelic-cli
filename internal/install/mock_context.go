package install

import "time"

type MockContext struct {
}

func (m MockContext) Deadline() (deadline time.Time, ok bool) {
	//TODO implement me
	//panic("implement me")
	return
}

func (m MockContext) Done() <-chan struct{} {
	//TODO implement me
	//panic("implement me")
	return nil
}

func (m MockContext) Err() error {
	//TODO implement me
	//panic("implement me")
	return nil
}

func (m MockContext) Value(key interface{}) interface{} {
	//TODO implement me
	//panic("implement me")
	return nil
}
