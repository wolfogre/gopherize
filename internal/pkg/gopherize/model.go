package gopherize

type Artwork struct {
	Categories        []ArtworkCategory `json:"categories"`
	TotalCombinations int64             `json:"total_combinations"`
}

type ArtworkCategory struct {
	Id     string                 `json:"id"`
	Name   string                 `json:"name"`
	Images []ArtworkCategoryImage `json:"images"`
}

type ArtworkCategoryImage struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	Href          string `json:"href"`
	ThumbnailHref string `json:"thumbnail_href"`
}
