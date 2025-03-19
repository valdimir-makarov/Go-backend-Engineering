package order

import (
	"context"
	"errors"
	"log"

	"github.com/valdimir-makarov/Go-backend-Engineering/Complete-Golang-Micro-Service-Project/account"
	"github.com/valdimir-makarov/Go-backend-Engineering/Complete-Golang-Micro-Service-Project/order"
	order "github.com/valdimir-makarov/Go-backend-Engineering/Complete-Golang-Micro-Service-Project/order/pb"
	"golang.org/x/text/message/catalog"
)

type gRPCServer struct {
	service       Service
	accountClient *account.Client
	orderClient   *order.Client
	catalogClient *catalog.Client
}

func NewClient(accountURL string, orderClient string, catalogClient string) error {

	accountClient, err := account.NewClient(accountURL)
	if err != nil {
		accountClient.Close()
		log.Fatalf("error While facing connecting to Client", err)
		return nil
	}
	orderClient, err = order.NewClient(orderClient)
	if err != nil {
		orderClient.Close()
		log.Fatalf("error while creating connection to orderClient", err)

	}
	catalogClient, err = catalog.NewClient(catalogClient)
	if err != nil {

		catalogClient.Close()
		log.Fatalf("error  while creating connecting to Catalog Client", err)
	}

	return err
}
func (service *gRPCServer) PostOrder(
	ctx context.Context,
	orderReq *order.GetOrderRequest,
) (*order.GetOrderResponse, error) {

	// Check if account exists
	_, err := service.accountClient.GetAccounts(ctx, orderReq.AccountId)
	if err != nil {
		log.Println("Account doesn't exist:", err)
		return nil, errors.New("account not found")
	}

	// Collect product IDs from the request
	productIDs := []string{}
	for _, p := range orderReq.Products { // Fix: Access orderReq.Products
		productIDs = append(productIDs, p.ProductId)
	}

	// Fetch product details from the catalog service
	orderedProducts, err := service.catalogClient.GetProducts(ctx, 0, 0, productIDs, "")
	if err != nil {
		log.Println("Failed to get products from catalog:", err)
		return nil, errors.New("failed to fetch product details")
	}

	// Prepare ordered products list
	products := []OrderedProduct{} // Fix: Correct struct name

	for _, p := range orderedProducts {
		product := OrderedProduct{
			ID:          p.ID,
			Quantity:    0, // Will be updated later
			Price:       p.Price,
			Name:        p.Name,
			Description: p.Description,
		}

		// Match request product quantity
		for _, rp := range orderReq.Products { // Fix: Reference correct request data
			if rp.ProductId == p.ID {
				product.Quantity = rp.Quantity
				break
			}
		}

		// Add product only if quantity is not zero
		if product.Quantity > 0 {
			products = append(products, product)
		}

	}

	order, err := service.service.PostOrder(ctx, orderReq.AccountId, products)
	if err != nil {
		log.Println("Error posting order: ", err)
		return nil, errors.New("could not post order")
	}

	// Make response order
	orderProto := &order.Order{
		Id:         order.ID,
		AccountId:  order.AccountID,
		TotalPrice: order.TotalPrice,
		Products:   []*order.Order_OrderProduct{},
	}
	orderProto.CreatedAt, _ = order.CreatedAt.MarshalBinary()
	for _, p := range order.Products {
		orderProto.Products = append(orderProto.Products, &order.Order_OrderProduct{
			Id:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Quantity:    p.Quantity,
		})
	}
	return &order.PostOrderResponse{
		Order: orderProto,
	}, nil

}
