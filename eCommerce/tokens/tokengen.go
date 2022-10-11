package tokens

import(
	"os"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/form3tech-oss/jwt-go"
)

type SignedDetails struct{
	Email string
	First_Name string
	Last_Name string
	Uid string
	jwt.StandardClaims
}

var UserData *mongo.Collection = database.userData(database.Client, "Users")

var SECRET_KEY = os.Getenv("SECRET_KEY")

func TokenGenerator(email string, firstname string, lastname string, uid string) (signedtoken string, signedrefreshtoken string, err error) {

	claims := &SignedDetails{
		Email: email,
		First_Name: firstname,
		Last_Name: lastname,
		Uid: uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshclaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))

	if err != nil{
		return "", "", err
	}

	refreshtoken, err := jwt.NewWithClaims(jwt.SigningMethodHS384, refreshclaims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Panic(err)
		return
	}

	return token, refreshtoken, err
}

func ValidateToken(signedtoken string) (claims *SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(signedtoken, &SignedDetails{}, func(token *jwt.Token)(interface{}, error) {
		return []byte(SECRET_KEY), nil
	})

	if err != nil {
		msg = err.Error
		return
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = "the token in invalid"
		return

	}

	if claims.ExpiresAt < time.Now().Local().unix(){
		msg = "token is already expired"
		return
	}
	return claims, msg

}

func UpdateAllTokens(signedtoken string, signedrefreshtoken string, userid string ){

	var ctx, cancel = context.WithTimeOut(context.Backgroud(), 100*time.Second)

	var updateobj primitive.DeleteAddress

	updateobj = append (updateobj,bson.E{Key:"token", Value: signedtoken})
	updateobj = append (updateobj.bson.E{Key:"refresh_token", Value: signedrefreshtoken})
	updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	updateobj = append (updateobj,bson.E{Key:"updatedat", Value: updated_at})

	upsert := true

	filter := bson.M{"user_id": userid}
	opt := options.UpdateOptions{
		upsert: &upsert,
	}

	_, err := UserData.UpdateOne(ctx, filter, bson.D{
		{Key:"$set", Value: updateobj},
	},
	&opt)
	defer cancel()

	if err != nil {
		log.Panic(err)
		return
	}

}