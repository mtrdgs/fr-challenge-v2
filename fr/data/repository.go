package data

// RepositoryPattern - interface to be used as repository
type RepositoryPattern interface {
	Insert(entry QuoteEntry) error
	FindSpecific(amount int64) (quotes []QuoteEntry, err error)
}
