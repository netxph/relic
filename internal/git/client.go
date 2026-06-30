package git

// RawCommit holds the raw data from a single git log entry.
type RawCommit struct {
	Hash    string
	Subject string
	Body    string
}

// GitClient abstracts git log querying.
type GitClient interface {
	Log(from, to string) ([]RawCommit, error)
}
