package usecase

import (
	"product-go-api/model"
	"product-go-api/repository"
)

type ProductUsecase struct {
	repository repository.ProductRepository
}

func NewProductUsecase(repository repository.ProductRepository) ProductUsecase {
	return ProductUsecase{
		repository: repository,
	}
}

func (pu *ProductUsecase) GetProducts() ([]model.Product, error) {
	return pu.repository.GetProducts()
}

func (pu *ProductUsecase) CreateProduct(product model.Product) (model.Product, error) {
	productId, err := pu.repository.CreateProduct(product)
	if err != nil {
		return model.Product{}, err
	}
	product.ID = productId
	return product, nil
}

func (pu *ProductUsecase) GetProductById(id_product int) (*model.Product, error) {
	product, err := pu.repository.GetProductById(id_product)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (pu *ProductUsecase) DeleteProduct(id_product int) error {
	err := pu.repository.DeleteProduct(id_product)
	if err != nil {
		return err
	}
	return nil
}

func (pu *ProductUsecase) UpdateProduct(product model.Product) (model.Product, error) {
	updatedProduct, err := pu.repository.UpdateProduct(product)
	if err != nil {
		return model.Product{}, err
	}
	return *updatedProduct, nil
}
