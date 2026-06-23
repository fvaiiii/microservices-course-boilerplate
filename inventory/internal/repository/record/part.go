package record

import "time"

type Part struct {
	UUID          string     `db:"uuid"`
	Name          string     `db:"name"`
	Description   string     `db:"description"`
	PartType      string     `db:"part_type"`
	Price         int64      `db:"price"`
	StockQuantity int64      `db:"stock_quantity"`
	CreatedAt     time.Time  `db:"created_at"`
	UpdatedAt     *time.Time `db:"updated_at"`
}
