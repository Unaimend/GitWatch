package main
import (
  "fmt"
  "bytes"
  "time"
  "log"
  "strings"
  "net/http"
  "html"
  "encoding/json"
  "io/ioutil"
  "os"
  "GitWatch/src"
	. "github.com/go-git/go-git/v5/_examples"
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
	fmt.Printf("Adder Received data: %+v\n", keyValueMap)

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


func (h *DataBaseHandler) statusHandler(w http.ResponseWriter, r *http.Request) {

  var respBody = ""
  err := h.Database.View(func(tx *bolt.Tx) error {
		// Assume you have a bucket named "mybucket"
    b := tx.Bucket([]byte(BUCKET_NAME))
    
    c := b.Cursor()
    

    for k, v := c.First(); k != nil; k, v = c.Next() {
      respBody +=  "Current repository: " + string(v) + "\n"
      var repo, _ = utils.OpenRepo(string(v))
	    wk, err := repo.Worktree()
      CheckIfError(err)
        
      status, _:= wk.Status()


    	if status.IsClean() == true {
        w.Write([]byte(""))
        return err
      }

      var status_string  []string = strings.Split(status.String(), "\n")
      
      for k, v := range status_string {
        status_string[k] = "\t" + v + "\n"
        respBody += status_string[k]
      }
    }
    return nil
	 
  })

	if err != nil {
		log.Fatal(err)
	}
  w.Write([]byte(respBody))
}


func (h *DataBaseHandler) removeHandler(w http.ResponseWriter, r *http.Request) {
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
	  fmt.Printf("Remover Received data: %+v\n", keyValueMap)

    fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))

	  // Write data to the bucket
	  key := []byte(keyValueMap["key1"])

		// Delete the key from the bucket.
		err = bucket.Delete(key)
		if err != nil {
			return err
		}
	  return nil
	})

	if err != nil {
		log.Fatal(err)
	}


  w.Write([]byte("this is my home page"))
}


func realMain(args []string) string  {
  mode :=  args[1]
  if mode == "server" {
    startServer()
  } else if mode == "client" {
    arg :=  os.Args[2]
    if arg == "add"  {
      
      return addRepository(os.Args[3])
    } else if arg == "status" {
      getStatus()
      return ""
    } else if arg == "remove" {
      return remove(os.Args[3])
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
		myHandler.add(w, r)
	}) 

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		myHandler.statusHandler(w, r)
	}) 


	http.HandleFunc("/remove", func(w http.ResponseWriter, r *http.Request) {
		myHandler.removeHandler(w, r)
	}) 

	// Start the HTTP server on port 8080
  port := 8000
  log.Printf("Server listening on :%d...\n", port)
  err = http.ListenAndServe("127.0.0.1:8000", nil)
  if err != nil {
  	log.Println("Error:", err)
  }
}



func addRepository(path string) string {
  
  if !utils.IsGitRepository(path) {
    return "Specified path is not a git repository"
  }
	// This key is hardcoded because we use it to access the dictionary that is later unmarshalled 
  // For the actual database we use path as key and as value
  postData := []byte(fmt.Sprintf(`{"key1": "%s"}`, path))
  log.Println("Sending" + string(postData))

  // Create new request with the current path
	req, err := http.NewRequest("POST", URL + "/add", bytes.NewBuffer(postData))
	if err != nil {
		log.Fatal("Error creating request:", err)
		return "Error during adding"
	}

	// Set the Content-Type header
	req.Header.Set("Content-Type", "application/json")

	// Create an HTTP client
	client := &http.Client{}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending request:", err)
		return "Error during sending"
	}
	defer resp.Body.Close()

	// Print the response status and body
	log.Println("Response Status:", resp.Status)
  if  resp.Status != "200 OK" {
	  log.Println("There was an error during adding the current path")
	  return "Error during adding the path"
	}
  return "Path added"
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


func remove(path string) string {
  if !utils.IsGitRepository(path) {
    return "Specified path is not a git repository"
  }
	// This key is hardcoded because we use it to access the dictionary that is later unmarshalled 
  // For the actual database we use path as key and as value
  postData := []byte(fmt.Sprintf(`{"key1": "%s"}`, path))
  log.Println("Sending" + string(postData))

  // Create new request with the current path
	req, err := http.NewRequest("POST", URL + "/remove", bytes.NewBuffer(postData))
	if err != nil {
		log.Fatal("Error creating request:", err)
		return "Error during creating"
	}

	// Set the Content-Type header
	req.Header.Set("Content-Type", "application/json")

	// Create an HTTP client
	client := &http.Client{}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending request:", err)
		return "Error during sending"
	}
	defer resp.Body.Close()

	// Print the response status and body
	log.Println("Response Status:", resp.Status)
  if  resp.Status != "200 OK" {
	  log.Println("There was an error during adding the current path")
	  return "Error during removing the path"
	}
  return path + " removed"
}
