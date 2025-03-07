package services

import (
	"fmt"
	"strings"
	"web-api/internal/pkg/database"
	"web-api/internal/pkg/models/request"
)

type CartService struct {
	*BaseService
}

var Cart = &CartService{}

func (s *CartService) AddToCart(userID int, cartItem *request.Cartrequest) (*request.Cartrequest, error) {
	// Kết nối cơ sở dữ liệu
	db, err := database.DB1Connection()
	if err != nil {
		fmt.Println("Database connection error:", err)
		return nil, err
	}
	// Kiểm tra dữ liệu đầu vào
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	// Debug: In các giá trị trước khi tính tổng
	fmt.Println("Base Price:", cartItem.Base.Price) // Truy xuất vào trường Price của Base
	fmt.Println("Size Price:", cartItem.Size.Price) // Truy xuất vào trường Price của Size

	// Truy vấn lấy giá Base và Size
	baseAndSizeQuery := `
		SELECT b.price AS base_price, s.price AS size_price
		FROM BaseSizes bs
		JOIN Bases b ON b.id = bs.base_id
		JOIN Sizes s ON s.id = bs.size_id
		WHERE bs.base_id = ? AND bs.size_id = ?
	`

	// Tạo một struct để nhận kết quả từ truy vấn
	type PriceResult struct {
		BasePrice float64 `json:"base_price"`
		SizePrice float64 `json:"size_price"`
	}

	var priceResult PriceResult

	// Truy vấn dữ liệu từ cơ sở dữ liệu và ánh xạ vào priceResult
	err = db.Raw(baseAndSizeQuery, cartItem.BaseID, cartItem.SizeID).Scan(&priceResult).Error
	if err != nil {
		fmt.Println("Error querying base and size prices:", err)
		return nil, err
	}

	// Gán giá trị vào các trường trong cartItem
	cartItem.Base.Price = priceResult.BasePrice
	cartItem.Size.Price = priceResult.SizePrice

	// Kiểm tra và tính toán giá cho Extras
	var extraPrices []request.Extrasrequest
	if cartItem.ExtraIDs != "" {
		// Truy vấn các phụ kiện từ cơ sở dữ liệu
		extraQuery := `SELECT id, name, price FROM Extras WHERE FIND_IN_SET(id, ?)`
		err = db.Raw(extraQuery, cartItem.ExtraIDs).Scan(&extraPrices).Error
		if err != nil {
			fmt.Println("Error querying extra prices:", err)
			return nil, err
		}

		// Cập nhật các phụ kiện vào cartItem
		cartItem.Extras = extraPrices
	}

	// Tính toán tổng tiền cho giỏ hàng
	s.calculateTotalPrice(cartItem) // Tính toán tổng tiền

	// Debug: In giá trị tổng tiền sau khi tính toán
	fmt.Println("Calculated Total Price:", cartItem.Price)

	// Kiểm tra xem giỏ hàng của người dùng có tồn tại hay không, nếu chưa thì tạo mới
	insertQuery := `INSERT INTO Cart (user_id, base_id, size_id, flavor_id, sweetness_id, ice_id, extra_ids, quantity, price) 
              VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	// Debug: In giá trị tổng tiền trước khi thực hiện INSERT
	fmt.Println("Inserting price into database:", cartItem.Price)

	// Thực hiện câu lệnh SQL INSERT
	result := db.Exec(insertQuery, userID, cartItem.BaseID, cartItem.SizeID, cartItem.FlavorID, cartItem.SweetnessID, cartItem.IceID, cartItem.ExtraIDs, cartItem.Quantity, cartItem.Price)
	if result.Error != nil {
		fmt.Println("Query execution error:", result.Error)
		return nil, result.Error
	}

	// Trả về cartItem đã thêm vào
	return cartItem, nil
}

func (s *CartService) calculateTotalPrice(cart *request.Cartrequest) {
	var totalPrice float64

	// Kiểm tra Base Price
	if cart.Base.Price > 0 {
		totalPrice += cart.Base.Price
	} else {
		fmt.Println("Base Price is invalid!")
	}

	// Kiểm tra Size Price
	if cart.Size.Price > 0 {
		totalPrice += cart.Size.Price
	} else {
		fmt.Println("Size Price is invalid!")
	}

	// Kiểm tra các Extras Price
	if len(cart.Extras) > 0 {
		for _, extra := range cart.Extras {
			if extra.Price > 0 {
				totalPrice += extra.Price
			} else {
				fmt.Println("Extra Price is invalid!")
			}
		}
	} else {
		fmt.Println("Extras is empty!")
	}

	// Cập nhật tổng giá trị vào giỏ hàng
	cart.Price = totalPrice

	// Debug: In tổng tiền tính được
	fmt.Printf("Total price after calculation: %.2f\n", cart.Price)
}

func (s *CartService) GetCart(userID int) ([]request.Cartrequest, error) {
	var carts []request.Cartrequest

	// Kết nối cơ sở dữ liệu
	db, err := database.DB1Connection()
	if err != nil {
		fmt.Println("Database connection error:", err)
		return nil, err
	}
	// Kiểm tra dữ liệu đầu vào
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	// Truy vấn lấy thông tin giỏ hàng với dữ liệu chi tiết (JOIN)
	query := `
	SELECT c.id, c.user_id, c.base_id, c.size_id, c.flavor_id, c.sweetness_id, c.ice_id, c.extra_ids, c.quantity, c.price,
	       b.name AS BaseName, b.price AS BasePrice,
	       s.name AS SizeName, s.price AS SizePrice,
	       f.name AS FlavorName,
	       t.name AS SweetnessName,
	       i.name AS IceName,
	       e.id AS ExtraID, e.name AS ExtraName, e.price AS ExtraPrice
	FROM Cart c
	LEFT JOIN Bases b ON c.base_id = b.id
	LEFT JOIN Sizes s ON c.size_id = s.id
	LEFT JOIN Flavors f ON c.flavor_id = f.id
	LEFT JOIN Sweetness t ON c.sweetness_id = t.id
	LEFT JOIN IceLevels i ON c.ice_id = i.id
	LEFT JOIN Extras e ON FIND_IN_SET(e.id, c.extra_ids) > 0
	WHERE c.user_id = ?
	`

	// Thực hiện câu lệnh query để lấy tất cả giỏ hàng
	var cartDetails []request.CartDetails
	err = db.Raw(query, userID).Scan(&cartDetails).Error
	if err != nil {
		fmt.Println("Query execution error:", err)
		return nil, fmt.Errorf("Error retrieving carts: %w", err)
	}

	// Dùng map để lọc các giỏ hàng không bị lặp
	seen := make(map[int]bool)
	for _, cartDetail := range cartDetails {
		// Kiểm tra xem giỏ hàng đã có chưa
		if _, exists := seen[cartDetail.ID]; !exists {
			seen[cartDetail.ID] = true

			var cart request.Cartrequest
			cart.ID = cartDetail.ID
			cart.UserID = cartDetail.UserID
			cart.BaseID = cartDetail.BaseID
			cart.SizeID = cartDetail.SizeID
			cart.FlavorID = cartDetail.FlavorID
			cart.SweetnessID = cartDetail.SweetnessID
			cart.IceID = cartDetail.IceID
			cart.ExtraIDs = cartDetail.ExtraIDs
			cart.Quantity = cartDetail.Quantity
			cart.Price = cartDetail.Price

			cart.Base = request.Basesrequest{
				Id:    cart.BaseID,
				Name:  cartDetail.BaseName,
				Price: cartDetail.BasePrice,
			}
			cart.Size = request.SizesRequest{
				ID:    cart.SizeID,
				Name:  cartDetail.SizeName,
				Price: cartDetail.SizePrice,
			}
			cart.Flavor = request.Flavorsrequest{
				Id:   cart.FlavorID,
				Name: cartDetail.FlavorName,
			}
			cart.Sweetness = request.Sweetnessrequest{
				Id:   cart.SweetnessID,
				Name: cartDetail.SweetnessName,
			}
			cart.Ice = request.IceLevelsrequest{
				Id:   cart.IceID,
				Name: cartDetail.IceName,
			}

			// Lấy danh sách các phụ kiện (Extras)
			var extras []request.Extrasrequest
			if cart.ExtraIDs != "" {
				err = db.Raw("SELECT * FROM Extras WHERE FIND_IN_SET(id, ?)", cart.ExtraIDs).Scan(&extras).Error
				if err != nil {
					fmt.Println("Error fetching extras:", err)
					return nil, fmt.Errorf("Error retrieving extras: %w", err)
				}
				cart.Extras = extras
			}

			// Thêm giỏ hàng vào danh sách
			carts = append(carts, cart)
		}
	}

	// Tính toán tổng giá trị giỏ hàng cho mỗi giỏ
	for i := range carts {
		s.calculateTotalPrice(&carts[i])
	}

	return carts, nil
}

// UpdateCart cập nhật thông tin giỏ hàng
func (s *CartService) UpdateCart(userID int, updatedCartItem *request.Cartrequest) (*request.Cartrequest, error) {
	// Kết nối cơ sở dữ liệu
	db, err := database.DB1Connection()
	if err != nil {
		fmt.Println("Database connection error:", err)
		return nil, err
	}
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	// Truy vấn lấy giá Base và Size
	baseAndSizeQuery := `
        SELECT b.price AS base_price, s.price AS size_price
        FROM BaseSizes bs
        JOIN Bases b ON b.id = bs.base_id
        JOIN Sizes s ON s.id = bs.size_id
        WHERE bs.base_id = ? AND bs.size_id = ?
    `
	var price struct {
		BasePrice float64 `json:"base_price"`
		SizePrice float64 `json:"size_price"`
	}

	err = db.Raw(baseAndSizeQuery, updatedCartItem.BaseID, updatedCartItem.SizeID).Scan(&price).Error
	if err != nil {
		fmt.Println("Error querying base and size prices:", err)
		return nil, err
	}

	// Cập nhật giá Base và Size trong cartItem
	updatedCartItem.Base.Price = price.BasePrice
	updatedCartItem.Size.Price = price.SizePrice

	// Truy vấn ExtraIDs từ bảng Cart
	var extraIDs string
	extraIDQuery := `SELECT extra_ids FROM Cart WHERE user_id = ? AND id = ?`
	err = db.Raw(extraIDQuery, userID, updatedCartItem.ID).Scan(&extraIDs).Error
	if err != nil {
		fmt.Println("Error querying extra IDs from Cart:", err)
		return nil, err
	}

	// Tách ExtraIDs thành danh sách các ID
	var extraIDList []string
	if extraIDs != "" {
		extraIDList = strings.Split(extraIDs, ",")
		fmt.Println("Extra ID List:", extraIDList) // Debugging line
	}

	// Truy vấn các phụ kiện từ cơ sở dữ liệu
	var extraPrices []request.Extrasrequest
	if len(extraIDList) > 0 {
		// Truy vấn các phụ kiện từ bảng Extras
		extraQuery := `SELECT id, name, price FROM Extras WHERE id IN (?)`
		err = db.Raw(extraQuery, extraIDList).Scan(&extraPrices).Error
		if err != nil {
			fmt.Println("Error querying extra prices:", err)
			return nil, err
		}

		// Debug: In danh sách giá phụ kiện
		fmt.Println("Extras from DB:", extraPrices)

		// Cập nhật các phụ kiện vào cartItem
		updatedCartItem.Extras = extraPrices
	}

	// Tính toán giá trị tổng của giỏ hàng (base + size + extras)
	totalExtraPrice := 0.0
	for _, extra := range updatedCartItem.Extras {
		totalExtraPrice += extra.Price
	}

	updatedCartItem.Price = updatedCartItem.Base.Price + updatedCartItem.Size.Price + totalExtraPrice

	// Cập nhật giỏ hàng trong cơ sở dữ liệu
	updateQuery := `
        UPDATE Cart 
        SET quantity = ?, price = ? 
        WHERE user_id = ? AND id = ?
    `
	result := db.Exec(updateQuery, updatedCartItem.Quantity, updatedCartItem.Price, userID, updatedCartItem.ID)
	if result.Error != nil {
		fmt.Println("Error updating cart item:", result.Error)
		return nil, result.Error
	}

	// Trả về cartItem đã cập nhật
	return updatedCartItem, nil
}

// RemoveFromCart xóa một mục khỏi giỏ hàng
func (s *CartService) RemoveFromCart(userID int, cartItemID int) error {
	// Kết nối cơ sở dữ liệu
	db, err := database.DB1Connection()
	if err != nil {
		fmt.Println("Database connection error:", err)
		return err
	}
	// Kiểm tra dữ liệu đầu vào
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	// Xóa một mục khỏi giỏ hàng
	query := `DELETE FROM Cart WHERE user_id = ? AND id = ?`

	// Thực thi câu lệnh SQL DELETE
	result := db.Exec(query, userID, cartItemID)
	if result.Error != nil {
		fmt.Println("Error removing cart item:", result.Error)
		return result.Error
	}

	return nil
}
