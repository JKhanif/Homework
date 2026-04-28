package api_model

import "perfume-bot/model/db_model"

// type ProductResponse struct {
// 	ID          int                `json:"id"`
// 	Title       string             `json:"title"`
// 	Description string             `json:"description"`
// 	Price       int                `json:"price"`
// 	Brand       *BrandResponse     `json:"brand,omitempty"`
// 	Categories  []CategoryResponse `json:"categories"`
// 	Photos      []PhotoResponse    `json:"photos"`
// }

// type BrandResponse struct {
// 	ID          int    `json:"id"`
// 	Title       string `json:"title"`
// 	Description string `json:"description,omitempty"`
// }

// type CategoryResponse struct {
// 	ID    int    `json:"id"`
// 	Title string `json:"title"`
// }

type CreateProductRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	BrandID     *int   `json:"brand_id"`
	CategoryIDs []int  `json:"category_ids"`
}

type UpdateProductRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Price       *int    `json:"price"`
	BrandID     *int64  `json:"brand_id"`
}

type PhotoResponse struct {
	ID       int64  `json:"id"`
	URL      string `json:"url"`
	IsMain   bool   `json:"is_main"`
	TgFileID string `json:"tg_file_id"`
}

type ProductFull struct {
	Product    db_model.Product
	Brand      *db_model.Brand
	Categories []db_model.Category
	Photos     []db_model.ProductPhoto
}
