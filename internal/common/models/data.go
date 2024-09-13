package models

type DataType int

const (
	DataTypeUnknown DataType = iota
	DataTypeCredentials
	DataTypeText
	DataTypeBinary
	DataTypeBankCard
)

type DataInfo struct {
	ID          uint     `json:"id"`
	Type        DataType `json:"type" validate:"required"`
	Description string   `json:"description"`
	Value       string   `json:"value" validate:"required"`
}
