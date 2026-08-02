package domain

// Repository is a placeholder extension point for the persistence layer.
// Part 4 will implement concrete adapters and SQL-backed repositories.
type Repository interface {
	SaveUser(User) error
	GetUserByID(string) (User, error)
}
