package main
 
import (
    "fmt"
    "net/http"
//    "strings"
    "io/ioutil"
    //"io"
    //"log"
    //"regexp"
    "crypto/tls"
    "os"
    "encoding/json"
)
 
type scoreboardTeam struct {
	Rank int `json:"ranks"`
	Team_id string `json:"team_id"`
	Score teamScore `json:"score"`
	Problems []teamProblem `json:"problems"`
}


type teamScore struct {
	Num_solved  int `json:"num_solved"`
	Total_time int `json:"total_time"`
}

type teamProblem struct {
	Problem_id string `json:"problem_id"`
	Num_judged int `json:"num_judged"`
	Num_pending int `json:"num_pending"`
	Solved bool `json:"solved"`
	Time uint32 `json:"time"`
}

var scoreboardPage []scoreboardTeam

func readJSONfile(jsonFile string) ([]byte, error) {
	fileContent, err := os.Open(jsonFile)
	if err != nil {
      		return nil, err
   	}
	defer fileContent.Close()
	byteResult, _ := ioutil.ReadAll(fileContent)
	return byteResult, nil
}


func getHtmlPage(webPage string,username string,password string) (string, error) {

    tr := &http.Transport {
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}
    client := &http.Client{Transport: tr}
    req, err := http.NewRequest("GET", webPage, nil)

    if err != nil {
        return "", err
    }
    req.Header.Add("Accept", "application/json")
    req.SetBasicAuth(username, password)
    resp, err := client.Do(req)
    if err != nil {
    		fmt.Printf("error %s", err)
    		return "", err
    }

    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }
    return string(body), nil
}


//Function, takes the scoreboard in JSON format, and return
//Num of pending problems, Num of solved problems

func parseScoreboard(jsonText []byte) (int,int, map[string]int, map[string]int) {
	var (
		NumPending = 0
		NumSolved = 0
		ProbTotal map[string]int
		ProbSolved map[string]int
	)
	json.Unmarshal(jsonText, &scoreboardPage)
	fmt.Println("You have ", len(scoreboardPage))
	ProbTotal = make(map[string]int)
	ProbSolved = make(map[string]int)
	
	for i :=0 ; i < len(scoreboardPage) ; i++ {
		//fmt.Print("Entry : ", i)
		//fmt.Print(" Rank : " , scoreboardPage[i].Rank, " team_id: " , scoreboardPage[i].Team_id)
		//fmt.Println(" Team score-Num-Solved : " , scoreboardPage[i].Score.Num_solved , " Team score-Total-Time : ", scoreboardPage[i].Score.Total_time)
		NumSolved = NumSolved + scoreboardPage[i].Score.Num_solved
		for j :=0 ; j < len(scoreboardPage[i].Problems) ; j++ {
			//fmt.Print("\tProblem: " , scoreboardPage[i].Problems[j].Problem_id)
				ProbTotal[scoreboardPage[i].Problems[j].Problem_id] = ProbTotal[scoreboardPage[i].Problems[j].Problem_id] + scoreboardPage[i].Problems[j].Num_judged
				NumPending = NumPending + scoreboardPage[i].Problems[j].Num_pending
			if scoreboardPage[i].Problems[j].Solved {
				ProbSolved[scoreboardPage[i].Problems[j].Problem_id] = ProbSolved[scoreboardPage[i].Problems[j].Problem_id] + 1
			}
			//ifmt.Println("")
		}
		//fmt.Println("")
	}
	return NumPending, NumSolved , ProbTotal, ProbSolved
}

