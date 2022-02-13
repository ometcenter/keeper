package utility

import "database/sql"

func GetAreasByStasus(DB *sql.DB, JobID, Stasus string) ([]string, error) {

	var argsquery []interface{}
	argsquery = append(argsquery, JobID)
	argsquery = append(argsquery, Stasus)

	QueryText := `select
		area
	from
		public.exchange_jobs
	where
		job_id = $1
		and "event" = $2;`

	rows, err := DB.Query(QueryText, argsquery...)
	if err != nil {
		return nil, err
	}

	var AreasForReturn []string
	for rows.Next() {
		var area string
		err = rows.Scan(&area)
		if err != nil {
			return nil, err
		}
		AreasForReturn = append(AreasForReturn, area)
	}

	defer rows.Close()

	return AreasForReturn, nil

}

func GetAreasNotEqualToStatus(DB *sql.DB, JobID, Stasus string) ([]string, error) {

	var argsquery []interface{}
	argsquery = append(argsquery, JobID)
	argsquery = append(argsquery, Stasus)

	QueryText := `SELECT area
	FROM public.exchange_jobs where job_id = $1 and "event" <> $2;`

	rows, err := DB.Query(QueryText, argsquery...)
	if err != nil {
		return nil, err
	}

	var AreasForReturn []string
	for rows.Next() {
		var area string
		err = rows.Scan(&area)
		if err != nil {
			return nil, err
		}
		AreasForReturn = append(AreasForReturn, area)
	}

	defer rows.Close()

	return AreasForReturn, nil

}