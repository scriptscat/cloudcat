package resource_entity

type Resource struct {
	URL        string `json:"url"`
	Content    string `json:"content"`
	Createtime int64  `json:"createtime"`
	Updatetime int64  `json:"updatetime"`
}
