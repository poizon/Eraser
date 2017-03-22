package main

import (
  "bufio"
  "bytes"
  "fmt"
  "os"
  "strings"
)

func main() {
  reader := bufio.NewReader(os.Stdin)
  var mydir string

  if len(os.Args) > 1 {
    mydir = os.Args[1]
    fmt.Println("Get dir from args")
  } else {
    fmt.Print("Enter full path to directory to delete recursive: ")
    mydir_from_term, _ := reader.ReadString('\n')
    mydir = strings.TrimSpace(mydir_from_term)
  }

  dir, err := os.Open(mydir)

  if err != nil {
    fmt.Println(err)
  }

  defer dir.Close()

  fileInfos, err := dir.Readdir(-1)

  for _, fi := range fileInfos {
    // fmt.Printf("\\%s\\%s\n", mydir, fi.Name())
    fmt.Printf("%s\n", mydir+"/"+fi.Name())
    if fi.IsDir() {
      fmt.Printf("======== \\%s\\%s ======\n", mydir, fi.Name())
      var buffer bytes.Buffer
      buffer.WriteString(mydir)
      buffer.WriteString("/")
      buffer.WriteString(fi.Name())
      readDirRec(buffer.String())
    } else {
      err := os.Remove(mydir + "/" + fi.Name())
      if err != nil {
        fmt.Println(err)
      }
    }
  }
}

func readDirRec(mydir string) {

  dir, err := os.Open(mydir)

  if err != nil {
    fmt.Println(err)
  }

  defer dir.Close()

  fileInfos, err := dir.Readdir(-1)

  if err != nil {
    fmt.Println(err)
  }

  for _, fi := range fileInfos {
    // fmt.Printf("\\%s\\%s\n", mydir, fi.Name())
    fmt.Printf("%s\n", mydir+"/"+fi.Name())
    if fi.IsDir() {
      fmt.Printf("======== \\%s\\%s ======\n", mydir, fi.Name())
      var buffer bytes.Buffer
      buffer.WriteString(mydir)
      buffer.WriteString("/")
      buffer.WriteString(fi.Name())
      readDirRec(buffer.String())
    } else {
      err := os.Remove(mydir + "/" + fi.Name())
      if err != nil {
        fmt.Println(err)
      }
    }
  }

}
