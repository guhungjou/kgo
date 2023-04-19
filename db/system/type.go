package system

import "time"

type StandardScaleScore struct {
	tableName struct{} `pg:"standard_scale_score" json:"-"`

	ID     string  `pg:"id,pk" json:"id"`
	Name   string  `pg:"name,notnull" json:"name"`
	Gender string  `pg:"gender,notnull" json:"gender"`
	Age    float64 `pg:"age,notnull,use_zero" json:"age"`
	Min    float64 `pg:"min,notnull,use_zero" json:"min"`
	Max    float64 `pg:"max,notnull,use_zero" json:"max"`
	Score  float64 `pg:"score,notnull,use_zero" json:"score"`

	CreatedAt time.Time `pg:"created_at" json:"created_at"`
	UpdatedAt time.Time `pg:"updated_at" json:"updated_at"`
}

type StandardScaleHWScore struct {
	tableName struct{} `pg:"standard_scale_hw_score" json:"-"`

	ID        string  `pg:"id,pk" json:"id"`
	Gender    string  `pg:"gender,notnull" json:"gender"`
	HeightMin float64 `pg:"height_min,notnull" json:"height_min"`
	HeightMax float64 `pg:"height_max,notnull" json:"height_max"`
	WeightMin float64 `pg:"weight_min,notnull" json:"weight_min"`
	WeightMax float64 `pg:"weight_max,notnull" json:"weight_max"`

	Score float64 `pg:"score,notnull,use_zero" json:"score"`

	CreatedAt time.Time `pg:"created_at" json:"created_at"`
	UpdatedAt time.Time `pg:"updated_at" json:"updated_at"`
}

type District struct {
	tableName struct{} `pg:"district" json:"-"`
	ID        string   `pg:"id,pk" json:"id"`
	Name      string   `pg:"name,notnull,use_zero" json:"name"`
	Fullname  string   `pg:"fullname,notnull,use_zero" json:"fullname"`
	Pinyin    string   `pg:"pinyin,notnull,use_zero" json:"pinyin"`
	ParentID  string   `pg:"parent_id,notnull,use_zero" json:"parent_id"`

	CreatedAt time.Time `pg:"created_at" json:"created_at"`
	UpdatedAt time.Time `pg:"updated_at" json:"updated_at"`

	Parent *District `pg:"rel:has-one" json:"parent,omitempty"`
}
