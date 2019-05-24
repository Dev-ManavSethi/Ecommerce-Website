package main

import (
	"text/template"

	"github.com/DevManavSethi/EcommerceWebsite/service"
)

var Tpl *template.Template

type ProductToSend struct {
	Product      *service.VariantProduct
	SuperProduct *service.Product
	AddedToCart  bool
}


