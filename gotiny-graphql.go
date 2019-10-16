package main

import (
	"context"
	"cors"
	"engine"
	"flag"
	"fmt"
	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"time"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

type query struct{}

func (_ *query) Hello() string { return "Hello, world!" }

func main() {
	var (
		dbinit = flag.Bool("i", true, "init database flag")
		mongo  = flag.String("m", "mongodb://root:root@192.168.101.68:27017,192.168.101.69:27017,192.168.101.70:27017/?authSource=admin&replicaSet=rs1", "mongod addr flag")
		//mongo = flag.String("m", "", "mongod addr flag")
		db = flag.String("db", "test_db", "mongod addr flag")
	)
	flag.Parse()
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	if _, err := os.Stat(dir + "/data"); err != nil {
		log.Println(err)
		return
	}

	//启动文件日志
	//logFile, logErr := os.OpenFile(dir+"/dal.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	//if logErr != nil {
	//	log.Printf("err: %v\n", logErr)
	//	return
	//
	//}
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	eng := engine.NewDbEngine()
	err = eng.Open(dir, *mongo, *db, *dbinit)
	if err != nil {
		log.Fatal("database connect error")
	}

	b, err := ioutil.ReadFile(dir + "/schema.graphql")
	if err != nil {
		log.Fatal(err.Error())
	}

	opts := []graphql.SchemaOpt{graphql.UseFieldResolvers(), graphql.UseStringDescriptions()}
	schema := graphql.MustParseSchema(string(b), eng, opts...)
	//1083
	mux := http.NewServeMux()
	mux.Handle("/", &relay.Handler{Schema: schema})
	srv := &http.Server{Handler: cors.CORS(mux), ErrorLog: nil}
	srv.Addr = ":8000"
	//cert, err := tls.LoadX509KeyPair(dir+"/server.pem", dir+"/server.key")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//config := &tls.Config{Certificates: []tls.Certificate{cert}}
	//srv.TLSConfig = config
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()
	log.Println("server on tls port", srv.Addr)

	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	cleanup := make(chan bool)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for range signalChan {
			ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)
			go func() {
				_ = srv.Shutdown(ctx)
				cleanup <- true
			}()
			<-cleanup
			eng.Close()
			fmt.Println("safe exit")
			cleanupDone <- true
		}
	}()
	<-cleanupDone
}
