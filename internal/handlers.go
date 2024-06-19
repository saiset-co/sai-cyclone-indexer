package internal

import (
	saiService "github.com/saiset-co/sai-service/service"
)

func (is *InternalService) NewHandler() saiService.Handler {
	return saiService.Handler{
		"add_address": saiService.HandlerElement{
			Name:        "add_address",
			Description: "Add new address for scan transactions",
			Function: func(data, meta interface{}) (interface{}, int, error) {
				return is.addAddress(data)
			},
		},
		"delete_address": saiService.HandlerElement{
			Name:        "delete_address",
			Description: "Delete address from addresses list",
			Function: func(data, meta interface{}) (interface{}, int, error) {
				return is.deleteAddress(data)
			},
		},
	}
}
