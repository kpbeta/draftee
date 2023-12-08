package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"
)

const site_template string = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/css/bootstrap.min.css">
	<title>Live Draft Stats</title>
	<style>
		body {
		font-size: 9pt;
		}
		.table-condensed>thead>tr>th, .table-condensed>tbody>tr>th, .table-condensed>tfoot>tr>th, .table-condensed>thead>tr>td, .table-condensed>tbody>tr>td, .table-condensed>tfoot>tr>td{
			padding: 1px;
			
		}
	</style>
</head>

<body>
	<center><h1>GAMEWEEK %d <h1></center>
	<div class="container">
		<div class="row">
			<div class="col-lg-10">
				<div class="row">
					%s
				</div>
			</div>
			<div class="col-lg-2">
				<div class="bg-primary text-light"><b><center>STANDINGS (Last GW)</center></b></div>
				%s
				<hr class="hr"> 
				<div class="bg-warning text-light"><b><center>FIXTURES (Next GW)</center></b></div>
				%s
			</div>
		</div>
	</div>
</body>
</html>
`
const matchup_template string = `
<div class="row bg-success text-white">
		<div class="col-lg-6">
			%s
        </div>
   
        <div class="col-lg-6">
			%s
        </div>
</div>
<hr class="hr">
</br>
`
const player_template string = `
 			<div>
                %s
            </div>
            <div>
                %s
            </div>
`

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

func readDraftLive() Draft {
	// TODO: This is insecure; use only in dev environments.
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	defer client.CloseIdleConnections()

	req, err := http.NewRequest("GET",
		"https://draft.premierleague.com/api/league/29143/details",
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

	var draft Draft
	err = json.NewDecoder(resp.Body).Decode(&draft)

	if err != nil {
		fmt.Println("Error:", err)
	}
	return draft
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
	// draft := readDraft()
	draft := readDraftLive()

	clubs := map[int]Club{}
	owners := map[int]string{}
	for _, user := range draft.LeagueEntries {
		owners[user.ID] = user.PlayerFirstName
		clubs[user.ID] = getDraftClubs(uint32(user.EntryID), event)
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

	done := 0
	clubOrder := []int{}
	for _, entry := range draft.Matches {
		if entry.Event == int(event) {
			clubOrder = append(clubOrder, entry.LeagueEntry1, entry.LeagueEntry2)
			done += 1
		}
		if done == 3 {
			break
		}
	}

	first_team := true
	first_team_disp, second_team_disp := "", ""
	for _, clid := range clubOrder {
		club := clubs[clid]
		total := 0
		var table string

		table += `<table class="table table-condensed table-striped table-bordered">` +
			"<tr>" +
			"<em><th>PLAYER</th><th>TM</th><th>POS</th>" +
			"<th>MP</th><th>GS</th><th>AS</th><th>GA</th><th>YC</th><th>BO</th>" +
			"<th>PT</th></em></tr>"

		for i, pl := range club.Squad {
			player := players[uint16(pl.Element)]
			playerLiveStat := live[uint16(pl.Element)].Stats

			var row_style string
			if i >= 11 {
				if playerLiveStat.TotalPoints > 5 {
					row_style = ` class="table-danger" style="font-weight:bold"`
				} else if playerLiveStat.TotalPoints > 1 {
					row_style = `class="table-danger style="font-weight:italic"`
				} else {
					row_style = ` class="table-danger"`
				}
			} else if playerLiveStat.Minutes > 0 {
				row_style = `class="table-dark text-light"`
			}

			table += fmt.Sprintf(
				"<tr %s> <td>%s</td> <td>%s</td> <td>%s</td>"+
					"<td>%d</td> <td>%d</td> <td>%d</td> <td>%d</td> <td>%d</td> <td>%d</td><td>%d</td></tr>",
				row_style,
				player.WebName, TEAMS[player.Team-1], POS[player.ElementType-1],
				playerLiveStat.Minutes,
				playerLiveStat.GoalsScored,
				playerLiveStat.Assists,
				playerLiveStat.GoalsConceded,
				playerLiveStat.YellowCards,
				playerLiveStat.Bonus,
				playerLiveStat.TotalPoints)
			if i < 11 {
				total += playerLiveStat.TotalPoints
			}

		}
		table += "</table>"

		player_deets := `<div><b>` + owners[clid] + " [Total Points: " + strconv.Itoa(total) + "]</b></div>"

		if first_team {
			first_team_disp = fmt.Sprintf(player_template, player_deets, table)
			first_team = false
		} else {
			second_team_disp = fmt.Sprintf(player_template, player_deets, table)
			out += fmt.Sprintf(matchup_template, first_team_disp, second_team_disp)
			first_team = true
		}
	}

	standings := `<table class="table table-condensed table-striped table-bordered">
			<tr> <th>#</th><th>Player</th><th>W-D-L</th><th>PTS</th></tr>`
	st := draft.Standings
	sort.Slice(st, func(i, j int) bool {
		if st[i].Total == st[j].Total {
			return st[i].PointsFor-st[i].PointsAgainst > st[j].PointsFor-st[j].PointsAgainst
		}
		return st[i].Total > st[j].Total
	})

	for i, pos := range draft.Standings {
		standings += fmt.Sprintf("<tr><td>%d</td><td>%s</td><td>%d-%d-%d</td><td>%d</td></tr>",
			i+1, owners[pos.LeagueEntry], pos.MatchesWon, pos.MatchesDrawn, pos.MatchesLost, pos.Total)
	}
	standings += `</table>`

	clubOrder = []int{}
	done = 0
	if event < 38 {
		for _, entry := range draft.Matches {
			if entry.Event == int(event)+1 {
				clubOrder = append(clubOrder, entry.LeagueEntry1, entry.LeagueEntry2)
				done += 1
			}
			if done == 3 {
				break
			}
		}
	}

	var fixtures string
	if len(clubOrder) == 6 {
		fixtures = fmt.Sprintf(
			`<div class="bg-secondary text-light"> <center><div>%s VS %s</div><div>%s VS %s</div><div>%s VS %s</div></center></div><br><p><p>`,
			owners[clubOrder[0]], owners[clubOrder[1]], owners[clubOrder[2]],
			owners[clubOrder[3]], owners[clubOrder[4]], owners[clubOrder[5]])
	} else {
		fixtures = "Could Not Load"
	}

	html := fmt.Sprintf(site_template, getCurrentEvent(), out, standings, fixtures)
	return html
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, getOutput())
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("0.0.0.0:80", nil))
	//log.Fatal(http.ListenAndServeTLS("0.0.0.0:443", "/etc/letsencrypt/live/draftee.kparajuli.com/fullchain.crt", "/etc/letsencrypt/live/draftee.kparajuli.com/privkey.crt", nil))
}
