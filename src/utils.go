package utils
//
//
import (
	"fmt"
  "errors"
  //"path"
	"github.com/go-git/go-git/v5"

)

////import bolt "go.etcd.io/bbolt"
//
//
//
func OpenRepo(repoPath string) (*git.Repository, error) {
	  //if(!filepath.isAbs(repoPath)) {
	  //	fmt.Println("Path must be absolute")
    //  os.Exit(1)
    //}
    //TODO: Check if path exists
    //TODO: Check if path contains a .git
	  repo, err := git.PlainOpen(repoPath)
    if err != nil {
	  	fmt.Println()
      return nil,errors.New("Error opening repository" )
	  }
    return repo, nil
}


func IsGitRepository(path string) bool {
	_, err := git.PlainOpen(path)
	return err == nil
}
