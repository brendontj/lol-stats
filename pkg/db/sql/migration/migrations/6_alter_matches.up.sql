ALTER TABLE schedule.matches
    ADD COLUMN team_a_5_form_ratio NUMERIC (4,3),
    ADD COLUMN team_a_3_form_ratio NUMERIC (4,3),
    ADD COLUMN team_b_5_form_ratio NUMERIC (4,3),
    ADD COLUMN team_b_3_form_ratio NUMERIC (4,3),
    ADD COLUMN team_a_5_gold_total_mean_at15 NUMERIC (9,2),
    ADD COLUMN team_a_3_gold_total_mean_at15 NUMERIC (9,2),
    ADD COLUMN team_b_5_gold_total_mean_at15 NUMERIC (9,2),
    ADD COLUMN team_b_3_gold_total_mean_at15 NUMERIC (9,2),
    ADD COLUMN team_a_5_kills_mean_at15 NUMERIC (3,1),
    ADD COLUMN team_a_3_kills_mean_at15 NUMERIC (3,1),
    ADD COLUMN team_b_5_kills_mean_at15 NUMERIC (3,1),
    ADD COLUMN team_b_3_kills_mean_at15 NUMERIC (3,1),
    ADD COLUMN team_a_5_dragons_mean_at15 NUMERIC (2,1),
    ADD COLUMN team_a_3_dragons_mean_at15 NUMERIC (2,1),
    ADD COLUMN team_b_5_dragons_mean_at15 NUMERIC (2,1),
    ADD COLUMN team_b_3_dragons_mean_at15 NUMERIC (2,1),
    ADD COLUMN team_a_5_towers_mean_at15 NUMERIC (3,1),
    ADD COLUMN team_a_3_towers_mean_at15 NUMERIC (3,1),
    ADD COLUMN team_b_5_towers_mean_at15 NUMERIC (3,1),
    ADD COLUMN team_b_3_towers_mean_at15 NUMERIC (3,1),
    ADD COLUMN team_a_5_gold_total_mean_at25 NUMERIC (9,2),
    ADD COLUMN team_a_3_gold_total_mean_at25 NUMERIC (9,2),
    ADD COLUMN team_b_5_gold_total_mean_at25 NUMERIC (9,2),
    ADD COLUMN team_b_3_gold_total_mean_at25 NUMERIC (9,2),
    ADD COLUMN team_a_5_kills_mean_at25 NUMERIC (3,1),
    ADD COLUMN team_a_3_kills_mean_at25 NUMERIC (3,1),
    ADD COLUMN team_b_5_kills_mean_at25 NUMERIC (3,1),
    ADD COLUMN team_b_3_kills_mean_at25 NUMERIC (3,1),
    ADD COLUMN team_a_5_dragons_mean_at25 NUMERIC (2,1),
    ADD COLUMN team_a_3_dragons_mean_at25 NUMERIC (2,1),
    ADD COLUMN team_b_5_dragons_mean_at25 NUMERIC (2,1),
    ADD COLUMN team_b_3_dragons_mean_at25 NUMERIC (2,1),
    ADD COLUMN team_a_5_barons_mean_at25 NUMERIC (2,1),
    ADD COLUMN team_a_3_barons_mean_at25 NUMERIC (2,1),
    ADD COLUMN team_b_5_barons_mean_at25 NUMERIC (2,1),
    ADD COLUMN team_b_3_barons_mean_at25 NUMERIC (2,1),
    ADD COLUMN team_a_5_towers_mean_at25 NUMERIC (3,1),
    ADD COLUMN team_a_3_towers_mean_at25 NUMERIC (3,1),
    ADD COLUMN team_b_5_towers_mean_at25 NUMERIC (3,1),
    ADD COLUMN team_b_3_towers_mean_at25 NUMERIC (3,1),
    ADD COLUMN team_a_5_inhibitors_mean_at25 NUMERIC (2,1),
    ADD COLUMN team_a_3_inhibitors_mean_at25 NUMERIC (2,1),
    ADD COLUMN team_b_5_inhibitors_mean_at25 NUMERIC (2,1),
    ADD COLUMN team_b_3_inhibitors_mean_at25 NUMERIC (2,1),
    ADD COLUMN team_a_5_inhibitors_mean_at15 NUMERIC (2,1),
    ADD COLUMN team_a_3_inhibitors_mean_at15 NUMERIC (2,1),
    ADD COLUMN team_b_5_inhibitors_mean_at15 NUMERIC (2,1),
    ADD COLUMN team_b_3_inhibitors_mean_at15 NUMERIC (2,1);





