package langwatch

const (
	AttributeNameUserID     = "user_id"
	AttributeNameThreadID   = "thread_id"
	AttributeNameCustomerID = "customer_id"
)

type attribute[T any] struct {
	Name  string
	Value T
}

func Attribute[T any](name string, value T) attribute[any] {
	return attribute[any]{
		Name:  name,
		Value: value,
	}
}
