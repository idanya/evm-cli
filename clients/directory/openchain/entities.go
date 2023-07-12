package openchain

type EntityType string

const (
	EntityTypeFunction EntityType = "function"
	EntityTypeEvent    EntityType = "event"
)

type LookupEntity struct {
	Name     string `json:"name"`
	Filtered bool   `json:"filtered"`
}

type LookupResponseResult struct {
	Function map[string][]*LookupEntity `json:"function"`
	Event    map[string][]*LookupEntity `json:"event"`
}

type LookupResponse struct {
	Ok     bool                 `json:"ok"`
	Result LookupResponseResult `json:"result"`
}
