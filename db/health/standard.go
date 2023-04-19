package health

import (
	log "github.com/sirupsen/logrus"
	"gitlab.com/ykgk/kgo/db"
)

// type Range struct {
// 	Min float64
// 	Max float64
// }

// func (r *Range) GetStatus(v float64) string {
// 	if r == nil {
// 		return ""
// 	}
// 	if r.Min > v {
// 		return StandardStatusLow
// 	} else if r.Max < v {
// 		return StandardStatusHigh
// 	}
// 	return StandardStatusNormal
// }

func (r *StandardScaleHL) GetStatus(v float64) string {
	if r == nil {
		return ""
	}
	if r.Min > v {
		return StandardStatusLow
	} else if r.Max != -1 && r.Max < v {
		return StandardStatusHigh
	}
	return StandardStatusNormal
}

/* https://www.youlai.cn/dise/imagedetail/157_32867.html */
// var MaleHeightStandard = map[float64]*Range{
// 	2: {
// 		Min: 84.3,
// 		Max: 91,
// 	},
// 	2.5: {
// 		Min: 88.9,
// 		Max: 98.5,
// 	},
// 	3: {
// 		Min: 91.1,
// 		Max: 98.7,
// 	},
// 	3.5: {
// 		Min: 95,
// 		Max: 103.1,
// 	},
// 	4: {
// 		Min: 98.7,
// 		Max: 107.2,
// 	},
// 	4.5: {
// 		Min: 102.1,
// 		Max: 111,
// 	},
// 	5: {
// 		Min: 105.3,
// 		Max: 114.5,
// 	},
// 	5.5: {
// 		Min: 108.4,
// 		Max: 117.8,
// 	},
// 	6: {
// 		Min: 111.2,
// 		Max: 121.0,
// 	},
// 	7: {
// 		Min: 116.6,
// 		Max: 126.8,
// 	},
// 	8: {
// 		Min: 121.6,
// 		Max: 132.2,
// 	},
// 	9: {
// 		Min: 126.5,
// 		Max: 137.8,
// 	},
// 	10: {
// 		Min: 131.4,
// 		Max: 143.6,
// 	},
// }

// var MaleWeightStandard = map[float64]*Range{
// 	2: {
// 		Min: 11.2,
// 		Max: 14.0,
// 	},
// 	2.5: {
// 		Min: 12.1,
// 		Max: 15.3,
// 	},
// 	3: {
// 		Min: 13.0,
// 		Max: 16.4,
// 	},
// 	3.5: {
// 		Min: 13.9,
// 		Max: 17.6,
// 	},
// 	4: {
// 		Min: 14.8,
// 		Max: 18.7,
// 	},
// 	4.5: {
// 		Min: 15.7,
// 		Max: 19.9,
// 	},
// 	5: {
// 		Min: 16.6,
// 		Max: 21.1,
// 	},
// 	5.5: {
// 		Min: 17.4,
// 		Max: 22.3,
// 	},
// 	6: {
// 		Min: 18.4,
// 		Max: 23.6,
// 	},
// 	7: {
// 		Min: 20.2,
// 		Max: 26.5,
// 	},
// 	8: {
// 		Min: 22.2,
// 		Max: 30.0,
// 	},
// 	9: {
// 		Min: 24.3,
// 		Max: 34.0,
// 	},
// 	10: {
// 		Min: 26.8,
// 		Max: 38.7,
// 	},
// }

// var FemaleHeightStandard = map[float64]*Range{
// 	2: {
// 		Min: 83.8,
// 		Max: 89.8,
// 	},
// 	2.5: {
// 		Min: 87.9,
// 		Max: 94.7,
// 	},
// 	3: {
// 		Min: 90.2,
// 		Max: 98.1,
// 	},
// 	3.5: {
// 		Min: 94.0,
// 		Max: 101.8,
// 	},
// 	4: {
// 		Min: 97.6,
// 		Max: 105.7,
// 	},
// 	4.5: {
// 		Min: 100.9,
// 		Max: 109.3,
// 	},
// 	5: {
// 		Min: 104.0,
// 		Max: 112.8,
// 	},
// 	5.5: {
// 		Min: 106.9,
// 		Max: 116.2,
// 	},
// 	6: {
// 		Min: 109.7,
// 		Max: 119.6,
// 	},
// 	7: {
// 		Min: 115.1,
// 		Max: 126.2,
// 	},
// 	8: {
// 		Min: 120.4,
// 		Max: 132.4,
// 	},
// 	9: {
// 		Min: 125.7,
// 		Max: 138.7,
// 	},
// 	10: {
// 		Min: 131.5,
// 		Max: 145.1,
// 	},
// }

// var FemaleWeightStandard = map[float64]*Range{
// 	2: {
// 		Min: 10.6,
// 		Max: 13.2,
// 	},
// 	2.5: {
// 		Min: 11.7,
// 		Max: 14.7,
// 	},
// 	3: {
// 		Min: 12.6,
// 		Max: 16.1,
// 	},
// 	3.5: {
// 		Min: 13.5,
// 		Max: 17.2,
// 	},
// 	4: {
// 		Min: 14.3,
// 		Max: 18.3,
// 	},
// 	4.5: {
// 		Min: 15,
// 		Max: 19.4,
// 	},
// 	5: {
// 		Min: 15.7,
// 		Max: 20.4,
// 	},
// 	5.5: {
// 		Min: 16.5,
// 		Max: 21.6,
// 	},
// 	6: {
// 		Min: 17.3,
// 		Max: 22.9,
// 	},
// 	7: {
// 		Min: 19.1,
// 		Max: 26,
// 	},
// 	8: {
// 		Min: 21.4,
// 		Max: 30.2,
// 	},
// 	9: {
// 		Min: 24.1,
// 		Max: 35.3,
// 	},
// 	10: {
// 		Min: 27.2,
// 		Max: 40.9,
// 	},
// }

const (
	StandardStatusNormal = "normal" /* 正常 */
	StandardStatusLow    = "low"    /* 偏低 */
	StandardStatusHigh   = "high"   /* 偏高 */
)

func GetStandardScaleHL(name, gender string, age float64) (*StandardScaleHL, error) {
	r := StandardScaleHL{}

	if err := db.PG().Model(&r).Where(`"name"=?`, name).Where(`"gender"=?`, gender).Where(`"age"=?`, age).Limit(1).Select(); err != nil {
		if err != db.ErrNoRows {
			log.Errorf("DB Error: %v", err)
			return nil, err
		}
		return nil, nil
	}
	return &r, nil
}

func GetHeightStatus(age float64, gender string, v float64) string {
	r, _ := GetStandardScaleHL("身高", gender, age)
	if r == nil {
		r, _ = GetStandardScaleHL("身高", gender, float64(int(age)))
	}
	return r.GetStatus(v)
}

func GetWeightStatus(age float64, gender string, v float64) string {
	r, _ := GetStandardScaleHL("体重", gender, age)
	if r == nil {
		r, _ = GetStandardScaleHL("体重", gender, float64(int(age)))
	}
	return r.GetStatus(v)
}

func GetBMIStatus(age float64, gender string, v float64) string {
	r, _ := GetStandardScaleHL("BMI", gender, age)
	if r == nil {
		r, _ = GetStandardScaleHL("BMI", gender, float64(int(age)))
	}
	return r.GetStatus(v)
}

// var HemoglobinStandard = Range{
// 	Min: 110,
// 	Max: 160,
// }

func GetHemoglobinStatus(v float64) string {
	r, _ := GetStandardScaleHL("血红蛋白", "", 0)
	return r.GetStatus(v)
}

// var ALTStandard = Range{
// 	Min: 0,
// 	Max: 40,
// }

func GetALTStatus(v float64) string {
	r, _ := GetStandardScaleHL("谷丙转氨酶", "", 0)
	return r.GetStatus(v)
}

// var SightStandard = Range{
// 	Min: 4.8,
// 	Max: 10,
// }

// func GetSightStatus(v float64) string {
// 	r, _ := GetStandardScaleHL("视力", "", 0)
// 	return r.GetStatus(v)
// }

func GetSightSStatus(v float64) string {
	r, _ := GetStandardScaleHL("球镜度", "", 0)
	return r.GetStatus(v)
}

func GetSightCStatus(v float64) string {
	r, _ := GetStandardScaleHL("柱镜度", "", 0)
	return r.GetStatus(v)
}

// var TemperatureStandard = Range{
// 	Min: 36,
// 	Max: 37.3,
// }

func GetTemperatureStatus(v float64) string {
	r, _ := GetStandardScaleHL("体温", "", 0)
	return r.GetStatus(v)
}

// /*
//  * https://www.cdc.gov/growthcharts/data/set1clinical/cj41l023.pdf
//  */
// var MaleBMIStandard = map[float64]*Range{
// 	2: {
// 		Min: 14.8,
// 		Max: 18.2,
// 	},
// 	2.5: {
// 		Min: 14.6,
// 		Max: 17.8,
// 	},
// 	3: {
// 		Min: 14.4,
// 		Max: 17.4,
// 	},
// 	3.5: {
// 		Min: 14.2,
// 		Max: 17.1,
// 	},
// 	4: {
// 		Min: 14,
// 		Max: 17,
// 	},
// 	4.5: {
// 		Min: 14,
// 		Max: 16.8,
// 	},
// 	5: {
// 		Min: 13.8,
// 		Max: 16.8,
// 	},
// 	5.5: {
// 		Min: 13.8,
// 		Max: 16.9,
// 	},
// 	6: {
// 		Min: 13.8,
// 		Max: 17,
// 	},
// 	6.5: {
// 		Min: 13.7,
// 		Max: 17.2,
// 	},
// 	7: {
// 		Min: 13.8,
// 		Max: 17.4,
// 	},
// 	7.5: {
// 		Min: 13.8,
// 		Max: 17.6,
// 	},
// 	8: {
// 		Min: 13.8,
// 		Max: 17.9,
// 	},
// 	8.5: {
// 		Min: 13.8,
// 		Max: 18.3,
// 	},
// 	9: {
// 		Min: 14,
// 		Max: 18.6,
// 	},
// 	9.5: {
// 		Min: 14.1,
// 		Max: 18.9,
// 	},
// 	10: {
// 		Min: 14.2,
// 		Max: 19.4,
// 	},
// }

// func GetMaleBMIStatus(birthday time.Time, v float64) string {
// 	age := getAge(birthday)
// 	r := MaleBMIStandard[age]
// 	return r.GetStatus(v)
// }

// /*
//  * https://www.cdc.gov/growthcharts/data/set1clinical/cj41l024.pdf
//  */
// var FemaleBMIStandard = map[float64]*Range{
// 	2: {
// 		Min: 14.4,
// 		Max: 18,
// 	},
// 	2.5: {
// 		Min: 14.2,
// 		Max: 17.6,
// 	},
// 	3: {
// 		Min: 14,
// 		Max: 17.2,
// 	},
// 	3.5: {
// 		Min: 13.8,
// 		Max: 17,
// 	},
// 	4: {
// 		Min: 13.8,
// 		Max: 16.8,
// 	},
// 	4.5: {
// 		Min: 13.6,
// 		Max: 16.8,
// 	},
// 	5: {
// 		Min: 13.5,
// 		Max: 16.8,
// 	},
// 	5.5: {
// 		Min: 13.5,
// 		Max: 16.9,
// 	},
// 	6: {
// 		Min: 13.4,
// 		Max: 17.1,
// 	},
// 	6.5: {
// 		Min: 13.4,
// 		Max: 17.3,
// 	},
// 	7: {
// 		Min: 13.4,
// 		Max: 17.6,
// 	},
// 	7.5: {
// 		Min: 13.5,
// 		Max: 17.9,
// 	},
// 	8: {
// 		Min: 13.6,
// 		Max: 18.3,
// 	},
// 	8.5: {
// 		Min: 13.6,
// 		Max: 18.5,
// 	},
// 	9: {
// 		Min: 13.8,
// 		Max: 19.1,
// 	},
// 	9.5: {
// 		Min: 13.9,
// 		Max: 19.3,
// 	},
// 	10: {
// 		Min: 14,
// 		Max: 19.7,
// 	},
// }

// func GetFemaleBMIStatus(birthday time.Time, v float64) string {
// 	age := getAge(birthday)
// 	r := FemaleBMIStandard[age]
// 	return r.GetStatus(v)
// }
