package main
import (
  "fmt"
  "bytes"
  "time"
  "log"
  "net/http"
  "io/ioutil"
  "os"
  "GitWatch/src"
)

import bolt "go.etcd.io/bbolt"


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
      return "Unknown argument, should be add|status|remove"
    }
  } 
  return "Unknown mode"
}


func main() {
  fmt.Println(realMain(os.Args))
}


func startServer() {
  myHandler := &utils.DataBaseHandler {}
  var err error
  myHandler.Database, err = bolt.Open(utils.DB_PATH, 0600, &bolt.Options{Timeout: 1 * time.Second})
  defer myHandler.Database.Close()


	// Define a handler function for the "/hello" endpoint
	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		myHandler.Add(w, r)
	}) 

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		myHandler.StatusHandler(w, r)
	}) 


	http.HandleFunc("/remove", func(w http.ResponseWriter, r *http.Request) {
		myHandler.RemoveHandler(w, r)
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
	req, err := http.NewRequest("POST", utils.URL + "/add", bytes.NewBuffer(postData))
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
    // TODO Habndle server not runniong here
		return "Error during sending"
	}
	defer resp.Body.Close()

	// Print the response status and body
	log.Println("Response Status:", resp.Status)
  if  resp.Status != "200 OK" {
	  log.Println("There was an error during adding the current path")
	  return "Error during adding the path"
	}
  return path + " added"
}


func getStatus() {
  // Make the GET request
	response, err := http.Get(utils.URL + "/status")
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
	req, err := http.NewRequest("POST", utils.URL + "/remove", bytes.NewBuffer(postData))
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


  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    fmt.Println("Error reading response body:", err)
    return ""
  }

	// Print the response status and body
	log.Println("Response Status:", resp.Status)
  if  resp.Status == "404 Not Found" {
	  return string(body)
	}

  if  resp.Status != "200 OK" {
	  log.Println("There was an error during removing the current path")
	  return "Error during removing the path"
	}

  print()
  return  string(body)
}
