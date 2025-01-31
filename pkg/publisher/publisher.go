package publisher

type EventType string

const (
	EventUserCreated EventType = "user.created"
	EventUserUpdated EventType = "user.updated"
	EventUserDeleted EventType = "user.deleted"
)

type Publisher interface {
	Publish(event EventType, payload map[string]interface{}) error
}
