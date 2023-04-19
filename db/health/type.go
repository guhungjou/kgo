package health

import (
	"time"

	kindergartendb "gitlab.com/ykgk/kgo/db/kindergarten"
)

type KindergartenStudentMorningCheck struct {
	tableName         struct{}  `pg:"kindergarten_student_morning_check" json:"-"`
	ID                int64     `pg:"id,notnull,pk" json:"id"`
	KindergartenID    int64     `pg:"kindergarten_id,notnull" json:"kindergarten_id"`
	StudentID         int64     `pg:"student_id,notnull" json:"student_id"`
	Date              time.Time `pg:"date" json:"date"`
	Temperature       float64   `pg:"temperature,notnull,use_zero" json:"temperature"`
	TemperatureStatus string    `pg:"temperature_status,notnull,use_zero" json:"temperature_status"`
	Hand              string    `pg:"hand,notnull,use_zero" json:"hand"`
	HandStatus        string    `pg:"hand_status,notnull,use_zero" json:"hand_status"`
	Mouth             string    `pg:"mouth,notnull,use_zero" json:"mouth"`
	MouthStatus       string    `pg:"mouth_status,notnull,use_zero" json:"mouth_status"`
	Eye               string    `pg:"eye,notnull,use_zero" json:"eye"`
	EyeStatus         string    `pg:"eye_status,notnull,use_zero" json:"eye_status"`

	CreatedAt time.Time `pg:"created_at" json:"created_at"`
	UpdatedAt time.Time `pg:"updated_at" json:"updated_at"`

	Deleted bool `pg:"deleted,notnull,use_zero" json:"-"`

	Student      *kindergartendb.KindergartenStudent `pg:"rel:has-one" json:"student,omitempty"`
	Kindergarten *kindergartendb.Kindergarten        `pg:"rel:has-one" json:"kindergarten,omitempty"`
}

type KindergartenStudentMedicalExamination struct {
	tableName      struct{}  `pg:"kindergarten_student_medical_examination" json:"-"`
	ID             int64     `pg:"id,notnull,pk" json:"id"`
	KindergartenID int64     `pg:"kindergarten_id,notnull" json:"kindergarten_id"`
	StudentID      int64     `pg:"student_id,notnull" json:"student_id"`
	Date           time.Time `pg:"date" json:"date"`

	Height          float64   `pg:"height,notnull,use_zero" json:"height"`
	HeightStatus    string    `pg:"height_status,notnull,use_zero" json:"height_status"`
	HeightUpdatedAt time.Time `pg:"height_updated_at" json:"height_updated_at"`

	Weight          float64   `pg:"weight,notnull,use_zero" json:"weight"`
	WeightStatus    string    `pg:"weight_status,notnull,use_zero" json:"weight_status"`
	WeightUpdatedAt time.Time `pg:"weight_updated_at" json:"weight_updated_at"`

	BMI          float64   `pg:"bmi,notnull,use_zero" json:"bmi"`
	BMIStatus    string    `pg:"bmi_status,notnull,use_zero" json:"bmi_status"`
	BMIUpdatedAt time.Time `pg:"bmi_updated_at" json:"bmi_updated_at"`

	Hemoglobin          float64   `pg:"hemoglobin,notnull,use_zero" json:"hemoglobin"`
	HemoglobinRemark    string    `pg:"hemoglobin_remark,notnull,use_zero" json:"hemoglobin_remark"`
	HemoglobinStatus    string    `pg:"hemoglobin_status,notnull,use_zero" json:"hemoglobin_status"`
	HemoglobinUpdatedAt time.Time `pg:"hemoglobin_updated_at" json:"hemoglobin_updated_at"`

	// SightL         float64   `pg:"sight_l,notnull,use_zero" json:"sight_l"`
	SightLS       string `pg:"sight_l_s,notnull,use_zero" json:"sight_l_s"`
	SightLC       string `pg:"sight_l_c,notnull,use_zero" json:"sight_l_c"`
	SightLRemark  string `pg:"sight_l_remark,notnull,use_zero" json:"sight_l_remark"`
	SightLSStatus string `pg:"sight_l_s_status,notnull,use_zero" json:"sight_l_s_status"`
	SightLCStatus string `pg:"sight_l_c_status,notnull,use_zero" json:"sight_l_c_status"`
	// SightLStatus string `pg:"sight_l_status,notnull,use_zero" json:"sight_l_status"`
	// SightR         float64   `pg:"sight_r,notnull,use_zero" json:"sight_r"`
	SightRS string `pg:"sight_r_s,notnull,use_zero" json:"sight_r_s"`
	SightRC string `pg:"sight_r_c,notnull,use_zero" json:"sight_r_c"`

	SightRRemark  string `pg:"sight_r_remark,notnull,use_zero" json:"sight_r_remark"`
	SightRSStatus string `pg:"sight_r_s_status,notnull,use_zero" json:"sight_r_s_status"`
	SightRCStatus string `pg:"sight_r_c_status,notnull,use_zero" json:"sight_r_c_status"`
	// SightRStatus   string    `pg:"sight_r_status,notnull,use_zero" json:"sight_r_status"`
	SightUpdatedAt time.Time `pg:"sight_updated_at" json:"sight_updated_at"`

	ToothCount       int       `pg:"tooth_count,notnull,use_zero" json:"tooth_count"`
	ToothCariesCount int       `pg:"tooth_caries_count,notnull,use_zero" json:"tooth_caries_count"`
	ToothRemark      string    `pg:"tooth_remark,notnull,use_zero" json:"tooth_remark"`
	ToothUpdatedAt   time.Time `pg:"tooth_updated_at" json:"tooth_updated_at"`

	ALT          float64   `pg:"alt,notnull,use_zero" json:"alt"`
	ALTRemark    string    `pg:"alt_remark,notnull,use_zero" json:"alt_remark"`
	ALTStatus    string    `pg:"alt_status,notnull,use_zero" json:"alt_status"`
	ALTUpdatedAt time.Time `pg:"alt_updated_at" json:"alt_updated_at"`

	CreatedAt time.Time `pg:"created_at" json:"created_at"`
	UpdatedAt time.Time `pg:"updated_at" json:"updated_at"`

	Deleted bool `pg:"deleted,notnull,use_zero" json:"-"`

	Advise string `pg:"-" json:"advise"`

	EyeLeft        float32                             `pg:"eye_left,notnull,use_zero" json:"eye_left"`
	EyeRight       float32                             `pg:"eye_right,notnull,use_zero" json:"eye_right"`
	EyeLeftStatus  string                              `pg:"eye_left_status,notnull,use_zero" json:"eye_left_status"`
	EyeRightStatus string                              `pg:"eye_right_status,notnull,use_zero" json:"eye_right_status"`
	Student        *kindergartendb.KindergartenStudent `pg:"rel:has-one" json:"student,omitempty"`
	Kindergarten   *kindergartendb.Kindergarten        `pg:"rel:has-one" json:"kindergarten,omitempty"`
}

type StandardScaleHL struct {
	tableName struct{} `pg:"standard_scale_hl" json:"-"`
	ID        int64    `pg:"id,notnull,pk" json:"id"`

	Name   string  `pg:"name,notnull" json:"name"`
	Gender string  `pg:"gender,notnull,use_zero" json:"gender"`
	Age    float64 `pg:"age,notnull,use_zero" json:"age"`

	Min float64 `pg:"min,notnull,use_zero" json:"min"`
	Max float64 `pg:"max,notnull,use_zero" json:"max"`

	CreatedAt time.Time `pg:"created_at" json:"created_at"`
	UpdatedAt time.Time `pg:"updated_at" json:"updated_at"`
}

type KindergartenStudentFitnessTest struct {
	tableName struct{} `pg:"kindergarten_student_fitness_test,discard_unknown_columns" json:"-"`

	ID             int64     `pg:"id,pk" json:"id"`
	KindergartenID int64     `pg:"kindergarten_id,notnull" json:"kindergarten_id"`
	StudentID      int64     `pg:"student_id,notnull" json:"student_id"`
	Date           time.Time `pg:"date" json:"date"`

	Height          float64   `pg:"height,notnull,use_zero" json:"height"`
	HeightUpdatedAt time.Time `pg:"height_updated_at" json:"height_updated_at"`
	Weight          float64   `pg:"weight,notnull,use_zero" json:"weight"`
	WeightUpdatedAt time.Time `pg:"weight_updated_at" json:"weight_updated_at"`

	HeightAndWeightScore float64 `pg:"height_and_weight_score,notnull,use_zero" json:"height_and_weight_score"`

	ShuttleRun10          float64   `pg:"shuttle_run_10,notnull,use_zero" json:"shuttle_run_10"`
	ShuttleRun10Score     float64   `pg:"shuttle_run_10_score,notnull,use_zero" json:"shuttle_run_10_score"`
	ShuttleRun10UpdatedAt time.Time `pg:"shuttle_run_10_updated_at" json:"shuttle_run_10_updated_at"`

	StandingLongJump          float64   `pg:"standing_long_jump,notnull,use_zero" json:"standing_long_jump"`
	StandingLongJumpScore     float64   `pg:"standing_long_jump_score,notnull,use_zero" json:"standing_long_jump_score"`
	StandingLongJumpUpdatedAt time.Time `pg:"standing_long_jump_updated_at" json:"standing_long_jump_updated_at"`

	BaseballThrow          float64   `pg:"baseball_throw,notnull,use_zero" json:"baseball_throw"`
	BaseballThrowScore     float64   `pg:"baseball_throw_score,notnull,use_zero" json:"baseball_throw_score"`
	BaseballThrowUpdatedAt time.Time `pg:"baseball_throw_updated_at" json:"baseball_throw_updated_at"`

	BunnyHopping          float64   `pg:"bunny_hopping,notnull,use_zero" json:"bunny_hopping"`
	BunnyHoppingScore     float64   `pg:"bunny_hopping_score,notnull,use_zero" json:"bunny_hopping_score"`
	BunnyHoppingUpdatedAt time.Time `pg:"bunny_hopping_updated_at" json:"bunny_hopping_updated_at"`

	SitAndReach          float64   `pg:"sit_and_reach,notnull,use_zero" json:"sit_and_reach"`
	SitAndReachScore     float64   `pg:"sit_and_reach_score,notnull,use_zero" json:"sit_and_reach_score"`
	SitAndReachUpdatedAt time.Time `pg:"sit_and_reach_updated_at" json:"sit_and_reach_updated_at"`

	BalanceBeam          float64   `pg:"balance_beam,notnull,use_zero" json:"balance_beam"`
	BalanceBeamScore     float64   `pg:"balance_beam_score,notnull,use_zero" json:"balance_beam_score"`
	BalanceBeamUpdatedAt time.Time `pg:"balance_beam_updated_at" json:"balance_beam_updated_at"`

	TotalScore  float64 `pg:"total_score,notnull,use_zero" json:"total_score"`
	TotalStatus string  `pg:"total_status,notnull,use_zero" json:"total_status"`

	Deleted   bool      `pg:"deleted,notnull,use_zero" json:"deleted"`
	CreatedAt time.Time `pg:"created_at" json:"created_at"`
	UpdatedAt time.Time `pg:"updated_at" json:"updated_at"`

	Student      *kindergartendb.KindergartenStudent `pg:"rel:has-one" json:"student,omitempty"`
	Kindergarten *kindergartendb.Kindergarten        `pg:"rel:has-one" json:"kindergarten,omitempty"`
}
