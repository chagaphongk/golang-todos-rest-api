package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/viper"
)

type todo struct {
	ID        bson.ObjectId `json:"id" bson:"_id"`
	Topic     string        `json:"topic" bson:"topic"`
	Done      bool          `json:"done" bson:"done"`
	CreatedAt time.Time     `json:"createdAt" bson:"createdAt"`
}

type handler struct {
	m *mgo.Session
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	mongoUser := viper.GetString("mongo.user")             //"mongodb"
	mongoHost := viper.GetString("mongo.host")             //"127.0.0.1"
	mongoPort := viper.GetString("mongo.port")             //"27017"
	mongoCollection := viper.GetString("mongo.collection") //"/local"

	connString := fmt.Sprintf("%v://%v:%v/%v", mongoUser, mongoHost, mongoPort, mongoCollection)
	session, err := mgo.Dial(connString)
	if err != nil {
		e.Logger.Fatal(err)
		return
	}

	h := &handler{
		m: session,
	}

	e.GET("/todos", h.list)
	e.GET("/todos/:id", h.view)
	e.POST("/todos", h.create)
	e.PUT("todos/:id", h.done)
	e.DELETE("todos/:id", h.delete)

	e.Logger.Fatal(e.Start(":1323"))
}

func (h *handler) create(c echo.Context) error {
	session := h.m.Copy()
	defer session.Close()

	var t todo
	if err := c.Bind(&t); err != nil {
		return err
	}

	t.ID = bson.NewObjectId()
	t.CreatedAt = time.Now()

	col := session.DB("local").C("todos")
	if err := col.Insert(t); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, t)
}

func (h *handler) list(c echo.Context) error {
	session := h.m.Copy()
	defer session.Close()

	var ts []todo

	col := session.DB("local").C("todos")
	if err := col.Find(nil).All(&ts); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, ts)
}

func (h *handler) view(c echo.Context) error {
	session := h.m.Copy()
	defer session.Close()

	id := bson.ObjectIdHex(c.Param("id"))

	var ts todo

	col := session.DB("local").C("todos")
	if err := col.FindId(id).One(&ts); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, ts)
}

func (h *handler) done(c echo.Context) error {
	session := h.m.Copy()
	defer session.Close()

	id := bson.ObjectIdHex(c.Param("id"))
	// topic := c.("topic")

	var ts todo

	col := session.DB("local").C("todos")
	if err := col.FindId(id).One(&ts); err != nil {
		return err
	}

	ts.Done = true
	// ts.Topic = topic
	if err := col.UpdateId(id, ts); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, ts)
}

func (h *handler) delete(c echo.Context) error {
	session := h.m.Copy()
	defer session.Close()

	id := bson.ObjectIdHex(c.Param("id"))
	col := session.DB("local").C("todos")
	if err := col.RemoveId(id); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"result": "success",
	})
}
