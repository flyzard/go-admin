package valueobject

// SmartTableConfig defines configuration for a smart table
type SmartTableConfig struct {
	Columns      []SmartTableColumn
	DefaultSort  string
	DefaultOrder string
	PageSizes    []int
	Actions      []SmartTableAction
}

// SmartTableColumn defines a column in the smart table
type SmartTableColumn struct {
	Field      string
	Label      string
	Sortable   bool
	Filterable bool
	FilterType string // "text", "select", "date", etc.
	FilterOpts []FilterOption
	Formatter  string // Optional JS function name for client-side formatting
	Width      string // CSS width
	Visible    bool
	Template   string // Optional custom template for rendering
}

// FilterOption for select filters
type FilterOption struct {
	Value string
	Label string
}

// SmartTableAction defines an action that can be performed on rows
type SmartTableAction struct {
	Label    string
	Icon     string
	Action   string // URL or JS function
	Confirm  bool
	Message  string
	Class    string
	ShowWhen func(entity interface{}) bool
}
