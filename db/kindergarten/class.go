package kindergarten

import (
	"fmt"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ykgk/kgo/db"
)

func GetKindergartenClassByName(kindergartenID int64, name string) (*KindergartenClass, error) {
	class := KindergartenClass{}
	if err := db.PG().Model(&class).Where(`"kindergarten_id"=?`, kindergartenID).Where(`"name"=?`, name).Where(`NOT "deleted"`).Limit(1).Select(); err != nil {
		if err != db.ErrNoRows {
			log.Errorf("DB Error: %v", err)
			return nil, err
		}
		return nil, nil
	}
	return &class, nil
}

func CreateKindergartenClassTx(tx *pg.Tx, name, remark string, kindergartenID int64) (*KindergartenClass, error) {
	kindergarten := Kindergarten{ID: kindergartenID}
	if err := tx.Model(&kindergarten).WherePK().Select(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}

	class := KindergartenClass{}
	if err := db.PG().Model(&class).Where(`"kindergarten_id"=?`, kindergarten.ID).Where(`"name"=?`, name).Where(`NOT "deleted"`).Limit(1).Select(); err != nil {
		if err != db.ErrNoRows {
			log.Errorf("DB Error: %v", err)
			return nil, err
		}
	} else {
		return &class, nil
	}

	class = KindergartenClass{
		Name:           name,
		Remark:         remark,
		KindergartenID: kindergarten.ID,
	}

	if _, err := tx.Model(&class).Insert(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}

	return &class, nil
}

/* 创建班级 */
func CreateKindergartenClass(name, remark string, kindergartenID int64) (*KindergartenClass, error) {
	tx, err := db.Begin()
	if err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}
	defer tx.Rollback()

	class, err := CreateKindergartenClassTx(tx, name, remark, kindergartenID)
	if err != nil {
		log.Errorf("CreateKindergartenClassTx Error: %v", err)
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}
	return class, nil
}

/* 查询班级 */
func FindKindergartenClasses(query string, kindergartenID int64, page, pageSize int) ([]*KindergartenClass, int, error) {
	classes := make([]*KindergartenClass, 0)

	q := db.PG().Model(&classes).Relation(`Kindergarten`).Where(`NOT "kindergarten_class"."deleted"`)

	if query != "" {
		qry := fmt.Sprintf("%%%s%%", query)
		q = q.WhereGroup(func(q *orm.Query) (*orm.Query, error) {
			q = q.WhereOr(`"kindergarten_class"."name" ILIKE ?`, qry)
			q = q.WhereOr(`"kindergarten_class"."remark" ILIKE ?`, qry)
			return q, nil
		})
	}
	if kindergartenID > 0 {
		q = q.Where(`"kindergarten_class"."kindergarten_id"=?`, kindergartenID)
	}

	if page > 0 && pageSize > 0 {
		q = q.Limit(pageSize).Offset((page - 1) * pageSize)
	}

	total, err := q.Order(`kindergarten_class.id DESC`).SelectAndCountEstimate(100000)
	if err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, 0, err
	}
	return classes, total, nil
}

func UpdateKindergartenClass(id int64, name, remark string) (*KindergartenClass, error) {
	tx, err := db.Begin()
	if err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}
	defer tx.Rollback()

	class := KindergartenClass{ID: id}

	if err := tx.Model(&class).WherePK().Where(`NOT "deleted"`).Select(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}

	if class.Name != name {
		old := KindergartenClass{}
		if err := tx.Model(&old).Where(`"kindergarten_id"=?`, class.KindergartenID).Where(`"name"=?`, name).Where(`NOT "deleted"`).Limit(1).Select(); err != nil {
			if err != db.ErrNoRows {
				log.Errorf("DB Error: %v", err)
				return nil, err
			}
		} else {
			return &old, nil
		}
	}

	if _, err := tx.Model(&class).WherePK().Set(`"name"=?`, name).Set(`"remark"=?`, remark).
		Set(`"updated_at"=CURRENT_TIMESTAMP`).Returning(`*`).Update(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}
	return &class, nil
}

func GetKindergartenClass(id int64) (*KindergartenClass, error) {
	class := KindergartenClass{ID: id}

	if err := db.PG().Model(&class).WherePK().Select(); err != nil {
		if err != db.ErrNoRows {
			log.Errorf("DB Error: %v", err)
			return nil, err
		}
		return nil, nil
	}
	return &class, nil
}

/* 删除班级 */
func DeleteKindergartenClassTx(tx *pg.Tx, id int64) (*KindergartenClass, error) {
	class := KindergartenClass{ID: id}

	if err := tx.Model(&class).WherePK().Where(`NOT "deleted"`).Select(); err != nil {
		if err != db.ErrNoRows {
			log.Errorf("DB Error: %v", err)
			return nil, err
		}
		return nil, nil
	}

	if _, err := tx.Model(&class).WherePK().Set(`"deleted"=?`, true).Returning(`*`).Update(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}
	return &class, nil
}

func DeleteKindergartenClass(id int64) (*KindergartenClass, error) {
	tx, err := db.Begin()
	if err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}
	defer tx.Rollback()

	class, err := DeleteKindergartenClassTx(tx, id)
	if err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		log.Errorf("DB Error: %v", err)
		return nil, err
	}

	return class, nil
}
