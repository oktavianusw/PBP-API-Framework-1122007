package controllers

import (
	"week2/app_revel/helpers"
	"week2/app_revel/models"

	"github.com/revel/revel"
)

type Product struct {
	*revel.Controller
}

var (
	product models.Product
	errors  helpers.Error
)

func (c Product) ListProducts() revel.Result {
	data, response := product.ListProducts()
	if response != nil {
		c.Response.Status = 400
		errors.Error = "Failed to get products"
		return c.RenderJSON(errors)
	}
	c.Response.Status = 200
	return c.RenderJSON(data)
}

func (c Product) CreateProduct() revel.Result {
	err := c.Params.BindJSON(&product)
	if err != nil {
		c.Response.Status = 400
		errors.Error = "Invalid JSON passed"
		return c.RenderJSON(errors)
	}
	response := product.AddProduct()
	if response != nil {
		c.Response.Status = 503
		errors.Error = "Failed to insert product"
		return c.RenderJSON(errors)
	}
	c.Response.Status = 201
	return c.RenderJSON(product)
}

func (c Product) GetProduct(id int64) revel.Result {
	data, response := product.GetProduct(id)
	if response != nil {
		c.Response.Status = 400
		errors.Error = "Product not found"
		return c.RenderJSON(errors)
	}
	c.Response.Status = 200
	return c.RenderJSON(data)
}

func (c Product) UpdateProduct(id int64) revel.Result {
	err := c.Params.BindJSON(&product)
	if err != nil {
		c.Response.Status = 400
		errors.Error = "Invalid JSON passed"
		return c.RenderJSON(errors)
	}
	response := product.UpdateProduct(id)
	if response != nil {
		c.Response.Status = 503
		errors.Error = "Failed to update product"
		return c.RenderJSON(errors)
	}
	c.Response.Status = 200
	return c.RenderJSON(product)
}

func (c Product) DeleteProduct(id int64) revel.Result {
	response := product.DeleteProduct(id)
	if response != nil {
		c.Response.Status = 503
		errors.Error = "Failed to delete product"
		return c.RenderJSON(errors)
	}
	c.Response.Status = 200
	return nil
}
