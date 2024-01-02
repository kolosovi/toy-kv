package walmanager

type reader struct {
}

func (r *reader) HasNext() bool {
	return false
}

func (r *reader) Scan(log *Log) error {
	_ = log
	return nil
}

func (r *reader) Err() error {
	return nil
}
