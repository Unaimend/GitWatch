package utils

import (
	"encoding/json"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	. "github.com/go-git/go-git/v5/_examples"
)

import bolt "go.etcd.io/bbolt"
  
type DataBaseHandler struct {
  Database *bolt.DB
}

func (h *DataBaseHandler) Add(w http.ResponseWriter, r *http.Request) {
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
    //.TODO ADD HTTP ERROR
		log.Println("Error unmarshaling JSON:", err)
		return err
	}

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


func (h *DataBaseHandler) StatusHandler(w http.ResponseWriter, r *http.Request) {

  var respBody = ""
  err := h.Database.View(func(tx *bolt.Tx) error {
		// Assume you have a bucket named "mybucket"
    b := tx.Bucket([]byte(BUCKET_NAME))
    
    c := b.Cursor()
    

    for k, v := c.First(); k != nil; k, v = c.Next() {
      respBody +=  "Current repository: " + string(v) + "\n"
      var repo, _ = OpenRepo(string(v))
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


func (h *DataBaseHandler) RemoveHandler(w http.ResponseWriter, r *http.Request) {
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

    value := bucket.Get(key)
    fmt.Println("WHY", value)
		if len(value) == 0  {
      fmt.Println("EMPTRY")
      w.Write([]byte("Repository " + string(key) + " does not exist"))
			return nil
		} 

		err = bucket.Delete(key)
		if err != nil {
			return err
		}
	  return nil
	})

	if err != nil {
		log.Fatal(err)
	}


}
