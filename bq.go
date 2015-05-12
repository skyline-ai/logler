package logler

import (
	"log"
	"os"
	"time"

	bqstream "github.com/rounds/go-bqstreamer"
	confy "github.com/spf13/viper"
	"google.golang.org/api/bigquery/v2"
)

var (
	ms *bqstream.MultiStreamer
)

func init() {
	//For GOPATH
	confy.BindEnv("gopath", "GOPATH")

	//Set Default (A good practice)
	//confy.SetDefault("gopath", "/go")
	confy.SetDefault("jwtconfig", "/Users/eranchetz/Certs/streamrail-41d1d158c64e.json")
	confy.SetDefault("numStreamers", 5)
	confy.SetDefault("maxRows", 50) //500
	confy.SetDefault("maxRetryInsert", 10)

	// This will read all environemnt vars starting with PUBZ
	// Vars will be evalueted when running Get()
	confy.SetEnvPrefix("WT") // will be uppercased automatically
	confy.AutomaticEnv()

	//Example:
	//os.Setenv("WT_JWTCONFIG", "/Certs/sr-key.json") // typically done outside of the app
	// id := confy.Get("JWTConfig") // uppercased
}

func initBQ() {
	templog := log.New(os.Stdout,
		"BQ: ",
		log.Ldate|log.Ltime)
	jwtConfig, err := bqstream.NewJWTConfig(confy.GetString("jwtconfig"))
	if err != nil {
		templog.Println("Failed to connect to BQ with default key ", confy.GetString("jwtconfig"),
			"\nPlease set env var WT_JWTCONFIG=<path_to_secret.json>")
		templog.Fatalln(err.Error())
	}

	// Set MultiStreamer configuration.
	numStreamers := confy.GetInt("numStreamers")     // Number of concurrent sub-streamers (workers) to use.
	maxRows := confy.GetInt("maxRows")               // Amount of rows queued before forcing insert to BigQuery.
	maxDelay := 1 * time.Second                      // Time to pass between forcing insert to BigQuery.
	sleepBeforeRetry := 1 * time.Second              // Time to wait between failed insert retries.
	maxRetryInsert := confy.GetInt("maxRetryInsert") // Maximum amount of failed insert retries before discarding rows and moving on.

	// Init a new multi-streamer.
	ms, err = bqstream.NewMultiStreamer(
		jwtConfig, numStreamers, maxRows, maxDelay, sleepBeforeRetry, maxRetryInsert)
	if err != nil {
		log.Fatalln(err.Error())
	}
	// Start multi-streamer and workers.
	ms.Start()
	//we should call this somewhere ...
	//defer ms.Stop()

	templog.Println("BQ MultiStreamer connected with", confy.GetString("jwtconfig"))

	// Worker errors are reported to MultiStreamer.Errors channel.
	// This inits a goroutine the reads from this channel and logs errors.
	//
	// It can be closed by sending "true" to the shutdownErrorChan channel.
	shutdownErrorChan := make(chan bool)
	go func() {
		var err error
		readErrors := true
		for readErrors {
			select {
			case <-shutdownErrorChan:
				readErrors = false
			case err = <-ms.Errors:
				log.Println(err)
			}
		}
	}()
	//defer func() { shutdownErrorChan <- true }()
}

func sendBQ(bqmap map[string]bigquery.JsonValue) {
	projectid := "streamrail"
	dataset := "bq_test"
	table := "test_table"
	log.Printf("sending :%v to %s:%s.%s ", bqmap, projectid, dataset, table)
	ms.QueueRow(projectid, dataset, table, bqmap)
}
