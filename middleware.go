package rufus

// Middleware struct for rufus
type Middleware struct {
	RedirectToNonWWW    bool `json:"redirect_to_non_www,omitempty"`
	EnableResponseCache bool `json:"enable_response_cache,omitempty"`
	Cache               responseCache
}
