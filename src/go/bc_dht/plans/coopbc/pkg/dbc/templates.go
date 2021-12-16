package dbc

type DocumentResponseStruct struct {
	Abstract string `json:"abstract"`
	ArxivID  string `json:"arxivId"`
	Authors  []struct {
		AuthorID string `json:"authorId"`
		Name     string `json:"name"`
		URL      string `json:"url"`
	} `json:"authors"`
	CitationVelocity int `json:"citationVelocity"`
	Citations        []struct {
		ArxivID interface{} `json:"arxivId"`
		Authors []struct {
			AuthorID string `json:"authorId"`
			Name     string `json:"name"`
		} `json:"authors"`
		Doi           interface{}   `json:"doi"`
		Intent        []interface{} `json:"intent"`
		IsInfluential bool          `json:"isInfluential"`
		PaperID       string        `json:"paperId"`
		Title         string        `json:"title"`
		URL           string        `json:"url"`
		Venue         string        `json:"venue"`
		Year          int           `json:"year"`
	} `json:"citations"`
	CorpusID                 int      `json:"corpusId"`
	Doi                      string   `json:"doi"`
	FieldsOfStudy            []string `json:"fieldsOfStudy"`
	InfluentialCitationCount int      `json:"influentialCitationCount"`
	IsOpenAccess             bool     `json:"isOpenAccess"`
	IsPublisherLicensed      bool     `json:"isPublisherLicensed"`
	NumCitedBy               int      `json:"numCitedBy"`
	NumCiting                int      `json:"numCiting"`
	PaperID                  string   `json:"paperId"`
	References               []struct {
		ArxivID interface{} `json:"arxivId"`
		Authors []struct {
			AuthorID string `json:"authorId"`
			Name     string `json:"name"`
		} `json:"authors"`
		Doi           string   `json:"doi"`
		Intent        []string `json:"intent"`
		IsInfluential bool     `json:"isInfluential"`
		PaperID       string   `json:"paperId"`
		Title         string   `json:"title"`
		URL           string   `json:"url"`
		Venue         string   `json:"venue"`
		Year          int      `json:"year"`
	} `json:"references"`
	Title  string `json:"title"`
	Topics []struct {
		Topic   string `json:"topic"`
		TopicID string `json:"topicId"`
		URL     string `json:"url"`
	} `json:"topics"`
	URL   string `json:"url"`
	Venue string `json:"venue"`
	Year  int    `json:"year"`
}
