package services

import (
	"context"
	"encoding/csv"
	"fmt"
	"github.com/brendontj/lol-stats/pkg/lol-transformer/queries"
	"github.com/jackc/pgx/v4/pgxpool"
	"os"
	"strconv"
	"time"
)

type Service interface {
	Extract(outputPath, prefixFilename string) error
}

type service struct {
	storage          *pgxpool.Pool
	DB Storage
}

func NewLolService(pgStorage *pgxpool.Pool) Service {
	return &service{
		storage:          pgStorage,
		DB: Storage{pool: pgStorage},
	}
}

func (s service) Extract(outputPath, prefixFilename string) error {
	date := time.Now()

	rows, err := s.storage.Query(context.Background(), queries.ExtractionQuery())
	if err != nil {
		return fmt.Errorf("[service error] Error extracting data: Date %v", date.GoString())
	}

	f, err := os.Create(outputPath + prefixFilename + date.Local().String() + ".csv")
	if err != nil {
		return fmt.Errorf("[service error] Error creating file")
	}
	defer  f.Close()

	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}

	w := csv.NewWriter(f)
	defer w.Flush()
	if err := w.Write([]string{
		"game_external_ref",
		"match_external_ref",
		"patch_version",
		"blue_team_ref",
		"red_team_ref",
		"game_statem",
		"blue_team_total_gold",
		"blue_team_inhibitors",
		"blue_team_towers",
		"blue_team_barons",
		"blue_team_total_kills",
		"blue_team_dragons",
		"red_team_total_gold",
		"red_team_inhibitors",
		"red_team_towers",
		"red_team_barons",
		"red_team_total_kills",
		"red_team_dragons",
		"participant_id",
		"champion_name",
		"esports_player_ref",
		"summoner_name",
		"role",
		"game_timestamp",
		"level",
		"kills",
		"deaths",
		"assists",
		"total_gold_earned",
		"creep_score",
		"kill_participation",
		"champion_damage_share",
		"wards_placed",
		"wards_destroyed",
		"event_external_ref",
		"game_number",
		"status",
		"team_a_external_ref",
		"team_b_external_ref",
		"team_a_side",
		"team_b_side",
		"match_external_ref",
		"team_a_name",
		"team_a_code",
		"team_b_name",
		"team_b_code",
		"team_a_record_wins",
		"team_a_record_losses",
		"team_b_record_wins",
		"team_b_record_losses",
		"team_a_game_wins",
		"team_b_game_wins",
		"best_of",
		"state",
		"league_name",
		"tournament_external_ref",
		"league_external_ref",
		"region"}); err != nil {
		return err
	}

	for rows.Next() {
		values, _ := rows.Values()
		s := make([]string, len(values))
		for i, v := range values {
			switch o := v.(type) {
			case bool:
				s[i] = strconv.FormatBool(o)
			case int:
				s[i] = strconv.Itoa(o)
			case time.Time:
				s[i] = fmt.Sprint(o.String())
			case string, int32:
				s[i] = fmt.Sprint(o)
			default:
				fmt.Printf("Need to implement this type convertion: %v", o)
				s[i] = fmt.Sprint(o)
			}
			s[i] = fmt.Sprint(v)
		}
		if err := w.Write(s); err != nil {
			return fmt.Errorf("error writing record to file: %v", err)
		}
	}
	return nil
}