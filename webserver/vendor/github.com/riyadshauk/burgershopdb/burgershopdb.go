package burgershopdb

import (
	"fmt"
	"regexp"
	"os"
	"log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"encoding/json"
)

// OrderItem ...
type OrderItem struct {
	Name string
	Price string
}

// OrderRequest ...
type OrderRequest struct {
	// ID string
	Username string
	Password string
	OrderItems [][]OrderItem
}

// User property Password shall be a hash of the user's password.
type User struct {
	Username string
	Password string
	Name string
	Email string
	Mobile string
	Orders [][]OrderItem
}

func initializeDBHostURL() (dbhost string, mongoDBUser string) {
	dbhost =  os.Getenv("MONGODB_URI")
	if dbhost != "" {
		re := regexp.MustCompile("mongodb://(.*):.*@")
		mongoDBUser = re.FindAllStringSubmatch(dbhost, 1)[0][1] // @todo handle fatal errors
	} else {
		dbhost = "127.0.0.1:27017"
		mongoDBUser = "test"
	}
	fmt.Printf("Using dbhost: %s", dbhost)
	return
}

// InsertUser ...
func InsertUser(username string, password string, name string, email string, mobile string) {
	dbhost, mongoDBUser := initializeDBHostURL()
	session, err := mgo.Dial(dbhost)
	if err != nil {
		log.Fatal("InsertUser mgo.Dial error:")
		panic(err)
	}
	defer session.Close()
	c := session.DB(mongoDBUser).C("user")
	err = c.Insert(&User{username, password, name, email, mobile, [][]OrderItem{}})
	if err != nil {
		log.Fatal("InsertUser c.Insert error:")
		log.Fatal(err)
	}
}
// GetUser ...
func GetUser(username string, password string) []byte {
	dbhost, mongoDBUser := initializeDBHostURL()
	session, err := mgo.Dial(dbhost)
	if err != nil {
		log.Fatal("InsertUser mgo.Dial error:")
		panic(err)
	}
	defer session.Close()
	c := session.DB(mongoDBUser).C("user")
	var oneUser User // []interface{} // it's not []User because User has capitalized vars...
	err = c.Find(bson.M{"username": username, "password": password}).One(&oneUser)
	if err != nil {
		log.Fatal("GetUser error:")
		panic(err)
	}
	userAsJSON, e := json.Marshal(oneUser)
	if e != nil {
		log.Fatal("error:")
		panic(e)
	}
	return userAsJSON
}

// GetOrders ...
func GetOrders(username string, password string) []byte {
	dbhost, mongoDBUser := initializeDBHostURL()
	session, err := mgo.Dial(dbhost)
	if err != nil {
		log.Fatal("GetOrders mgo.Dia error:")
		panic(err)
	}
	defer session.Close()
	c := session.DB(mongoDBUser).C("user")
	var user User
	err = c.Find(bson.M{"username": username, "password": password}).One(&user)
	if err != nil {
		log.Fatal("GetUser error:")
		panic(err)
	}
	ordersAsJSON, e := json.Marshal(user.Orders)
	if e != nil {
		log.Fatal("error:")
		panic(e)
	}
	return ordersAsJSON
}

// InsertOrder ...
func InsertOrder(orderRequest OrderRequest) {
	dbhost, mongoDBUser := initializeDBHostURL()
	session, err := mgo.Dial(dbhost)
	if err != nil {
		log.Fatal("InsertOrder mgo.Dial error:")
		panic(err)
	}
	defer session.Close()
	c := session.DB(mongoDBUser).C("user")

	result := User{}
	err = c.Find(bson.M{"username": orderRequest.Username, "password": orderRequest.Password}).One(&result)
	if err != nil {
		log.Fatal("InsertOrder c.Find error:")
		panic(err)
	}

	result.Orders = append(result.Orders, orderRequest.OrderItems)
	// @todo update instead of insert new record into db... (try using _id??)
	// @todo https://godoc.org/gopkg.in/mgo.v2#Collection.Update
	// err = c.Insert(&result)
	err = c.Update(bson.M{"username": orderRequest.Username, "password": orderRequest.Password}, &result)
	if err != nil {
		log.Fatal("InsertOrder c.Update error:")
		log.Fatal(err)
	}
}

// InsertOrderItem ...
func InsertOrderItem(authcode string, name string, price string) {
	dbhost, mongoDBUser := initializeDBHostURL()
	session, err := mgo.Dial(dbhost)
	if err != nil {
		log.Fatal("InsertOrderItem mgo.Dial error:")
		panic(err)
	}
	defer session.Close()
	c := session.DB(mongoDBUser).C("orderItem")
	err = c.Insert(&OrderItem{name, price})
	if err != nil {
		log.Fatal("InsertOrderItem c.Insert error:")
		log.Fatal(err)
	}
}