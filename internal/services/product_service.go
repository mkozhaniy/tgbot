package services

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/tgbot/internal/domain/interfaces"
	"github.com/tgbot/internal/domain/models"
)

func ParseStringToProduct(str string) (*models.Product, error) {
	str = strings.Trim(str, " ")
	tokens := strings.Split(str, ";")
	if len(tokens) < 6 {
		return nil, fmt.Errorf("Неверное количество аргументов")
	}
	weight, err := strconv.ParseFloat(tokens[1], 32)
	if err != nil {
		return nil, err
	}
	cost, err := strconv.ParseFloat(tokens[2], 32)
	if err != nil {
		return nil, err
	}
	amount, err := strconv.ParseUint(tokens[3], 10, 32)
	if err != nil {
		return nil, err
	}
	res := models.Product{Name: tokens[0],
		Weight: float32(weight), Cost: float32(cost), Amount: uint(amount),
		Kind: tokens[4], Description: tokens[5]}
	return &res, err
}

func SaveProduct(source interfaces.ProductStorage,
	args string) (*models.Product, error) {
	var product *models.Product
	product, err := ParseStringToProduct(args)
	if err != nil {
		return nil, err
	}
	product1, err := source.Save(*product)
	if err != nil {
		return nil, fmt.Errorf("Product with args %s not saved", args)
	}
	*product = product1.(models.Product)
	return product, nil
}

func AddPhotosByName(source interfaces.ProductStorage, name string, photo string) error {
	product, err := source.GetPruductByName(name)
	if err != nil {
		return fmt.Errorf("Product with name %s not found", name)
	}
	product.Photo = append(product.Photo, photo)
	source.Save(*product)
	return nil
}

func GetProducts(product_storage interfaces.ProductStorage) map[string][]models.Product {
	kinds, err := product_storage.GetAllKinds()
	if err != nil {
		panic(err)
	}
	result := make(map[string][]models.Product)
	for _, val := range kinds {
		result[val], err = product_storage.GetProductsByKind(val)
	}
	return result
}

func GetProductsByKind(source interfaces.ProductStorage, kind string) ([]models.Product, error) {
	products, err := source.GetProductsByKind(kind)
	if err != nil {
		return nil, fmt.Errorf("Products with kind %s not found", kind)
	}
	return products, nil
}

func DeleteProduct(source interfaces.ProductStorage, name string) (*models.Product, error) {
	product, err := source.GetPruductByName(name)
	if err != nil {
		return nil, fmt.Errorf("Product with name %s not found", name)
	}
	_, err = source.Delete(*product)
	if err != nil {
		return nil, fmt.Errorf("Product with name %s cannot delete", name)
	}
	return product, nil
}
