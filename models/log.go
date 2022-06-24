package models

import (
	"github.com/go-playground/validator"
)

type LogCustom struct {
	Id             int    `json:"-" gorm:"primaryKey, index"`
	Server         string `json:"server"`
	Time           string `json:"time"`
	UserName       string `json:"user_name"`
	IPAddress      string `json:"ip"`
	Method         string `json:"method"`
	Status         string `json:"status"`
	NameStatus     string `json:"names_tatus"`
	EndPoint       string `json:"end_point"`
	DataRequest    string `json:"data_request"`
	DataBeforeEdit string `json:"data_before_edit"`
	InfoResponse   string `json:"info_response"`
	InfoRequest    string `json:"info_request"`
}

func (LogCustom) TableName() string {
	return "tbl_log"
}

type ErrorResponseLogCustom struct {
	Field string
	Tag   string
	Value string
}

func ValidateLogCustom(logCustom LogCustom) []*ErrorResponseLogCustom {
	var errors []*ErrorResponseLogCustom
	validate := validator.New()
	err := validate.Struct(logCustom)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponseLogCustom
			element.Field = err.StructField()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}
