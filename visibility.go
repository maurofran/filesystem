package filesystem

// Visibility enumeration.
type Visibility int

// Visibility values (public and private).
const (
	VisibilityPublic Visibility = iota + 1
	VisibilityPrivate
)

var visibilities = [...]string{"Public", "Private"}

func (v Visibility) String() string {
	return visibilities[v-1]
}
