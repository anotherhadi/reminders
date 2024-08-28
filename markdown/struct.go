package markdown

type TodoItem struct {
	Label string
	State string

	Link string

	Importance string // Stored in the link markdown's metadata
	Content    string // The content of the link markdown's body
}

type List struct {
	Title     string
	TodoItems []TodoItem
}

type Md struct {
	DirPath string
	File    string
	Title   string
	Lists   []List
}
