package queries

func ExtractionQuery() string {
	return `SELECT DISTINCT
                g.gameid as game_external_ref,
                matchid as match_external_ref,
                patch_version as patch_version,
                g.blueteamid as blue_team_ref,
                g.redteamid as red_team_ref,
                gamestate as game_statem,
                blue_team_total_gold,
                blue_team_inhibitors,
                blue_team_towers,
                blue_team_barons,
                blue_team_total_kills,
                blue_team_dragons,
                red_team_total_gold,
                red_team_inhibitors,
                red_team_towers,
                red_team_barons,
                red_team_total_kills,
                red_team_dragons,
                pi.participantid as participant_id,
                championid as champion_name,
                esportsplayerid as esports_player_ref,
                summonername as summoner_name,
                role,
                game_timestamp,
                level,
                kills,
                deaths,
                assists,
                total_gold_earned,
                creep_score,
                kill_participation,
                champion_damage_share,
                wards_placed,
                wards_destroyed,
                eg.event_external_ref as event_external_ref,
                game_number,
                status,
                team_a_external_ref,
                team_b_external_ref,
                team_a_side,
                team_b_side,
                m.external_reference as match_external_ref,
                team_a_name,
                team_a_code,
                team_b_name,
                team_b_code,
                team_a_record_wins,
                team_a_record_losses,
                team_b_record_wins,
                team_b_record_losses,
                team_a_game_wins,
                team_b_game_wins,
                best_of,
                state,
                league_name,
                tournament_external_ref,
                league_external_ref,
                region
FROM game.participants_stats ps
    LEFT JOIN game.games_stats gs
        ON gs.gameid = ps.gameid
    LEFT JOIN game.participants_info pi
        ON pi.gameid = gs.gameid
    LEFT JOIN game.games g
        ON g.id = ps.gameid
    LEFT JOIN schedule.events_games eg
        ON g.gameid = eg.game_external_ref
    LEFT JOIN schedule.matches m
        ON eg.event_external_ref = m.external_reference
    LEFT JOIN schedule.events_detail e
        ON m.external_reference = e.event_external_ref
    LEFT JOIN league.leagues l
        ON e.league_external_ref = l.external_reference
    WHERE eg.status = 'completed'
        and ps.gameid = g.id
        and gs.timestamp = ps.game_timestamp
        and ps.game_externalid = eg.game_external_ref
        and pi.participantid = ps.participantid
order by game_external_ref, game_number, game_timestamp, pi.participantid;`
}