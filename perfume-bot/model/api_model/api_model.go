package api_model

type PhotoResponse struct {
	ID       int64  `json:"id"`
	URL      string `json:"url"`
	IsMain   bool   `json:"is_main"`
	TgFileID string `json:"tg_file_id"`
}

// BRAND BRAND BRAND BRAND BRAND BRAND BRAND BRAND BRAND BRAND BRAND BRAND BRAND BRAND BRAND BRAND BRAND BRAND

type BrandResponse struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
}

type CreateBrandRequest struct {
	Title       string  `json:"title" binding:"required"`
	Description *string `json:"description"`
}

type UpdateBrandRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description,omitempty"`
}

// CATEGORY CATEGORY CATEGORY CATEGORY CATEGORY CATEGORY CATEGORY CATEGORY CATEGORY CATEGORY CATEGORY CATEGORY

type CategoryResponse struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

type CreateCategoryRequest struct {
	Title string `json:"title" binding:"required"`
}

type UpdateCategoryRequest struct {
	Title *string `json:"title"`
}

// PRODUCT PRODUCT PRODUCT PRODUCT PRODUCT PRODUCT PRODUCT PRODUCT PRODUCT PRODUCT PRODUCT PRODUCT PRODUCT

type ProductResponse struct {
	ID          int                `json:"id"`
	Title       string             `json:"title"`
	Description *string            `json:"description"`
	Price       int                `json:"price"`
	Brand       *BrandResponse     `json:"brand,omitempty"`
	Categories  []CategoryResponse `json:"categories"`
	Photos      []PhotoResponse    `json:"photos"`
}

type CreateProductRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	BrandID     *int   `json:"brand_id"`
	CategoryIDs []int  `json:"category_ids"`
}

type UpdateProductRequest struct {
	Title       *string   `json:"title,omitempty"`
	BrandID     *int      `json:"brand_id,omitempty"`
	Price       *int      `json:"price,omitempty"`
	Description *string   `json:"description,omitempty"`
	Images      *[]string `json:"images,omitempty"`
}
