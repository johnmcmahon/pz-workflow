package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/venicegeo/pz-gocommon"
	"gopkg.in/olivere/elastic.v3"
	"log"
	"net/http"
	"os"
	"strings"
)

//---------------------------------------------------------------------------

var esClient *elastic.Client

func runAlertServer(discoveryURL string, port string) error {

	conditionDB := newConditionDB()
	eventDB := newEventDB()
	alertDB := newAlertDB()

	myAddress := fmt.Sprintf("%s:%s", "localhost", port)
	myURL := fmt.Sprintf("http://%s/alerts", myAddress)

	piazza.RegistryInit(discoveryURL)
	err := piazza.RegisterService("pz-alerter", "core-service", myURL)
	if err != nil {
		return err
	}

	///////////////////////////////////////////////////////
	// Create a client
	esClient, err = elastic.NewClient()
	if err != nil {
		// Handle error
		panic(err)
	}

	// Delete the index, just in case
	_, err = esClient.DeleteIndex("events").Do()
	if err != nil {
		// Handle error
		panic(err)
	}

	// Create the index
	_, err = esClient.CreateIndex("events").Do()
	if err != nil {
		// Handle error
		panic(err)
	}
	///////////////////////////////////////////////////////////

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	//router.Use(gin.Logger())
	//router.Use(gin.Recovery())

	//---------------------------------

	router.POST("/events", func(c *gin.Context) {
		var event Event
		err := c.BindJSON(&event)
		if err != nil {
			log.Println(err)
			c.Error(err)
			return
		}
		err = eventDB.write(&event)
		if err != nil {
			log.Println("bbbb")
			c.Error(err)
			return
		}
		log.Println("cccc")
		c.JSON(http.StatusCreated, gin.H{"id": event.ID})

		alertDB.checkConditions(event, conditionDB)
	})

	router.GET("/events", func(c *gin.Context) {
		c.JSON(http.StatusOK, eventDB.getAll())
	})

	//---------------------------------

	router.GET("/alerts/:id", func(c *gin.Context) {
		id := c.Param("id")
		v := alertDB.getByID(id)
		if v == nil {
			c.JSON(http.StatusNotFound, gin.H{"condition_id": id})
			return
		}
		c.JSON(http.StatusOK, v)
	})

	router.GET("/alerts", func(c *gin.Context) {
		c.JSON(http.StatusOK, alertDB.getAll())
	})

	//---------------------------------
	router.POST("/conditions", func(c *gin.Context) {
		var condition Condition
		err := c.BindJSON(&condition)
		if err != nil {
			c.Error(err)
			return
		}
		err = conditionDB.write(&condition)
		if err != nil {
			c.Error(err)
			return
		}
		c.JSON(http.StatusCreated, gin.H{"id": condition.ID})
	})

	router.PUT("/conditions", func(c *gin.Context) {
		var condition Condition
		err := c.BindJSON(&condition)
		if err != nil {
			c.Error(err)
			return
		}
		ok := conditionDB.update(&condition)
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"id": condition.ID})
			return
		}
		c.JSON(http.StatusOK, gin.H{"id": condition.ID})
	})

	router.GET("/conditions", func(c *gin.Context) {
		c.JSON(http.StatusOK, conditionDB.data)
	})

	router.GET("/conditions/:id", func(c *gin.Context) {
		id := c.Param("id")
		v := conditionDB.readByID(id)
		if v == nil {
			c.JSON(http.StatusNotFound, gin.H{"id": id})
			return
		}
		c.JSON(http.StatusOK, v)
	})

	router.DELETE("/conditions/:id", func(c *gin.Context) {
		id := c.Param("id")
		ok := conditionDB.deleteByID(id)
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"id": id})
			return
		}
		c.JSON(http.StatusOK, nil)
	})

	//---------------------------------

	err = router.Run("localhost:" + port)
	return err
}

func app() int {
	var discoveryURL = flag.String("discovery", "http://localhost:3000", "URL of pz-discovery")
	var port = flag.String("port", "12342", "port number of this pz-alerter")

	flag.Parse()

	log.Printf("starting: discovery=%s, port=%s", *discoveryURL, *port)

	err := runAlertServer(*discoveryURL, *port)
	if err != nil {
		fmt.Print(err)
		return 1
	}

	// not reached
	return 1
}

func main2(cmd string) int {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = strings.Fields("main_tester " + cmd)
	return app()
}

func main() {
	os.Exit(app())
}