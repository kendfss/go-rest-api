package main

import (
    "fmt"
    "log"
    "net/http"
    "strconv"
    "encoding/json"
    "io/ioutil"
    "sync"
    "os"
    "time"
    "bufio"
    "bytes"
    
    "github.com/gorilla/mux"
)

var counter int
var mutex = &sync.Mutex{}

// Track ...
type Track struct {
    Id      string `json:"Id"`
    Title   string `json:"Title"`
    Artist  string `json:"author"`
    Ep      string `json:"ep"`
    Se      string `json:"se"`
    Link    string `json:"link"`
}
func (self Track) eq(other Track) bool {
    return self.Id==other.Id && self.Title==other.Title && self.Artist==other.Artist 
}




type Season struct {
    Tracks []Track `json:"tracks"`
    Id      string `json:"link"`
}
func (self *Season) add(track Track) (redundant bool) {
    if (*self).has(track) {
        redundant = true
    } else {
        self.Tracks = append(self.Tracks, track)
    }
    return    
}
func (self Season) has(track Track) bool {
    for _, elem := range self.Tracks {
        if elem.Id == track.Id {
            return true
        }
    }
    return false
}
type Arxiv map[string]Season
func (self Arxiv) seasons() (rack []Season) {
    for _, val := range self {
        rack = append(rack, val)
    }
    return
}
func (self *Arxiv) add(season Season) (redundant bool) {
    if (*self).hasId(season.Id) {
        redundant = true
    } else {
        (*self)[season.Id] = season
        // self[season.Id]
        // self[season.Id] = season
    }
    return    
}
func (self Arxiv) has(season Season) bool {
    for _, elem := range self {
        if elem.Id == season.Id {
            return true
        }
    }
    return false
}
func (self Arxiv) hasId(season string) bool {
    for _, elem := range self {
        if elem.Id == season {
            return true
        }
    }
    return false
}
 
// Tracks ...
var Tracks []Track
var Archive Arxiv

// readLines reads a whole file into memory
// and returns a slice of its lines.
func read(path string) (string, error) {
    lines, err := readLines(path)
    if err != nil {
        // log.Printf("Couldn't read %q \n", path)
        return "", err
    }
    str := ""
    for _, line := range lines {
        str += line
    }
    return str, nil
}
func readLines(path string) ([]string, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var lines []string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }
    return lines, scanner.Err()
}

// writeLines writes the lines to the given file.
func writeLines(lines []string, path string) error {
    file, err := os.Create(path)
    if err != nil {
        return err
    }
    defer file.Close()

    w := bufio.NewWriter(file)
    for _, line := range lines {
        fmt.Fprintln(w, line)
    }
    return w.Flush()
}
 
func homePage(w http.ResponseWriter, r *http.Request) {
    // fmt.Fprintf(w, "<b>Welcome to the HomePage!</b>")
    text, err := read("index.html")
    if err != nil {
        log.Printf("Couldn't read %q\n", "index.html")
    }
    fmt.Fprintf(w, text)
    fmt.Println("Endpoint Hit: homePage")
}

func echoString(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "hello")
}

func incrementCounter(w http.ResponseWriter, r *http.Request) {
    mutex.Lock()
    counter++
    fmt.Fprintf(w, strconv.Itoa(counter))
    mutex.Unlock()
}
func returnAllTracks(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Endpoint Hit: returnAllTracks")
    json.NewEncoder(w).Encode(Tracks)
}
func returnWholeSeason(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Endpoint Hit: returnWholeSeason")
    
    vars := mux.Vars(r)
    id := vars["id"]
    // json.NewEncoder(w).Encode(Tracks)
    // results := []Track{}
    // for _, track := range Tracks {
    //     // if muxi
    //     results = append(results, track)
    // }
    json.NewEncoder(w).Encode(Archive[id])
}
func returnAllSeasons(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Endpoint Hit: returnWholeSeason")
    
    json.NewEncoder(w).Encode(Archive.seasons())
}
func returnSingleTrack(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Endpoint Hit: returnSingleTrack")
    vars := mux.Vars(r)
    id := vars["id"]
 
    for _, track := range Tracks {
        if track.Id == id {
            json.NewEncoder(w).Encode(track)
        }
    }
}
func returnEpisode(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Endpoint Hit: returnEpisode")
    vars := mux.Vars(r)
    ep := vars["ep"]
    se := vars["se"]
 
    for _, track := range Tracks {
        if track.Ep == ep && track.Se == se {
            json.NewEncoder(w).Encode(track)
        }
    }
}
func createNewTrack(w http.ResponseWriter, r *http.Request) {
    reqBody, _ := ioutil.ReadAll(r.Body)
    var track Track 
    json.Unmarshal(reqBody, &track)
    Tracks = append(Tracks, track)
 
    json.NewEncoder(w).Encode(track)}
func deleteTrack(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]
 
    for index, track := range Tracks {
        if track.Id == id {
            Tracks = append(Tracks[:index], Tracks[index+1:]...)
        }
    }
}
func computeSeasons(tracks []Track, archive Arxiv) Arxiv {
    // Seasons := []Season{}
    // archive = (archive).(map[string]Season)
    for _, track := range tracks {
        // archive[track.Se].Tracks = append(archive[track.Se].Tracks, track)
        if archive.hasId(track.Se) {
            // (&archive)[track.Se].add(track)
            se := archive[track.Se]
            se.add(track)
            // if !archive[track.Se].has(track) {
            //     archive[track.Se].Tracks = append(archive[track.Se].Tracks, track)
            // }
            // archive[track.Se].add(track)
        } else {
            // archive[track.Se] = Season{
            se := Season{
                Id: track.Se,
                Tracks: []Track{track},
            }
            archive.add(se)
        }
    }
    return archive
} 
func downloadHandler(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    StoredAs := r.Form.Get("StoredAs") // file name
    data, err := ioutil.ReadFile("files/"+StoredAs)
    if err != nil { fmt.Fprint(w, err) }
    http.ServeContent(w, r, StoredAs, time.Now(), bytes.NewReader(data))
}
func handleRequests(port int) {
    portNum := strconv.Itoa(port)
    log.Println("serving on", "http://localhost:"+portNum)
    router := mux.NewRouter().StrictSlash(true)
    router.HandleFunc("/", homePage)
    router.HandleFunc("/tracks", returnAllTracks)
    router.HandleFunc("/seasons", returnAllSeasons)
    router.HandleFunc("/season/{id}", returnWholeSeason)
    router.HandleFunc("/season/{se}/{ep}", returnEpisode)
    router.HandleFunc("/track/{id}", returnSingleTrack)
    router.HandleFunc("/track", createNewTrack).Methods("POST")
    router.HandleFunc("/track/{id}", deleteTrack).Methods("DELETE")
    log.Fatal(http.ListenAndServe(":"+portNum, router))
}

func main() {
    Tracks = []Track{
        Track{
            Id:     "1",
            Title:  "track_one",
            Artist: "eli2and40",
            Link:   "https://www.soundcloud.com/eli2and40/track_one",
            Se: "1",
            Ep: "1",
        },
        Track{
            Id:     "2",
            Title:  "track_two",
            Artist: "eli2and40",
            Link:   "https://www.soundcloud.com/eli2and40/track_two",
            Se: "1",
            Ep: "2",
        },
        Track{
            Id:     "3",
            Title:  "track_three",
            Artist: "eli2and40",
            Link:   "https://www.soundcloud.com/eli2and40/track_three",
            Se: "1",
            Ep: "2",
        },
    }
    // Archive = computeSeasons(Tracks, Archive)
    Archive = computeSeasons(Tracks, Arxiv{})
    handleRequests(8000)
}
