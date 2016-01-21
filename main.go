package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/venicegeo/pz-gocommon"
	"gopkg.in/olivere/elastic.v3"
	"log"
	"net/http"
	"os"
	"strings"
)

//---------------------------------------------------------------------------

func makeESIndex(client *elastic.Client, index string) error {
	exists, err := client.IndexExists(index).Do()
	if err != nil {
		return err
	}

	if exists {
		_, err = client.DeleteIndex(index).Do()
		if err != nil {
			return err
		}
	}

	_, err = client.CreateIndex(index).Do()
	if err != nil {
		return err
	}
	return nil
}

func newESClient() (*elastic.Client, error) {
  client, err := elastic.NewClient(
    elastic.SetURL("https://search-venice-es-pjebjkdaueu2gukocyccj4r5m4.us-east-1.es.amazonaws.com"),
    elastic.SetSniff(false),
    elastic.SetMaxRetries(5),
    elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
    elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)))
	if err != nil {
		return nil, err
	}
	return client, nil
}

///////////////////////////////////////////////////////////

func runAlertServer(serviceAddress string, discoverAddress string, debug bool) error {

	esClient, err := newESClient()
	if err != nil {
		return err
	}

	conditionDB, err := newConditionDB(esClient, "conditions")
	if err != nil {
		return err
	}
	eventDB, err := newEventDB(esClient, "events")
	if err != nil {
		return err
	}
	alertDB, err := newAlertDB(esClient, "alerts")
	if err != nil {
		return err
	}

/*	myAddress := fmt.Sprintf(":%s", port)
	myURL := fmt.Sprintf("http://%s/alerts", myAddress)

	piazza.RegistryInit(discoveryURL)
	err = piazza.RegisterService("pz-alerter", "core-service", myURL)
	if err != nil {
		return err
	}*/

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	//router.Use(gin.Logger())
	//router.Use(gin.Recovery())

	//---------------------------------

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "Hi. I'm pz-alerter.")
	})

	//---------------------------------

	router.POST("/events", func(c *gin.Context) {
		event := &Event{}
		err := c.BindJSON(event)
		if err != nil {
			log.Println(err)
			c.Error(err)
			return
		}
		err = eventDB.write(event)
		if err != nil {
			c.Error(err)
			return
		}
		c.JSON(http.StatusCreated, gin.H{"id": event.ID})

		alertDB.checkConditions(*event, conditionDB)
	})

	router.GET("/events", func(c *gin.Context) {
		m, err := eventDB.getAll()
		if err != nil {
			c.Error(err)
			return
		}
		c.JSON(http.StatusOK, m)
	})

	//---------------------------------

	router.GET("/alerts/:id", func(c *gin.Context) {
		id := c.Param("id")
		v, err := alertDB.getByConditionID(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"condition_id": id})
			return
		}
		if v == nil {
			c.JSON(http.StatusNotFound, gin.H{"condition_id": id})
			return
		}
		c.JSON(http.StatusOK, v)
	})

	router.GET("/alerts", func(c *gin.Context) {
		all, err := alertDB.getAll()
		if err != nil {
			c.Error(err)
			return
		}
		c.JSON(http.StatusOK, all)
	})

	//---------------------------------
	router.POST("/conditions", func(c *gin.Context) {
		var condition Condition
		err := c.BindJSON(&condition)
		if err != nil {
			c.Error(err)
			log.Println(err)
			return
		}
		err = conditionDB.write(&condition)
		if err != nil {
			c.Error(err)
			return
		}
		c.JSON(http.StatusCreated, gin.H{"id": condition.ID})
	})

	/*router.PUT("/conditions", func(c *gin.Context) {
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
	})*/

	router.GET("/conditions", func(c *gin.Context) {
		all, err := conditionDB.getAll()
		if err != nil {
			c.Error(err)
			return
		}
		c.JSON(http.StatusOK, all)
	})

	router.GET("/conditions/:id", func(c *gin.Context) {
		id := c.Param("id")
		v, err := conditionDB.readByID(id)
		if err != nil {
			c.Error(err)
			return
		}
		if v == nil {
			c.JSON(http.StatusNotFound, gin.H{"id": id})
			return
		}
		c.JSON(http.StatusOK, v)
	})

	router.DELETE("/conditions/:id", func(c *gin.Context) {
		id := c.Param("id")
		ok, err := conditionDB.deleteByID(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"id": id})
			return
		}
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"id": id})
			return
		}
		c.JSON(http.StatusOK, nil)
	})

	//---------------------------------

	err = router.Run(serviceAddress)
	return err
}

func app() int {

	var err error

	// handles the command line flags, finds the discover service, registers us,
	// and figures out our own server address
	svc, err := piazza.NewDiscoverService(os.Args[0], "localhost:12342", "localhost:3000")
	if err != nil {
		log.Print(err)
		return 1
	}

	err = runAlertServer(svc.BindTo, svc.DiscoverAddress, *svc.DebugFlag)
	if err != nil {
		log.Print(err)
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
