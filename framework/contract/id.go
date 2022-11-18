package contract

const IDKey = "fire:id"

type IDService interface {
	NewID() string
}
