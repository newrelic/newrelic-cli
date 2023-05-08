package recipes

const (
	embeddedEventsPath = "files/events.src"
)

type EmbeddedEventsFetcher struct{}

func NewEmbeddedEventsFetcher() *EmbeddedEventsFetcher {
	return &EmbeddedEventsFetcher{}
}

func (f *EmbeddedEventsFetcher) GetWriteKey() (string, error) {
	data, err := EmbeddedFS.ReadFile(embeddedEventsPath)
	if err != nil {
		return "", err
	}
	key := string(data)

	return key, nil
}
