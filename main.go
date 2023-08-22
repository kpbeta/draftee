package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Stats struct {
	Minutes                  int     `json:"minutes"`
	GoalsScored              int     `json:"goals_scored"`
	Assists                  int     `json:"assists"`
	CleanSheets              int     `json:"clean_sheets"`
	GoalsConceded            int     `json:"goals_conceded"`
	OwnGoals                 int     `json:"own_goals"`
	PenaltiesSaved           int     `json:"penalties_saved"`
	PenaltiesMissed          int     `json:"penalties_missed"`
	YellowCards              int     `json:"yellow_cards"`
	RedCards                 int     `json:"red_cards"`
	Saves                    int     `json:"saves"`
	Bonus                    int     `json:"bonus"`
	Bps                      int     `json:"bps"`
	Influence                float64 `json:"influence"`
	Creativity               float64 `json:"creativity"`
	Threat                   float64 `json:"threat"`
	IctIndex                 float64 `json:"ict_index"`
	Starts                   int     `json:"starts"`
	ExpectedGoals            float64 `json:"expected_goals"`
	ExpectedAssists          float64 `json:"expected_assists"`
	ExpectedGoalInvolvements float64 `json:"expected_goal_involvements"`
	ExpectedGoalsConceded    float64 `json:"expected_goals_conceded"`
	TotalPoints              int     `json:"total_points"`
	InDreamteam              bool    `json:"in_dreamteam"`
}
type Element struct {
	Explain string `json:"-"`
	Stats   Stats  `json:"stats"`
}
type Live struct {
	El map[uint16]Element `json:"elements"`
}
type Game struct {
	CurrentEvent          uint8  `json:"current_event"`
	CurrentEventFinished  bool   `json:"current_event_finished"`
	NextEvent             uint8  `json:"next_event"`
	ProcessingStatus      string `json:"processing_status"`
	TradesTimeForApproval bool   `json:"trades_time_for_approval"`
	WaiversProcessed      bool   `json:"waivers_processed"`
}
type League struct {
	AdminEntry         int       `json:"admin_entry"`
	Closed             bool      `json:"closed"`
	DraftDt            time.Time `json:"draft_dt"`
	DraftPickTimeLimit int       `json:"draft_pick_time_limit"`
	DraftStatus        string    `json:"draft_status"`
	DraftTzShow        string    `json:"draft_tz_show"`
	ID                 int       `json:"id"`
	KoRounds           int       `json:"ko_rounds"`
	MakeCodePublic     bool      `json:"make_code_public"`
	MaxEntries         int       `json:"max_entries"`
	MinEntries         int       `json:"min_entries"`
	Name               string    `json:"name"`
	Scoring            string    `json:"scoring"`
	StartEvent         int       `json:"start_event"`
	StopEvent          int       `json:"stop_event"`
	Trades             string    `json:"trades"`
	TransactionMode    string    `json:"transaction_mode"`
	Variety            string    `json:"variety"`
}
type LeagueEntries []struct {
	EntryID         int       `json:"entry_id"`
	EntryName       string    `json:"entry_name"`
	ID              int       `json:"id"`
	JoinedTime      time.Time `json:"joined_time"`
	PlayerFirstName string    `json:"player_first_name"`
	PlayerLastName  string    `json:"player_last_name"`
	ShortName       string    `json:"short_name"`
	WaiverPick      int       `json:"waiver_pick"`
}
type Matches []struct {
	Event              int         `json:"event"`
	Finished           bool        `json:"finished"`
	LeagueEntry1       int         `json:"league_entry_1"`
	LeagueEntry1Points int         `json:"league_entry_1_points"`
	LeagueEntry2       int         `json:"league_entry_2"`
	LeagueEntry2Points int         `json:"league_entry_2_points"`
	Started            bool        `json:"started"`
	WinningLeagueEntry interface{} `json:"winning_league_entry"`
	WinningMethod      interface{} `json:"winning_method"`
}
type Standings []struct {
	LastRank      interface{} `json:"last_rank"`
	LeagueEntry   int         `json:"league_entry"`
	MatchesDrawn  int         `json:"matches_drawn"`
	MatchesLost   int         `json:"matches_lost"`
	MatchesPlayed int         `json:"matches_played"`
	MatchesWon    int         `json:"matches_won"`
	PointsAgainst int         `json:"points_against"`
	PointsFor     int         `json:"points_for"`
	Rank          interface{} `json:"rank"`
	RankSort      interface{} `json:"rank_sort"`
	Total         int         `json:"total"`
}
type Draft struct {
	League        League        `json:"league"`
	LeagueEntries LeagueEntries `json:"league_entries"`
	Matches       Matches       `json:"matches"`
	Standings     Standings     `json:"standings"`
}
type Club struct {
	Squad        Squad         `json:"picks"`
	EntryHistory struct{}      `json:"-"`
	Subs         []interface{} `json:"-"`
}
type Squad []struct {
	Element       int  `json:"element"`
	Position      int  `json:"position"`
	IsCaptain     bool `json:"is_captain"`
	IsViceCaptain bool `json:"is_vice_captain"`
	Multiplier    int  `json:"multiplier"`
}
type Players []struct {
	ID                               int         `json:"id"`
	Assists                          int         `json:"assists"`
	Bonus                            int         `json:"bonus"`
	Bps                              int         `json:"bps"`
	CleanSheets                      int         `json:"clean_sheets"`
	Creativity                       string      `json:"creativity"`
	GoalsConceded                    int         `json:"goals_conceded"`
	GoalsScored                      int         `json:"goals_scored"`
	IctIndex                         string      `json:"ict_index"`
	Influence                        string      `json:"influence"`
	Minutes                          int         `json:"minutes"`
	OwnGoals                         int         `json:"own_goals"`
	PenaltiesMissed                  int         `json:"penalties_missed"`
	PenaltiesSaved                   int         `json:"penalties_saved"`
	RedCards                         int         `json:"red_cards"`
	Saves                            int         `json:"saves"`
	Threat                           string      `json:"threat"`
	YellowCards                      int         `json:"yellow_cards"`
	Starts                           int         `json:"starts"`
	ExpectedGoals                    string      `json:"expected_goals"`
	ExpectedAssists                  string      `json:"expected_assists"`
	ExpectedGoalInvolvements         string      `json:"expected_goal_involvements"`
	ExpectedGoalsConceded            string      `json:"expected_goals_conceded"`
	Added                            time.Time   `json:"added"`
	ChanceOfPlayingNextRound         int         `json:"chance_of_playing_next_round"`
	ChanceOfPlayingThisRound         int         `json:"chance_of_playing_this_round"`
	Code                             int         `json:"code"`
	DraftRank                        int         `json:"draft_rank"`
	DreamteamCount                   int         `json:"dreamteam_count"`
	EpNext                           interface{} `json:"ep_next"`
	EpThis                           interface{} `json:"ep_this"`
	EventPoints                      int         `json:"event_points"`
	FirstName                        string      `json:"first_name"`
	Form                             string      `json:"form"`
	InDreamteam                      bool        `json:"in_dreamteam"`
	News                             string      `json:"news"`
	NewsAdded                        time.Time   `json:"news_added"`
	NewsReturn                       interface{} `json:"news_return"`
	NewsUpdated                      interface{} `json:"news_updated"`
	PointsPerGame                    string      `json:"points_per_game"`
	SecondName                       string      `json:"second_name"`
	SquadNumber                      interface{} `json:"squad_number"`
	Status                           string      `json:"status"`
	TotalPoints                      int         `json:"total_points"`
	WebName                          string      `json:"web_name"`
	InfluenceRank                    int         `json:"influence_rank"`
	InfluenceRankType                int         `json:"influence_rank_type"`
	CreativityRank                   int         `json:"creativity_rank"`
	CreativityRankType               int         `json:"creativity_rank_type"`
	ThreatRank                       int         `json:"threat_rank"`
	ThreatRankType                   int         `json:"threat_rank_type"`
	IctIndexRank                     int         `json:"ict_index_rank"`
	IctIndexRankType                 int         `json:"ict_index_rank_type"`
	FormRank                         interface{} `json:"form_rank"`
	FormRankType                     interface{} `json:"form_rank_type"`
	PointsPerGameRank                interface{} `json:"points_per_game_rank"`
	PointsPerGameRankType            interface{} `json:"points_per_game_rank_type"`
	CornersAndIndirectFreekicksOrder interface{} `json:"corners_and_indirect_freekicks_order"`
	CornersAndIndirectFreekicksText  string      `json:"corners_and_indirect_freekicks_text"`
	DirectFreekicksOrder             interface{} `json:"direct_freekicks_order"`
	DirectFreekicksText              string      `json:"direct_freekicks_text"`
	PenaltiesOrder                   interface{} `json:"penalties_order"`
	PenaltiesText                    string      `json:"penalties_text"`
	ElementType                      int         `json:"element_type"`
	Team                             int         `json:"team"`
}
type Player struct {
	ID                               int         `json:"id"`
	Assists                          int         `json:"assists"`
	Bonus                            int         `json:"bonus"`
	Bps                              int         `json:"bps"`
	CleanSheets                      int         `json:"clean_sheets"`
	Creativity                       string      `json:"creativity"`
	GoalsConceded                    int         `json:"goals_conceded"`
	GoalsScored                      int         `json:"goals_scored"`
	IctIndex                         string      `json:"ict_index"`
	Influence                        string      `json:"influence"`
	Minutes                          int         `json:"minutes"`
	OwnGoals                         int         `json:"own_goals"`
	PenaltiesMissed                  int         `json:"penalties_missed"`
	PenaltiesSaved                   int         `json:"penalties_saved"`
	RedCards                         int         `json:"red_cards"`
	Saves                            int         `json:"saves"`
	Threat                           string      `json:"threat"`
	YellowCards                      int         `json:"yellow_cards"`
	Starts                           int         `json:"starts"`
	ExpectedGoals                    string      `json:"expected_goals"`
	ExpectedAssists                  string      `json:"expected_assists"`
	ExpectedGoalInvolvements         string      `json:"expected_goal_involvements"`
	ExpectedGoalsConceded            string      `json:"expected_goals_conceded"`
	Added                            time.Time   `json:"added"`
	ChanceOfPlayingNextRound         int         `json:"chance_of_playing_next_round"`
	ChanceOfPlayingThisRound         int         `json:"chance_of_playing_this_round"`
	Code                             int         `json:"code"`
	DraftRank                        int         `json:"draft_rank"`
	DreamteamCount                   int         `json:"dreamteam_count"`
	EpNext                           interface{} `json:"ep_next"`
	EpThis                           interface{} `json:"ep_this"`
	EventPoints                      int         `json:"event_points"`
	FirstName                        string      `json:"first_name"`
	Form                             string      `json:"form"`
	InDreamteam                      bool        `json:"in_dreamteam"`
	News                             string      `json:"news"`
	NewsAdded                        time.Time   `json:"news_added"`
	NewsReturn                       interface{} `json:"news_return"`
	NewsUpdated                      interface{} `json:"news_updated"`
	PointsPerGame                    string      `json:"points_per_game"`
	SecondName                       string      `json:"second_name"`
	SquadNumber                      interface{} `json:"squad_number"`
	Status                           string      `json:"status"`
	TotalPoints                      int         `json:"total_points"`
	WebName                          string      `json:"web_name"`
	InfluenceRank                    int         `json:"influence_rank"`
	InfluenceRankType                int         `json:"influence_rank_type"`
	CreativityRank                   int         `json:"creativity_rank"`
	CreativityRankType               int         `json:"creativity_rank_type"`
	ThreatRank                       int         `json:"threat_rank"`
	ThreatRankType                   int         `json:"threat_rank_type"`
	IctIndexRank                     int         `json:"ict_index_rank"`
	IctIndexRankType                 int         `json:"ict_index_rank_type"`
	FormRank                         interface{} `json:"form_rank"`
	FormRankType                     interface{} `json:"form_rank_type"`
	PointsPerGameRank                interface{} `json:"points_per_game_rank"`
	PointsPerGameRankType            interface{} `json:"points_per_game_rank_type"`
	CornersAndIndirectFreekicksOrder interface{} `json:"corners_and_indirect_freekicks_order"`
	CornersAndIndirectFreekicksText  string      `json:"corners_and_indirect_freekicks_text"`
	DirectFreekicksOrder             interface{} `json:"direct_freekicks_order"`
	DirectFreekicksText              string      `json:"direct_freekicks_text"`
	PenaltiesOrder                   interface{} `json:"penalties_order"`
	PenaltiesText                    string      `json:"penalties_text"`
	ElementType                      int         `json:"element_type"`
	Team                             int         `json:"team"`
}
type Bootstrap struct {
	Players Players `json:"elements"`
}

func readDraft() Draft {
	file, err := os.Open("data-draft-league.json")
	if err != nil {
		fmt.Println(err)
	}

	// defer the closing of the file
	defer file.Close()

	// Create a new decoder
	var draft Draft
	err = json.NewDecoder(file).Decode(&draft)
	if err != nil {
		fmt.Println("Error: Contact Admin")
	}

	return draft
}
func getCurrentEvent() uint8 {
	currEvent := "https://draft.premierleague.com/api/game"
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("GET", currEvent, nil)
	if err != nil {
		fmt.Println("Error: Contact admin")
	}

	req.Header.Set("Authority", "draft.premierleague.com")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error: ask administrator")
	}

	var event Game
	err = json.NewDecoder(resp.Body).Decode(&event)
	if err != nil {
		fmt.Println("Error: ask administrator")
	}

	return event.CurrentEvent
}
func getLiveRequest(gw uint8) Live {
	// TODO: This is insecure; use only in dev environments.
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	defer client.CloseIdleConnections()

	req, err := http.NewRequest("GET",
		"https://draft.premierleague.com/api/event/"+
			strconv.Itoa(int(gw))+"/live",
		nil)

	if err != nil {
		// handle err
	}
	req.Header.Set("Authority", "draft.premierleague.com")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	resp, err := client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		fmt.Println("Error: ask administrator")
	}

	var vals Live
	err = json.NewDecoder(resp.Body).Decode(&vals)

	if err != nil {
		fmt.Println("Error:", err)
		return Live{}
	}

	return vals
}
func getDraftClubs(player uint32, gw uint8) Club {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}
	uri := "https://draft.premierleague.com/api/entry/" +
		strconv.Itoa(int(player)) + "/event/" + strconv.Itoa(int(gw))
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		fmt.Println("Err: contact Admin")
	}
	req.Header.Set("Authority", "draft.premierleague.com")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error: ask administrator")
	}
	defer resp.Body.Close()

	var club Club
	err = json.NewDecoder(resp.Body).Decode(&club)
	if err != nil {
		fmt.Println("Error: contact Admin")
		return Club{}
	}

	return club
}

//	func unnamed(live Live, clubs map[string]Club) {
//		fmt.Println(live)
//		fmt.Println(clubs)
//
//		for player, club := range clubs {
//			fmt.Println(player, club)
//			points := 0
//			for _, player := range club.Squad {
//				playerId := uint16(player.Element)
//				stats := live.El[playerId].Stats
//				fmt.Println(stats)
//				if stats.Minutes > 0 {
//
//					points += stats.TotalPoints
//				}
//			}
//		}
//	}
func getPlayers() {
	req, err := http.NewRequest("GET", "https://draft.premierleague.com/api/bootstrap-static", nil)
	if err != nil {
		// handle err
	}
	req.Header.Set("Authority", "draft.premierleague.com")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Mobile Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error: contact Admin")
	}
	defer resp.Body.Close()

}
func readPlayers() Players {
	file, err := os.Open("data-bootstrap-static.json")
	if err != nil {
		fmt.Println(err)
		return Players{}
	}

	// defer the closing of the file
	defer file.Close()

	// Create a new decoder
	var bootstrap Bootstrap
	err = json.NewDecoder(file).Decode(&bootstrap)
	if err != nil {
		fmt.Println("Error: Contact Admin")
	}

	//fmt.Println(bootstrap.Players)
	return bootstrap.Players
}

func getOutput() string {
	event := getCurrentEvent()
	draft := readDraft()

	clubs := map[string]Club{}
	for _, user := range draft.LeagueEntries {
		clubs[user.PlayerFirstName] = getDraftClubs(uint32(user.EntryID), event)
	}

	TEAMS := []string{"ARS", "AVL", "BOU", "BRE", "BHA", "BUR", "CHE", "CRY", "EVE", "FUL",
		"LIV", "LUT", "MCI", "MUN", "NEW", "NFO", "SHU", "TOT", "WHU", "WOL"}
	POS := []string{"GK", "DF", "MD", "FD"}

	players := map[uint16]Player{}
	for _, pl := range readPlayers() {
		players[uint16(pl.ID)] = pl
	} // Fastness by serializing deserealizing this?

	var out string
	live := getLiveRequest(event).El
	for owner, club := range clubs {

		total := 0
		table := `<table class="col-lg-5 col-md-5 col-sm-12 striped bordered">` +
			"<tr>" +
			"<th>Player</th><th>Club</th><th>Pos</th>" +
			"<th>GS</th><th>A</th><th>GA</th><th>YC</th><th>BON</th>" +
			"<th>PTS</th></tr>"
		for _, pl := range club.Squad {
			player := players[uint16(pl.Element)]
			playerLiveStat := live[uint16(pl.Element)].Stats

			table += fmt.Sprintf(
				"<tr> <td>%s</td> <td>%s</td> <td>%s</td>"+
					" <td>%d</td> <td>%d</td> <td>%d</td> <td>%d</td><td>%d</td><td>%d</td></tr>",
				player.WebName, TEAMS[player.Team-1], POS[player.ElementType-1],
				playerLiveStat.GoalsScored,
				playerLiveStat.Assists,
				playerLiveStat.GoalsConceded,
				playerLiveStat.YellowCards,
				playerLiveStat.Bonus,
				playerLiveStat.TotalPoints)
			total += playerLiveStat.TotalPoints
		}
		table += "</table>"
		out += "<b>" + owner + " [Total Points: " + strconv.Itoa(total) + "]</b>"
		out += table + "</br>"
	}

	html := `<!DOCTYPE html><html><head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<link href="https://cdn.jsdelivr.net/npm/bootstrap@5.0.2/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-EVSTQN3/azprG1Anm3QDgpJLIm9Nao0Yz1ztcQTwFspd3yD65VohhpuuCOmLASjC" crossorigin="anonymous">
<style>
body {
font-size: 8pt;
}
.bordered {
border: 2px solid #dee2e6; /* Outline style */
}

.striped tbody tr:nth-of-type(odd) {
background-color: #f5f5f5; /* Alternate row coloring */
}

</style>
`
	html += "<title>Live Draft Stats</title></head><body>"
	html += "<center><h1>GAMEWEEK " + strconv.Itoa(int(event)) + "<h1></center>"

	html += `<div class="col">` + out + "</div>"
	html += "</body></html>"

	fmt.Println(html)
	return html
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, getOutput())
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}