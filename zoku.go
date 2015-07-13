package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "math/rand"
    // "os"
    "net/http"
    "time"
)

const (
    // txt dir path
    ZUKI = "/home/ahmad/Documents/gira/zuki"
    // config file
    AOMI = "/home/ahmad/Documents/gira/mochi/aomi.json"
)

var (
    Now time.Time
    Start *Starts
    Rng *rand.Rand
)

type Starts struct {
    Count int
    Time int64
    Xav []byte
}

func NewStarts() *Starts {
    s0 := Starts{}
    s0.Count = 0
    // s0.Time = time.Now().Unix()
    s0.Xav = make([]byte, 1)
    s0.Xav[0] = s0.Rand()
    return &s0
}

func (s *Starts) Now() {
    s.Count = s.Count + 1
    s.Time = time.Now().Unix()
    s.Xav = append(s.Xav, s.Rand())
}

func (s *Starts) Rand() byte {
    return byte(Rng.Intn(255))
}

func Load() {
    // keep data as byte stream until ready to write
    b0, err := ioutil.ReadFile(AOMI)
    if err != nil {
        fmt.Println(err)
        // create data struct
        Start = NewStarts()
        // Start.Count = 1
        // Start.Time = time.Now().Unix()
        // else load from file
        return
    }
    Start = NewStarts()
    err = json.Unmarshal(b0, Start)
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println(b0, Start)
}

func Save() {
    b0, err := json.Marshal(Start)
    if err != nil {
        fmt.Println(err)
    }
    err = ioutil.WriteFile(AOMI, b0, 0644)
    if err != nil {
        fmt.Println(err)
    }
    // starts should save an array of start times
}

func Dir() {
    fi, err := ioutil.ReadDir(ZUKI)
    if err != nil {
        fmt.Println(err)
    }
    for i0, fi0 := range fi {
        fmt.Println(i0)
        fmt.Printf("File name: %s.\n", fi0.Name())
        fmt.Printf("Size in bytes: %d.\n", fi0.Size())
        fmt.Printf("Mode permission: %vi.\n", fi0.Mode())
        fmt.Printf("Creation date: %s.\n", fi0.ModTime().Format(time.RFC1123Z))
        fmt.Println(fi0.Sys())
    }
    path := ""
    path += ZUKI
    path += "/"
    path += fi[33].Name()
    b0, err0 := ioutil.ReadFile(path)
    if err0 != nil {
        fmt.Println(err0)
    }
    // fmt.Println(path, b0, string(b0))
    fmt.Println(len(b0))
    // new line byte = 10

    // read line by line
    // file := os.Open(path)
    // defer Close
    // scanner := bufio.NewScanner(file)
    // scanner.Scan()
    // scanner.Text()
}

func Web() {
    fmt.Println("listening...")
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "boki.html")
        fmt.Println(r)
    })
    // register other handlers
    err := http.ListenAndServe(":8008", nil)
    if err != nil {
        fmt.Println(err)
    }
}

func main() {
    Now = time.Now()
    Rng = rand.New(rand.NewSource(Now.UnixNano()))
    fmt.Println("okaq mochi web app ok")
    fmt.Printf("Started at: %s.\n", Now.Format(time.RFC1123Z))
    fmt.Printf("Loading config file %s.\n", AOMI)
    Load()
    Start.Now()
    fmt.Println(Start)
    Save()
    fmt.Printf("Opening dir: %s.\n", ZUKI)
    Dir()
    fmt.Println("starting web server on port 8008")
    Web()
    // thread wait after listen is called
}
