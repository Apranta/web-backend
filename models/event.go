package models

import (
	basemodel "web-backend-patal/config"
)

type (
	Event struct {
		basemodel.BaseModel
		Nama string `json:"nama" gorm:"column:nama"`
	}
)

func (Event) TableName() string {
	return "Event"
}

func (p *Event) Create() error {
	err := basemodel.Create(&p)
	return err
}

func (p *Event) Save() error {
	err := basemodel.Save(&p)
	return err
}

func (p *Event) Delete() error {
	err := basemodel.Delete(&p)
	return err
}

func (p *Event) FindbyID(id int) error {
	err := basemodel.FindbyID(&p, id)
	return err
}

func (p *Event) SingleFind(filter interface{}) error {
	err := basemodel.SingleFindFilter(&p, filter)
	return err
}

func (p *Event) GetAll(filter interface{}) (result interface{}, err error) {
	dist := []Event{}
	result, err = basemodel.GetAll(&dist, filter)

	return result, err
}

func (b *Event) PagedFilterSearch(page int, rows int, orderby string, sort string, filter interface{}) (result basemodel.PagedFindResult, err error) {
	Question := []Event{}
	orders := []string{orderby}
	sorts := []string{sort}
	result, err = basemodel.PagedFindFilter(&Question, page, rows, orders, sorts, filter, []string{})

	return result, err
}

func (model *Event) FindFilter(order string, sort string, limit int, offset int, filter interface{}) (result interface{}, err error) {
	cat := []Event{}
	orders := []string{order}
	sorts := []string{sort}
	result, err = basemodel.FindFilter(&cat, orders, sorts, limit, offset, filter)
	return result, err
}
