package main

//--go_out=plugins=grpc:.
//e2BbLeLTr8dcSvm
import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"google.golang.org/appengine"

	"github.com/fatih/color"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/DevManavSethi/EcommerceWebsite/service"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"google.golang.org/grpc"
)

//type server struct{}

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile) //to log errors with line number and column of error

	mux := http.NewServeMux()

	mux.HandleFunc("/", home)

	mux.HandleFunc("/login", login)
	mux.HandleFunc("/signup", signup)

	mux.HandleFunc("/editUser", editUser)

	mux.HandleFunc("/search", search)
	mux.HandleFunc("/product", product)wishlist
	mux.HandleFunc("/AddToWishlist", AddToWishlist)

	mux.HandleFunc("/checkout", checkout)
	mux.HandleFunc("/pay", pay)

	mux.HandleFunc("/cart", cart)
	mux.HandleFunc("/AddToCart", addToCart)

	mux.Handle("/favicon.ico", http.NotFoundHandler())

	log.Println("Web Server listening at  localhost:8080")

	appengine.Main()

}

func addToCart(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		err01 := r.ParseForm()
		if err01 != nil {

		}
		ProductID := r.FormValue("id")
		ProductID_int, _ := strconv.Atoi(ProductID)
		quantity := r.FormValue("quantity")
		quantity_int, _ := strconv.Atoi(quantity)
		comment := r.FormValue("comment")

		cookie, err0001 := r.Cookie("ecommerce_user")

		if err0001 != nil || err01 == http.ErrNoCookie {

			//prompt to login
			//if yes, then store in mongo

			//else
			cookie, err02 := r.Cookie("ecommerce_cart")
			if err02 != nil || err02 == http.ErrNoCookie {

				http.SetCookie(w, &http.Cookie{
					Name:  "ecommerce_cart",
					Value: ProductID + ":" + quantity + ":" + comment + "|",
				})
			} else {
				cookie.Value = cookie.Value + ProductID + ":" + quantity + ":" + comment + "|"
				http.SetCookie(w, cookie)
			}

		} else {

			//store in mongo

			userId := cookie.Value

			cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
			if err != nil {
				log.Fatalf("could not connect: %v", err)
			}
			defer cc.Close()

			c := service.NewEcommerceClient(cc)

			request := &service.AddToCartRequest{
				Quantity:  int32(quantity_int),
				Comments:  comment,
				UserID:    userId,
				ProductID: int32(ProductID_int),
			}

			response, err004 := c.AddToCart(context.TODO(), request)
			if err004 != nil {

				fmt.Fprintf(w, "Error adding to cart!")
			}

			color.Red("", response)

			if response.Success {

				Tpl.ExecuteTemplate(w, "cart_success.html", nil)

			} else {

			}

		}
	}
}

func search(w http.ResponseWriter, r *http.Request) {

	keys, ok := r.URL.Query()["query"]

	if !ok || len(keys[0]) < 1 {
		log.Println("Url Param 'key' is missing")

	}

	query := keys[0]

	//search products mongo db

	mongoDBclient, err001 := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err001 != nil {
		FatalOnError("Cannot connect to MongoDBclient", err001)

	}

	regex := bson.M{"name": primitive.Regex{Pattern: query, Options: "i"}}

	resultCursor, err002 := mongoDBclient.Database("ecommerce").Collection("products").Find(context.TODO(), regex)
	if err002 != nil {
		FatalOnError("Error finding from mongo", err002)
	} else {

		var SearchResults []*service.VariantProduct

		for resultCursor.Next(context.TODO()) {

			var result *service.VariantProduct
			err01 := resultCursor.Decode(&result)
			FatalOnError("Error decoding result", err01)
			SearchResults = append(SearchResults, result)

		}

		Tpl.ExecuteTemplate(w, "search.html", SearchResults)

	}

}

func AddToWishlist(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["id"]

	if !ok || len(keys[0]) < 1 {
		log.Println("Url Param 'key' is missing")
		return
	}

	id := keys[0]
	id_int, _ := strconv.Atoi(id)

	//--------------------------------------------------------------------------

	mongoDBclient, err001 := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err001 != nil {
		FatalOnError("Cannot connect to MongoDBclient", err001)
		return
	}

	log.Println(id_int)

	result := mongoDBclient.Database("ecommerce").Collection("products").FindOne(context.TODO(), bson.M{"id": id_int})

	var product *service.VariantProduct

	er01 := result.Decode(&product)

	//	log.Println(product)

	if er01 != nil {
		FatalOnError("Error finding product from database1", er01)
	} else {

		ck, er1 := r.Cookie("ecommerce_user")
		if er1 != nil || er1 == http.ErrNoCookie {

			Tpl.ExecuteTemplate(w, "pls_login.html", nil)
		} else {
			u_id := ck.Value

			//u_id_int, _ := strconv.Atoi(u_id)

			user_res := mongoDBclient.Database("ecommerce").Collection("users").FindOne(context.TODO(), bson.M{"id": u_id})

			var user *service.User

			er2 := user_res.Decode(&user)
			//	log.Println(user)
			if er2 != nil || user == nil {

			} else {

				for _, productInWishlist := range user.WishList {
					if product.ID == productInWishlist.ID {

						//PRODUCT ALREDY PRESENT

					}
				}

				user.WishList = append(user.WishList, product)
				log.Println(user.WishList)

				_, err3 := mongoDBclient.Database("ecommerce").Collection("users").DeleteOne(context.TODO(), bson.M{"id": u_id})
				if err3 != nil {
					log.Println(err3)
					return

				} else {

					_, err4 := mongoDBclient.Database("ecommerce").Collection("users").InsertOne(context.TODO(), *user)
					if err4 != nil {

					}
					Tpl.ExecuteTemplate(w, "AddToWishlist_success.html", nil)

				}
			}
		}

	}

}
func product(w http.ResponseWriter, r *http.Request) {

	keys, ok := r.URL.Query()["id"]

	if !ok || len(keys[0]) < 1 {
		log.Println("Url Param 'key' is missing")
		return
	}

	id := keys[0]
	id_int, _ := strconv.Atoi(id)

	mongoDBclient, err001 := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err001 != nil {
		FatalOnError("Cannot connect to MongoDBclient", err001)
		return
	}

	result := mongoDBclient.Database("ecommerce").Collection("products").FindOne(context.TODO(), bson.M{"id": id_int})

	var product *service.VariantProduct

	er01 := result.Decode(&product)

	if product == nil || er01 != nil {
		FatalOnError("Error finding user from database", er01)
	} else {

		MainProductName := product.GetMainProductName()
		result := mongoDBclient.Database("ecommerce").Collection("main_products").FindOne(context.TODO(), bson.M{"name": MainProductName})

		var MainProduct *service.Product

		er01 := result.Decode(&MainProduct)

		if MainProduct == nil || er01 != nil {
			FatalOnError("Error finding user from database", er01)
		} else {

			var ProductToSend ProductToSend
			ProductToSend.Product = product
			ProductToSend.SuperProduct = MainProduct
			ProductToSend.AddedToCart = true

			err02 := Tpl.ExecuteTemplate(w, "product.html", ProductToSend)
			if err02 != nil {

			}

		}

	}

}

func home(w http.ResponseWriter, r *http.Request) {

	cookie, err01 := r.Cookie("ecommerce_user")
	if err01 == http.ErrNoCookie || err01 != nil {

		var user service.User

		user.FirstName = "Guest"

		err1 := Tpl.ExecuteTemplate(w, "home.html", user)
		if err1 != nil {

		}

		//prompt to login
	} else {

		mongoDBclient, err001 := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
		if err001 != nil {
			FatalOnError("Cannot connect to MongoDBclient", err001)
			return
		}

		result := mongoDBclient.Database("ecommerce").Collection("users").FindOne(context.TODO(), bson.M{"id": cookie.Value})

		var user *service.User

		er01 := result.Decode(&user)

		FatalOnError("Error finding user from database", er01)

		if user == nil {

			err02 := Tpl.ExecuteTemplate(w, "home.html", nil)

			if err02 != nil {
				FatalOnError("Error executing template", err02)
				http.Redirect(w, r, "/", http.StatusNotFound)
			}

		}

		err02 := Tpl.ExecuteTemplate(w, "home.html", user)

		if err02 != nil {
			FatalOnError("Error executing template", err02)
			http.Redirect(w, r, "/", http.StatusNotFound)
		}
	}

}
func login(w http.ResponseWriter, r *http.Request) {

	// if r.Response.StatusCode == http.StatusNonAuthoritativeInfo {
	// 	// write wrong email or password
	// }
	//--------------------------------------------------------------------

	cookie, err0001 := r.Cookie("ecommerce_user")
	if err0001 != nil || err0001 == http.ErrNoCookie {

		if r.Method == http.MethodPost {

			cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
			if err != nil {
				log.Fatalf("could not connect: %v", err)
			}
			defer cc.Close()

			c := service.NewEcommerceClient(cc)

			//-----------------------------------------------------------------------

			err01 := r.ParseForm()
			FatalOnError("Error parsing form. ", err01)

			//---------------------------------------------------------------------

			req := &service.LoginRequest{
				Email:          r.FormValue("email"),
				Password:       r.FormValue("password"),
				ToKeepLoggedIn: true,
			}

			res, err02 := c.Login(context.TODO(), req)
			if err02 != nil {
				FatalOnError("Error logging in: ", err02)

				http.Redirect(w, r, "/login", http.StatusNonAuthoritativeInfo)
			}

			//-----------------------------------------------------------------------

			user := res.GetUser()

			http.SetCookie(w, &http.Cookie{
				Name:  "ecommerce_user",
				Value: user.GetID(),
			})
			log.Println("im here")
			err011 := Tpl.ExecuteTemplate(w, "login_success.html", *user)
			if err011 != nil {
				FatalOnError("Error parsing template login.html", err01)

			}

		} else {

			err01 := Tpl.ExecuteTemplate(w, "login.html", nil)
			if err01 != nil {
				FatalOnError("Error parsing template login.html", err01)

			}
		}

	} else {

		mongoDBclient, err001 := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
		if err001 != nil {
			FatalOnError("Cannot connect to MongoDBclient", err001)
			return
		}

		result := mongoDBclient.Database("ecommerce").Collection("users").FindOne(context.TODO(), bson.M{"id": cookie.Value})

		var user *service.User

		er01 := result.Decode(&user)

		FatalOnError("Error finding user from database", er01)

		err01 := Tpl.ExecuteTemplate(w, "login_user_already_logged_in.html", *user)
		if err01 != nil {
			FatalOnError("Error parsing template login.html", err01)

		}
	}

}
func signup(w http.ResponseWriter, r *http.Request) {

	cookie, err1 := r.Cookie("ecommerce_user")

	if err1 != nil || err1 == http.ErrNoCookie {

		if r.Method == http.MethodPost {

			cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
			if err != nil {
				log.Fatalf("could not connect grpc server: %v", err)
			}
			defer cc.Close()

			c := service.NewEcommerceClient(cc)

			//-----------------------------------------------------------------------

			err01 := r.ParseForm()
			FatalOnError("Error parsing form. ", err01)

			//---------------------------------------------------------------------

			phone := r.FormValue("phone")
			phone_int, _ := strconv.Atoi(phone)

			pin := r.FormValue("pin")
			pin_int, _ := strconv.Atoi(pin)

			req := &service.SignupRequest{
				User: &service.User{
					FirstName: r.FormValue("fname"),
					LastName:  r.FormValue("lname"),
					Phone:     int64(phone_int),
					Email:     r.FormValue("email"),
					Address1:  r.FormValue("add1"),
					Address2:  r.FormValue("add2"),
					City:      r.FormValue("city"),
					State:     r.FormValue("state"),
					Country:   r.FormValue("country"),
					Pincode:   int64(pin_int),
					Cart:      nil,
					Funds:     0,
					Member:    false,
					Password:  r.FormValue("pass"),
				},
			}

			res, err02 := c.Signup(context.TODO(), req)
			if err02 != nil {

				FatalOnError("Error signup ", err02)
				http.Redirect(w, r, "/signup", http.StatusNotImplemented)

			}
			http.SetCookie(w, &http.Cookie{
				Name:  "ecommerce_user",
				Value: res.GetUser().GetID(),
			})

			user := res.GetUser()

			err03 := Tpl.ExecuteTemplate(w, "signup.html", *user)
			if err03 != nil {
				FatalOnError("Error execute template signup.html", err01)
			}

		} else {

			err01 := Tpl.ExecuteTemplate(w, "signup.html", nil)
			if err01 != nil {
				FatalOnError("Error execute template signup.html", err01)
			}
		}

	} else {

		mongoDBclient, err001 := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
		if err001 != nil {
			FatalOnError("Cannot connect to MongoDBclient", err001)
			return
		}

		result := mongoDBclient.Database("ecommerce").Collection("users").FindOne(context.TODO(), bson.M{"id": cookie.Value})

		var user *service.User

		er01 := result.Decode(&user)

		FatalOnError("Error finding user from database", er01)

		err01 := Tpl.ExecuteTemplate(w, "signup_user_already_logged_in.html", *user)
		if err01 != nil {
			FatalOnError("Error execute template signup.html", err01)
		}
	}

}
func editUser(w http.ResponseWriter, r *http.Request) {

}
func cart(w http.ResponseWriter, r *http.Request) {

}
func checkout(w http.ResponseWriter, r *http.Request) {

}
func pay(w http.ResponseWriter, r *http.Request) {

}

//{"name": "Rajdhani Besan", "variants": [{"name": "Rajdhani Besan 500 grams", "categoryid":1, "id":1, "imagepaths":["/images/rajdhanibesan500_1.jpg","/images/rajdhanibesan500_2.jpg","/images/rajdhanibesan500_3.jpg", "/images/rajdhanibesan500_4.jpg", "/images/rajdhanibesan500_5.jpg"],"unitprice":45, "mrp": 55, "cashback":0, "unit": "Kg","size":0.5,"memberprice":43,"mainproductname":"Rajdhani Besan", "xxx_nounkeyedliteral":{}, "xxx_unrecognized": null, "xxx_sizecache":0},{"name": "Rajdhani Besan 1 Kg", "categoryid":1, "id":2, "imagepaths":["/images/rajdhanibesan1000_1.jpg","/images/rajdhanibesan1000_2.jpg","/images/rajdhanibesan1000_3.jpg", "/images/rajdhanibesan1000_4.jpg", "/images/rajdhanibesan1000_5.jpg"],"unitprice":80, "mrp": 105, "cashback":0, "unit": "Kg","size":1,"memberprice":85,"mainproductname":"Rajdhani Besan", "xxx_nounkeyedliteral":{}, "xxx_unrecognized": null, "xxx_sizecache":0}], "xxx_nounkeyedliteral":{}, "xxx_unrecognized": null, "xxx_sizecache":0}
