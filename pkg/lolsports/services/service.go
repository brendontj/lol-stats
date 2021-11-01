package services

import (
	"context"
	"fmt"
	lol_transformer "github.com/brendontj/lol-stats/pkg/lol-transformer"
	"github.com/brendontj/lol-stats/pkg/lolsports"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"time"
)

type Service interface {
	PopulateLeagues() error
	GetLeagues() ([]lolsports.League,error)
	GetEventsExternalRef() ([]*string, error)
	PopulateDBScheduleOfLeague(leagueExternalReference string) error
	PopulateDBWithEventDetail(eventExternalReference string) error
	GetGamesReference() ([]string,error)
	PopulateDBWithGameData(gameID string) error
	GetLiveGames() (*lolsports.EventsLiveData, error)
	GetTeamsHistoricalData(redTeamName, blueTeamName string) (lolsports.HistoricalData, error)
}

type lolService struct {
	storage          *pgxpool.Pool
	esportsApiClient lolsports.EsportsAPIScrapper
	feedApiClient    lolsports.FeedAPIScrapper
	DB Storage
}

func NewLolService(pgStorage *pgxpool.Pool, esportsApiClient lolsports.EsportsAPIScrapper, feedApiClient lolsports.FeedAPIScrapper) Service {
	return &lolService{
		storage:          pgStorage,
		esportsApiClient: esportsApiClient,
		feedApiClient:    feedApiClient,
		DB: Storage{pool: pgStorage},
	}
}

func (l *lolService) PopulateLeagues() error {
	leagueData, err := l.esportsApiClient.GetLeagues("pt-BR")
	if err != nil {
		return errors.Wrap(err, "[service error] unable to get leagues")
	}
	if err := l.DB.InsertLeagues(leagueData); err != nil {
		return errors.Wrap(err, " [service error] unable to insert leagues")
	}
	return nil
}

func (l *lolService) GetLeagues() ([]lolsports.League,error) {
	return l.DB.GetLeagues()
}

func (l *lolService) PopulateDBScheduleOfLeague(leagueExternalReference string) error {
	olderPage, err := l.populateWithMostRecentScheduleByLeague(leagueExternalReference)
	if err != nil {
		return errors.Wrap(err, "unable to populate with the most recent schedule")
	}

	if olderPage != lolsports.EmptyField {
		for {
			scheduleData, err := l.esportsApiClient.GetSchedule("pt-BR", leagueExternalReference, olderPage)
			if err != nil {
				return errors.Wrapf(err, "[service error] unable to get schedules with older page = %s",olderPage)
			}
			op, err := l.saveScheduleContent(scheduleData.Data.Schedule)
			if err != nil {
				return err
			}

			if op == lolsports.EmptyField {
				break
			}
			olderPage = op
		}
	}
	return nil
}

func (l *lolService) PopulateDBWithEventDetail(eventExternalReference string) error {
	exists, err := l.DB.ExistsEventExternalRef(eventExternalReference)
	if err != nil {
		return errors.Wrapf(err, "[service error] unable to verify if event was sync: %v", eventExternalReference)
	}

	if exists {
		return nil
	}

	eventDetail, err := l.esportsApiClient.GetEventDetail("pt-BR", eventExternalReference)
	if err != nil {
		return errors.Wrapf(err, "[service error] unable to get event detail for event external reference: %v", eventExternalReference)
	}

	if eventDetail.Data.Event.ID == lolsports.EmptyField || eventDetail.Data.Event.Tournament.TournamentID == lolsports.EmptyField || eventDetail.Data.Event.League.ID == lolsports.EmptyField {
		return nil
	}

	if err := l.DB.SaveEventDetail(eventDetail, eventExternalReference); err != nil {
		return errors.Wrapf(err, "[service error] unable to save event detail for event with external reference: %v", eventExternalReference)
	}

	if err = l.DB.SaveGameInfo(eventDetail); err != nil {
		return errors.Wrapf(err, "[service error] unable to save game event for event with external reference: %v", eventExternalReference)
	}
	return nil
}

func (l *lolService) populateWithMostRecentScheduleByLeague(leagueExternalReference string) (string, error) {
	scheduleData, err := l.esportsApiClient.GetSchedule("pt-BR", leagueExternalReference, lolsports.EmptyField)
	if err != nil {
		return lolsports.EmptyField, errors.Wrap(err, "[service error] unable to get leagues")
	}

	s, err := l.saveScheduleContent(scheduleData.Data.Schedule)
	if err != nil {
		return s, err
	}
	return scheduleData.Data.Schedule.Pages.Older, nil
}

func (l *lolService) saveScheduleContent(scheduleContent lolsports.Schedule) (string, error) {
	for _, event := range scheduleContent.Events {
		time.Sleep(20 * time.Microsecond)
		exists, err := l.DB.ExistsEvent(event.Match.ID)
		if err != nil {
			return lolsports.EmptyField, err
		}

		if exists {
			continue
		}

		err = l.DB.SaveEvent(event)
		if err != nil {
			return lolsports.EmptyField, err
		}
	}
	return scheduleContent.Pages.Older, nil
}

func (l *lolService) GetEventsExternalRef() ([]*string,error) {
	queryGetAllEventsRef := `
SELECT 
	e.external_reference
FROM schedule.matches e
where e.state = 'completed';
`
	rows, err := l.storage.Query(context.Background(), queryGetAllEventsRef)
	if err != nil {
		return nil, errors.Wrap(err,"unable to get events from storage")
	}
	defer rows.Close()

	var eventIDs []*string
	for rows.Next() {
		var s *string
		err = rows.Scan(&s)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to scan event with external reference = %s", s)
		}
		eventIDs = append(eventIDs, s)
	}

	return eventIDs, nil
}

func (l *lolService) GetGamesReference() ([]string, error) {
	queryGetAllGames := `
SELECT DISTINCT
	game_external_ref
FROM schedule.events_games
WHERE status = 'completed'
`
	rows, err := l.storage.Query(context.Background(), queryGetAllGames)
	if err != nil {
		return nil, errors.Wrap(err,"unable to get games from storage")
	}
	defer rows.Close()

	var gamesReference []string
	for rows.Next() {
		var g string
		err = rows.Scan(&g)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to scan game reference with external reference = %s", g)
		}
		gamesReference = append(gamesReference, g)
	}

	return gamesReference, nil
}

func (l *lolService) PopulateDBWithGameData(gameID string) error {
	exists, err := l.DB.ExistsGameID(gameID)
	if err != nil {
		return errors.Wrapf(err, "[Service error] unable to verify if gameID was registered: %v", gameID)
	}

	if exists {
		return nil
	}

	firstFrame, err := l.getFirstFrameOfMatchGame(gameID)
	if err != nil {
		return err
	}

	if firstFrame == nil || len(firstFrame.Participants) == 0  {
		return nil
	}

	y, m, d := firstFrame.Rfc460Timestamp.Date()
	h := firstFrame.Rfc460Timestamp.Hour()
	min := firstFrame.Rfc460Timestamp.Minute()

	startTime := time.Date(y, m, d, h, min, 0, 0, time.UTC)

	gameDetails, err := l.feedApiClient.GetDetailsFromLiveMatch(gameID, startTime)
	if err != nil {
		return err
	}
	if gameDetails == nil {
		return nil
	}

	gameSystemID, err := l.saveGame(*gameDetails)
	if err != nil {
		return err
	}
	currentTime := startTime
	for {
		liveMatchData, err := l.feedApiClient.GetDataFromLiveMatch(gameID,currentTime)
		if err != nil {
			return err
		}

		if liveMatchData == nil || len(liveMatchData.Frames) == 0  {
			return nil
		}

		for _, p := range liveMatchData.Frames[0].Participants {
				if err := l.saveParticipantMetadata(gameSystemID, gameID, p, currentTime); err != nil {
					return err
				}
		}

		detailsLiveMatch, err := l.feedApiClient.GetDetailsFromLiveMatch(gameID, currentTime)
		if err != nil {
			return err
		}

		if err := l.saveGameMetadata(gameSystemID, *detailsLiveMatch, currentTime); err != nil {
			return err
		}

		if detailsLiveMatch.Frames[0].GameState == "finished" || currentTime.After(startTime.Add(time.Minute * 25)){
			break
		}
		currentTime = currentTime.Add(time.Minute * 5)
	}
	return nil
}

func (l *lolService) getFirstFrameOfMatchGame(gameID string) (*lolsports.Frames, error) {
	timeOfBegin := time.Date(1950,1,1,0,0,0,0,time.UTC)
	liveMatch, err := l.feedApiClient.GetDataFromLiveMatch(gameID, timeOfBegin)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to get data from live match, game reference = %s", gameID)
	}

	if liveMatch == nil || len(liveMatch.Frames) == 0 {
		return nil, nil
	}

	for _, f := range liveMatch.Frames {
		if f.Participants[0].TotalGoldEarned != 0 {
			return &f, nil
		}
	}

	y, m, d := liveMatch.Frames[len(liveMatch.Frames)-1].Rfc460Timestamp.Date()
	h := liveMatch.Frames[len(liveMatch.Frames)-1].Rfc460Timestamp.Hour()
	min := liveMatch.Frames[len(liveMatch.Frames)-1].Rfc460Timestamp.Minute()

	timeOfBegin = time.Date(y, m, d, h, min, 0, 0, time.UTC)
	for {
		liveMatch, err = l.feedApiClient.GetDataFromLiveMatch(gameID, timeOfBegin)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to get data from live match, game reference = %s", gameID)
		}

		if liveMatch == nil || len(liveMatch.Frames) == 0{
				return nil, nil
		}

		for _, f := range liveMatch.Frames {
			if f.Participants[0].TotalGoldEarned != 0 {
				return &f, nil
			}
		}

		timeOfBegin = timeOfBegin.Add(60 * time.Second)
	}
}

func (l *lolService) saveGame(gameDetail lolsports.LiveMatchDetailData) (uuid.UUID, error) {
	tx, err := l.storage.Begin(context.Background())
	if err != nil {
		return uuid.Nil, err
	}

	queryInsertGameMetadata :=
		`INSERT INTO game.games 
(id, gameID, matchID, patch_version, blueTeamID, redTeamID) 
VALUES ($1, $2, $3, $4, $5, $6);`

	gameID := uuid.NewV4()
	_, err = tx.Exec(
		context.Background(),
		queryInsertGameMetadata,
		gameID,
		gameDetail.EsportsGameID,
		gameDetail.EsportsMatchID,
		gameDetail.GameMetadata.PatchVersion,
		gameDetail.GameMetadata.BlueTeamMetadata.EsportsTeamID,
		gameDetail.GameMetadata.RedTeamMetadata.EsportsTeamID)
	if err != nil {
		_ = tx.Rollback(context.Background())
		return uuid.Nil, err
	}

	queryInsertGameParticipantsInfo :=
		`INSERT INTO game.participants_info 
(gameID, participantID, championID, esportsPlayerID, summonerName, role) 
VALUES ($1, $2, $3, $4, $5, $6);`

		for _, b := range gameDetail.GameMetadata.BlueTeamMetadata.ParticipantMetadata {
			_, err := tx.Exec(
				context.Background(),
				queryInsertGameParticipantsInfo,
				gameID,
				b.ParticipantID,
				b.ChampionID,
				b.EsportsPlayerID,
				b.SummonerName,
				b.Role)
			if err != nil {
				_ = tx.Rollback(context.Background())
				return uuid.Nil, err
			}
		}

	for _, b := range gameDetail.GameMetadata.RedTeamMetadata.ParticipantMetadata {
		_, err := tx.Exec(
			context.Background(),
			queryInsertGameParticipantsInfo,
			gameID,
			b.ParticipantID,
			b.ChampionID,
			b.EsportsPlayerID,
			b.SummonerName,
			b.Role)
		if err != nil {
			_ = tx.Rollback(context.Background())
			return uuid.Nil, err
		}
	}
	_ = tx.Commit(context.Background())
	return gameID, nil
}

func (l *lolService) saveParticipantMetadata(gameID uuid.UUID, gameExternalID string, participantData lolsports.Participants, currentTime time.Time) error {
	queryInsertParticipantsStatsMetadata :=
		`INSERT INTO game.participants_stats
(gameID, game_externalID, participantID, game_timestamp, level, kills, deaths, assists, total_gold_earned, creep_score, kill_participation, champion_damage_share, wards_placed, wards_destroyed) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14);`

	_, err := l.storage.Exec(
		context.Background(),
		queryInsertParticipantsStatsMetadata,
		gameID,
		gameExternalID,
		participantData.ParticipantID,
		currentTime,
		participantData.Level,
		participantData.Kills,
		participantData.Deaths,
		participantData.Assists,
		participantData.TotalGoldEarned,
		participantData.CreepScore,
		participantData.KillParticipation,
		participantData.ChampionDamageShare,
		participantData.WardsPlaced,
		participantData.WardsDestroyed)
	if err != nil {
		return  err
	}
	return nil
}

func (l *lolService) saveGameMetadata(gameID uuid.UUID, liveGameDetail lolsports.LiveMatchDetailData, currentTime time.Time) error {
	queryInsertGameMetadata :=
		`INSERT INTO game.games_stats
(gameID, timestamp, gameState, blueTeamID, redTeamID, blue_team_total_gold, blue_team_inhibitors, blue_team_towers, blue_team_barons, blue_team_total_kills, blue_team_dragons, red_team_total_gold, red_team_inhibitors, red_team_towers, red_team_barons, red_team_total_kills, red_team_dragons) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17);`

	_, err := l.storage.Exec(
		context.Background(),
		queryInsertGameMetadata,
		gameID,
		currentTime,
		liveGameDetail.Frames[0].GameState,
		liveGameDetail.GameMetadata.BlueTeamMetadata.EsportsTeamID,
		liveGameDetail.GameMetadata.RedTeamMetadata.EsportsTeamID,
		liveGameDetail.Frames[0].BlueTeam.TotalGold,
		liveGameDetail.Frames[0].BlueTeam.Inhibitors,
		liveGameDetail.Frames[0].BlueTeam.Towers,
		liveGameDetail.Frames[0].BlueTeam.Barons,
		liveGameDetail.Frames[0].BlueTeam.TotalKills,
		liveGameDetail.Frames[0].BlueTeam.Dragons,
		liveGameDetail.Frames[0].RedTeam.TotalGold,
		liveGameDetail.Frames[0].RedTeam.Inhibitors,
		liveGameDetail.Frames[0].RedTeam.Towers,
		liveGameDetail.Frames[0].RedTeam.Barons,
		liveGameDetail.Frames[0].RedTeam.TotalKills,
		liveGameDetail.Frames[0].RedTeam.Dragons)
	if err != nil {
		return  err
	}
	return nil
}

func (l *lolService) GetLiveGames() (*lolsports.EventsLiveData, error) {
	data, err := l.esportsApiClient.GetEventsLive("pt-BR")
	if err != nil {
		return nil, fmt.Errorf("[SERVICE ERROR] unable to get live games, err: %v", err)
	}
	return data, nil
}

func (l *lolService) GetTeamsHistoricalData(teamRedName, teamBlueName string) (lolsports.HistoricalData, error) {
	fr3RedTeam, err := l.getFormRatio(teamRedName, 3)
	if err != nil {
		return nil, err
	}

	fr5RedTeam, err := l.getFormRatio(teamRedName, 5)
	if err != nil {
		return nil, err
	}

	fr3BlueTeam, err := l.getFormRatio(teamBlueName, 3)
	if err != nil {
		return nil, err
	}

	fr5BlueTeam, err := l.getFormRatio(teamBlueName, 5)
	if err != nil {
		return nil, err
	}

}

func (l *lolService) getFormRatio(teamName string, numberOfPastGames int) (float64, error) {
	panic("implement me")
}

func (l *lolService) getPastGamesStats(teamName string, numberOfPastGames int, gameMoment int) (lol_transformer.StatsInfo, error) {
	panic("implement me")
}