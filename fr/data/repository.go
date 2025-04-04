package data

type RepositoryPattern interface {
	Insert(entry QuoteEntry) error
	FindSpecific(amount int64) (quotes []QuoteEntry, err error)
}
