package main

import (
    "bufio"
    "encoding/json"
    "fmt"
    "io"
    "io/ioutil"
    "math/rand"
    "os"
    "net/http"
    "runtime"
    "strings"
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
    Fi []os.FileInfo
    // Ideas []*Idea
    Ideas map[string]*Idea
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

type Idea struct {
    Path string
    Date string
    Title string
    Creator string
    Txt []byte
}

func NewIdea() *Idea {
    return &Idea{}
}

func (id0 *Idea) Pop() {
    // after obtaining path
    // read in file data
    f0, err := os.Open(id0.Path)
    if err != nil {
        fmt.Println(err)
    }
    defer f0.Close()
    b0 := bufio.NewReader(f0)
    t0, err := b0.ReadString('\n')
    if err != nil {
        fmt.Println(err)
    }
    t0 = strings.TrimSpace(t0)
    fmt.Printf("Title is: %s.\n", t0)
    c0, err := b0.ReadString('\n')
    if err != nil {
        fmt.Println(err)
    }
    c0 = strings.TrimSpace(c0)
    fmt.Printf("Creator is: %s.\n", c0)
    d0, err := b0.ReadString('\n')
    if err != nil {
        fmt.Println(err)
    }
    d0 = strings.TrimSpace(d0)
    fmt.Printf("Date is: %s.\n", d0)
    id0.Title = t0
    id0.Date = d0
    id0.Creator = c0
    s0, err := b0.ReadBytes(255)
    if err == io.EOF {
        fmt.Println(err)
    }
    fmt.Printf("Bytes: %d.\n", len(s0))
    // fmt.Printf("Text: %s.\n", string(s0))
    id0.Txt = s0
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
    var err error
    Fi, err = ioutil.ReadDir(ZUKI)
    if err != nil {
        fmt.Println(err)
    }
    for i0, fi0 := range Fi {
        fmt.Println(i0)
        fmt.Printf("File name: %s.\n", fi0.Name())
        fmt.Printf("Size in bytes: %d.\n", fi0.Size())
        fmt.Printf("Mode permission: %vi.\n", fi0.Mode())
        fmt.Printf("Creation date: %s.\n", fi0.ModTime().Format(time.RFC1123Z))
        fmt.Println(fi0.Sys())
    }
    /*
    path := ""
    path += ZUKI
    path += "/"
    path += fi[33].Name()
    b0, err0 := ioutil.ReadFile(path)
    if err0 != nil {
        fmt.Println(err0)
    }
    */
    // fmt.Println(path, b0, string(b0))
    // fmt.Println(len(b0))
    // new line byte = 10

    // read line by line
    // file := os.Open(path)
    // defer Close
    // scanner := bufio.NewScanner(file)
    // scanner.Scan()
    // scanner.Text()
}

func Pop() {
    // populate ideas file data cache
    // Ideas = make([]*Idea, len(Fi))
    Ideas = make(map[string]*Idea)
    for _, f0 := range Fi {
        path := ""
        path += ZUKI
        path += "/"
        path += f0.Name()
        id0 := NewIdea()
        id0.Path = path
        id0.Pop()
        Ideas[f0.Name()] = id0
        // map name to idea for easy access
    }
}

func StatHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println(r)
    // w.Write([]byte("stats live"))
    ms := &runtime.MemStats{}
    runtime.ReadMemStats(ms)
    b0, err := json.Marshal(ms)
    if err != nil {
        fmt.Println(err)
    }
    w.Write(b0)
}

func IdeasHandler(w http.ResponseWriter, r *http.Request) {
    b0, err := json.Marshal(Ideas)
    if err != nil {
        fmt.Println(err)
    }
    w.Write(b0)
}

func Web() {
    fmt.Println("listening...")
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "boki.html")
        fmt.Println(r)
    })
    // register other handlers
    http.HandleFunc("/s", StatHandler)
    http.HandleFunc("/i", IdeasHandler)
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
    fmt.Println("Populating file data cache")
    Pop()
    fmt.Println("starting web server on port 8008")
    Web()
    // thread wait after listen is called
}
