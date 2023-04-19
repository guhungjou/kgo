package system

import (
	log "github.com/sirupsen/logrus"
	"gitlab.com/ykgk/kgo/db"
)

func FindDistricts(query string, parentID string) ([]*District, error) {
	districts := make([]*District, 0)
	q := db.PG().Model(&districts).Where(`"parent_id"=?`, parentID)

	if query != "" {
		qry := "%" + query + "%"
		q = q.Where(`"name" ILIKE ? OR "fullname" ILIKE ?`, qry, qry)
	}

	if err := q.Select(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}
	return districts, nil
}

func GetDistrict(id string) (*District, error) {
	district := District{ID: id}
	if err := db.PG().Model(&district).WherePK().
		Relation(`Parent`).Relation(`Parent.Parent`).Select(); err != nil {
		if err != db.ErrNoRows {
			log.Errorf("DB Error: %v", err)
			return nil, err
		}
		return nil, nil
	}
	return &district, nil
}
