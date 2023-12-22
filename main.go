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

var DB_PATH string = "/home/td/.GitWatch/mydatabase.db"
var BUCKET_NAME string = "MyBucket"
var	URL string = "http://localhost:8000"
type DataBaseHandler struct {
  Database *bolt.DB
}


func (h *DataBaseHandler) add(w http.ResponseWriter, r *http.Request) {
  fmt.Println("This is my home page")

  // Update the database
	var err = h.Database.Update(func(tx *bolt.Tx) error {
		// Get or create a bucket (similar to a table in relational databases)
		bucket, err := tx.CreateBucketIfNotExists([]byte(BUCKET_NAME))
		if err != nil {
			return err
		}

  body, err := ioutil.ReadAll(r.Body)
  if err != nil {
    http.Error(w, "Error reading request body", http.StatusBadRequest)
    return err
  }

  var keyValueMap map[string]string

	// Unmarshal the JSON data into the map
	err = json.Unmarshal([]byte(body), &keyValueMap)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return err
	}

	// Access the key-value pair
	// Access the data
	fmt.Printf("Received data: %+v\n", keyValueMap)

  fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))

	// Write data to the bucket
	key := []byte(keyValueMap["key1"])
	value := []byte(keyValueMap["key1"])
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


func (h *DataBaseHandler) status(w http.ResponseWriter, r *http.Request) {

  var respBody = "["
  err := h.Database.View(func(tx *bolt.Tx) error {
		// Assume you have a bucket named "mybucket"
    b := tx.Bucket([]byte(BUCKET_NAME))
    
    fmt.Printf("status")
    c := b.Cursor()
    


    for k, v := c.First(); k != nil; k, v = c.Next() {
      respBody = respBody + fmt.Sprintf(`{"%s, %s"},`, string(k), string(v))
    }


    str := respBody

	  // Convert the string to a rune slice
	  runes := []rune(str)

	  // Check if the string is not empty
	  if len(runes) > 0 {
	  	// Change the last character
	  	runes[len(runes)-1] = ']'
	  }

	  // Convert the rune slice back to a string
	  newStr := string(runes)
    respBody = newStr
    return nil
	 
  })
  fmt.Printf("RESP" + respBody)

	if err != nil {
		log.Fatal(err)
	}
  w.Write([]byte(respBody))
}

func realMain(args []string) string  {
  mode :=  args[1]
  if mode == "server" {
    startServer()
  } else if mode == "client" {
    arg :=  os.Args[2]
    if arg == "add"  {
      addRepository(os.Args[3])
      return "Added repository"
    } else if arg == "status" {
      getStatus()
      return ""
    } else {
      return "Unknown argument, should be add|status"
    }
  } 
  return "Unknown mode"
}


func main() {
  fmt.Println(realMain(os.Args))
}




func startServer() {
  myHandler := &DataBaseHandler {}
  var err error
  myHandler.Database, err = bolt.Open(DB_PATH, 0600, &bolt.Options{Timeout: 1 * time.Second})
  defer myHandler.Database.Close()


	// Define a handler function for the "/hello" endpoint
	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		// Call the ServeHTTP method of your custom handler
    fmt.Println("ADD")
		myHandler.add(w, r)
	}) 

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
    fmt.Println("STATUS")
		myHandler.status(w, r)
	}) 

	// Start the HTTP server on port 8080
  port := 8000
  fmt.Printf("Server listening on :%d...\n", port)
  err = http.ListenAndServe("127.0.0.1:8000", nil)
  if err != nil {
  	fmt.Println("Error:", err)
  }
}



func addRepository(path string) {
// URL to send the POST request to

	// This key is hardcoded because we use it to access the dictionary that is later unmarshalled 
  // For the actual database we use path as key and as value
  postData := []byte(fmt.Sprintf(`{"key1": "%s"}`, path))
  log.Println("Sending" + string(postData))
  // Create new request with the current path
	req, err := http.NewRequest("POST", URL + "/add", bytes.NewBuffer(postData))
	if err != nil {
		log.Fatal("Error creating request:", err)
		return
	}

	// Set the Content-Type header
	req.Header.Set("Content-Type", "application/json")

	// Create an HTTP client
	client := &http.Client{}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// Print the response status and body
	log.Println("Response Status:", resp.Status)
  if  resp.Status != "200 OK" {
	  fmt.Println("There was an error during adding the current path")
	  return
	}
}

func getStatus() {
  // Make the GET request
	response, err := http.Get(URL + "/status")
	if err != nil {
		fmt.Println("Error making GET request:", err)
		return
	}
	defer response.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// Print the response body
	fmt.Println(string(body))
}
