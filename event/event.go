package event

type Event interface {
	MatchString() string
	ShortString() string
	LongString() string
	Sent()
}
