package publisher

type Publisher interface {
	Publish(event string, payload map[string]interface{}) error
}
