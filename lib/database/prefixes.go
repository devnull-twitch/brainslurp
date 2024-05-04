package database

const (
	// Issues contain infromation as protobuf byte data
	IssuePrefix byte = 0x1
	// Views contain infromation as protobuf byte data
	ViewPrefix byte = 0x2
	// ViewIssue contains a list uint64 issue IDs for a given view
	ViewIssuesPrefix byte = 0x3
	// Project is a container for issues and views
	ProjectPrefix byte = 0x4
	// User contains login credentials and a name
	UserPrefix byte = 0x5
	// Username lookup maps usernames to user keys
	UsernameLookupPrefix byte = 0x6
	// Flows describe actions that change category or tags on an issue
	FlowPrefix byte = 0x7
	// Sequence are the prefix for 2 byte
	SequencePrefix byte = 0xFA
)
