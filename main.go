package main
import (
  "fmt"
  "bytes"
  "time"
  "log"
  "net/http"
  "html"
  "encoding/json"
  "io/ioutil"
  "os"
)


import bolt "go.etcd.io/bbolt"

type DataBaseHandler struct {
  Database *bolt.DB
}

type MyData struct {
	Key1 string `json:"key1"`
}


func (h *DataBaseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  fmt.Println("This is my home page")

  // Update the database
	var err = h.Database.Update(func(tx *bolt.Tx) error {
		// Get or create a bucket (similar to a table in relational databases)
		bucket, err := tx.CreateBucketIfNotExists([]byte("MyBucket"))
		if err != nil {
			return err
		}

  body, err := ioutil.ReadAll(r.Body)
  if err != nil {
    http.Error(w, "Error reading request body", http.StatusBadRequest)
    return err
  }

  // Unmarshal the JSON data into a struct
  var requestData MyData
  err = json.Unmarshal(body, &requestData)
  if err != nil {
    http.Error(w, "Error decoding JSON", http.StatusBadRequest)
    return err
  }

	// Access the data
	fmt.Printf("Received data: %+v\n", requestData)

  fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))

	// Write data to the bucket
	key := []byte(requestData.Key1)
	value := []byte(requestData.Key1)
	err = bucket.Put(key, value)
	if err != nil {
		return err
	}

	fmt.Println("Data written to the database:", string(key), string(value))
	return nil
	})

	if err != nil {
		log.Fatal(err)
	}


  w.Write([]byte("This is my home page"))
}

func main() {
  

  mode :=  os.Args[1]

  if mode == "server" {
    startServer()
  } else if mode == "client" {
    fmt.Println("Client")
    addRepository(os.Args[2])
  } else {
    fmt.Println("Unknow mode")
  }
}




func startServer() {
  myHandler := &DataBaseHandler {}
  var err error
  myHandler.Database, err = bolt.Open("/home/td/.GitWatch/mydatabase.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
  defer myHandler.Database.Close()


	// Define a handler function for the "/hello" endpoint
	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		// Call the ServeHTTP method of your custom handler

		myHandler.ServeHTTP(w, r)
	}) 

	// Start the HTTP server on port 8080
  port := 8080
  fmt.Printf("Server listening on :%d...\n", port)
  err = http.ListenAndServe("127.0.0.1:8000", nil)
  if err != nil {
  	fmt.Println("Error:", err)
  }
}


func addRepository(path string) {
// URL to send the POST request to
	url := "http://localhost:8000/add"

	// Data to be sent in the POST request
  postData := []byte(fmt.Sprintf(`{"key1": "%s", "key2": "%s"}`, path, path ))
  fmt.Println(string(postData))
	// Create a new POST request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(postData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Set the Content-Type header
	req.Header.Set("Content-Type", "application/json")

	// Create an HTTP client
	client := &http.Client{}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// Print the response status and body
	fmt.Println("Response Status:", resp.Status)
	fmt.Println("Response Body:")
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	fmt.Println(buf.String())
}
