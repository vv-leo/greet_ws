package resp

import "git.exclouds.org/ww/global/model"

type ContactModel struct {
	model.ContactModel
	LastContent MessageWithUserModel `json:"lastContent"`
}
