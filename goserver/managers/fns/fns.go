package fns

import (
	"fmt"
)

type Item struct {
	UL *ULInfo `json:"ЮЛ,omitempty"`
	IP *IPInfo `json:"ИП,omitempty"`
}

func (i *Item) Title() string {
	if i.UL != nil {
		return i.UL.Title()
	}

	if i.IP != nil {
		return i.IP.Title()
	}

	return ""
}

type OrgIDs struct {
	INN  *string `json:"ИНН,omitempty"`
	OGRN *string `json:"ОГРН,omitempty"`
}

type BaseItemInfo struct {
	OrgIDs
	Address          *string `json:"АдресПолн,omitempty"`
	KindOfActivity   *string `json:"ОснВидДеят,omitempty"`
	FoundAt          *string `json:"ГдеНайдено,omitempty"`
	Status           *string `json:"Статус,omitempty"`
	RegistrationDate *string `json:"ДатаРег,omitempty"`
	StopDate         *string `json:"ДатаПрекр,omitempty"`
}

type IPInfo struct {
	BaseItemInfo
	FIO *string `json:"ФИОПолн,omitempty"`
}

func (ip *IPInfo) Title() string {
	return *ip.FIO
}

type ULInfo struct {
	BaseItemInfo
	FullTitle  *string `json:"НаимПолнЮЛ,omitempty"`
	ShortTitle *string `json:"НаимСокрЮЛ,omitempty"`
}

func (ul *ULInfo) Title() string {
	return *ul.FullTitle
}

type OrganizationInfo struct {
	Items []*Item `json:"items,omitempty"`
	Count *int    `json:"Count,omitempty"`
}

type Error struct {
	Type string
	Data string
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Data)
}

func NewError(data string) *Error {
	return &Error{
		Type: "fns error",
		Data: data,
	}
}
