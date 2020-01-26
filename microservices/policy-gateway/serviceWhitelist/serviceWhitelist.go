package serviceWhitelist

type ServiceWhitelistEntry struct {
	Method string
	URL    string
	Allow  []string
}

type ServiceWhitelist struct {
	Entries []ServiceWhitelistEntry
}
