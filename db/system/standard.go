package system

import (
	"github.com/go-pg/pg/v10"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ykgk/kgo/db"
)

func FindStandardScaleScoresByName(name string) ([]*StandardScaleScore, error) {
	scores := make([]*StandardScaleScore, 0)

	if err := db.PG().Model(&scores).Where(`"name"=?`, name).Select(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}
	return scores, nil
}

func GetStandardScaleScoreTx(tx *pg.Tx, name string, gender string, age float64, v float64) (float64, error) {
	score := StandardScaleScore{}

	if err := tx.Model(&score).Where(`"name"=?`, name).Where(`"gender"=?`, gender).
		Where(`"age"=?`, age).Where(`"min"<=?`, v).
		Where(`"max">=?`, v).Limit(1).Select(); err != nil {
		if err != db.ErrNoRows {
			log.Errorf("DB Error: %v", err)
			return 0, err
		}
		return 0, nil
	}
	return score.Score, nil
}

func GetStandardSclaeHWScoreTx(tx *pg.Tx, gender string, height, weight float64) (float64, error) {
	score := StandardScaleHWScore{}

	if err := tx.Model(&score).Where(`"gender"=?`, gender).
		Where(`"height_min"<=?`, height).Where(`"height_max">=?`, height).
		Where(`"weight_min"<=?`, weight).Where(`"weight_max">=?`, weight).Limit(1).Select(); err != nil {
		if err != db.ErrNoRows {
			log.Errorf("DB Error: %v", err)
			return 0, err
		}
		return 0, nil
	}
	return score.Score, nil
}
