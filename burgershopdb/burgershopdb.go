package burgershopdb

import (
	"os"
	"log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// OrderItem ...
type OrderItem struct {
	Name string
	Price string
}

// User property Password shall be a hash of the user's password.
type User struct {
	Username string
	Password string
	Name string
	Email string
	Mobile string
	Orders []OrderItem
}

func initializeDBHostURL() string {
	dbhost :=  os.Getenv("MONGODB_URI")
	if dbhost == "" {
		dbhost = "127.0.0.1:27017"
	}
	return dbhost
}

// InsertUser ...
func InsertUser(username string, password string, name string, email string, mobile string) {
	dbhost := initializeDBHostURL()
	session, err := mgo.Dial(dbhost)
	if err != nil {
		log.Fatal("InsertUser mgo.Dial error:")
		panic(err)
	}
	defer session.Close()
	c := session.DB("test").C("user")
	err = c.Insert(&User{username, password, name, email, mobile, []OrderItem{}})
	if err != nil {
		log.Fatal("InsertUser c.Insert error:")
		log.Fatal(err)
	}
}

// InsertOrder ...
func InsertOrder(username string, password string, order OrderItem) {
	dbhost := initializeDBHostURL()
	session, err := mgo.Dial(dbhost)
	if err != nil {
		log.Fatal("InsertOrder mgo.Dial error:")
		panic(err)
	}
	defer session.Close()
	c := session.DB("test").C("user")

	result := User{}
	err = c.Find(bson.M{"username": username, "password": password}).One(&result)
	if err != nil {
		log.Fatal("InsertOrder c.Find error:")
		panic(err)
	}

	result.Orders = append(result.Orders, order)
	err = c.Insert(&result)
	if err != nil {
		log.Fatal("InsertOrder c.Insert error:")
		log.Fatal(err)
	}
}

// InsertOrderItem ...
func InsertOrderItem(authcode string, name string, price string) {
	dbhost := initializeDBHostURL()
	session, err := mgo.Dial(dbhost)
	if err != nil {
		log.Fatal("InsertOrderItem mgo.Dial error:")
		panic(err)
	}
	defer session.Close()
	c := session.DB("test").C("orderItem")
	err = c.Insert(&OrderItem{name, price})
	if err != nil {
		log.Fatal("InsertOrderItem c.Insert error:")
		log.Fatal(err)
	}
}