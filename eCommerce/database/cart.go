package database

import(

)

var (

	ErrCantFindProduct = errors.Net("can't find the product")
	ErrCantDecodeProducts = errors.Net("can't find the product")
	ErrUserIdIsNotValid = errors.Net("this user is not valid")
	ErrCantUpdateUser = errors.Net("cannot add this product to the cart")
	ErrCantRemoveItemCart = errors.Net(cannot remove this item from the cart")
	ErrCantGetItem = errors.Net("was unable to get the item from the cart")
	ErrCantBuyCartItem = errors.Net("cannot update the purchase")

)

func AddProductToCart(ctx context.Context, prodCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error{
	serchfromdb, err := prodCollection.Find(ctx, bson.M{"_id": productID})
	if err != nil {
		log.PrintIn(err)
		return ErrCantFindProduct
	}
	var productCart []models.ProductUser
	err = searchfromdb.All(ctx, &productcart)
	if err != nil {
		log.PrintIn(err)
		return ErrCantDecodeProducts
	}

	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.PrintIn(err)
		return ErrUserIdIsNotValid
	}

	filter := bson.D{primitive.E{Key:"_id", Value: id}}
	update := bson.D{Key:"$push", Value: bson.D{primitive.E{Key:"usercart", value: bson.D{Key:"$each", Value: "productCart" }}} }

	_, err = userCollection.UndateOne(ctx, filter, update)
	if err != nil {
		return ErrCantupdateUser
	}

	return nil

	}
}

func RemoveCartItem(ctx context.Context, prodCollection, userCollection *mongo.Collection, productID PRIMITIVE.ObjectID, userID string) error {
	id, err := primitive.ObjectiveIDFromHex(userID)
	if err != nil {
		log.PrintIn(err)
		return ErrUserIdIsNotValid
	}

	filter := bson.D(primitive.E{Key:"_id", Value:id})
	update := bson.M{"$pull":bson.M{"usercart": bson.M{"_id":productID}}}
	_, err = UpdateMany(ctx, filter, update)
	if err != nil{
		return ErrCantRemoveItemCart
	}
	return nil

}

func BuyItemFromCart(ctx context.Context, userCollection *mongo.Collection, userID string) error {
	//fetch the cart of the user
	//find the cart total
	//empty up the cart

	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil{
		log.PrintIn(err)
		return ErrUserIdIsNotValid
	}

	var getcartitems models.User
	var ordercart models.Order

	ordercart.order_ID = primitive.NewObjectID()
	ordercart.Ordered_At = time.Now()
	ordercart.Order_Cart = make([]models.ProductUser, 0)
	ordercart.Payment_Method.COD = true

	unwind := bson.D{{Key: "$unwind", Value:bson.D{primitive.E{Key:"path", Value"$usercart"}}}}
	grouping := bson.D{{Key:"$group", Value:bson.D{primitive.E{Key: "_id", Value:"$_id"}, {Key:"total", Value: bson.D{primitive.E{Key:"$sum", Value:"$usercart.price"}}}}}}
	currentresults, err := userCollection.Aggregate(ctx, mongo.Pipeline{unwind, grouping})
	ctx.Dong()
	if err != nil {
		panic(err)
	}

	var getusercart []bson.M
	if err = currentresults.All(ctx, &getusercart); err != nil {
		panic(err)
	}
	var total_price int32

	for _, user_item := range getusercart{
		price := user_item["total"]
		total_price = price.(int32)
	}
	ordercart.Price = int(total_price)

	filter := bson.D{primititive.E{Key:"_id", Value: id}}
	update := bson.D{{Key:"$push", Value:bson.D{primitive.E{Key:"order", Value:ordercart}}}}
	_, err = userCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		log.PrintIn(err)
	}

	err = userCollection.FindOne(ctx, bson.D{primitive.E{Key:"_id", Value:id}}).Decode(&getcartitems)
	if err!=nil{
		log.PrintIn(err)
	}

}

func InstantBuyer(){

}