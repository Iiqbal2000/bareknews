package bareknews

type Status string

const (
	Publish Status = "publish"
	Draft   Status = "draft"
	Deleted Status = "deleted"
)
