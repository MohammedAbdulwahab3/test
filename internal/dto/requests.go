package dto

type ItemCreateRequest struct {
 	Name        string  `json:"name"`
 	Description string  `json:"description"`
 	Price       float64 `json:"price"`
 	Status      string  `json:"status"`
}

type ItemUpdateRequest struct {
 	Name        *string  `json:"name"`
 	Description *string  `json:"description"`
 	Price       *float64 `json:"price"`
 	Status      *string  `json:"status"`
 	Version     uint     `json:"version"`
}

type BulkDeleteRequest struct {
 	IDs  []uint `json:"ids"`
 	Hard bool   `json:"hard"`
}
