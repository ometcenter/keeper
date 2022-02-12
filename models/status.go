package models

import "strings"

type JobsStatussGroupSwap struct {
	Base          string
	ExchangeJobID string
}

type JobsStatusGroup struct {
	Areas         string
	Base          string
	ExchangeJobID string
	AmountAreas   int
}

func (JSB *JobsStatusGroup) GetJobsStatusGroupSlice(SettingsJobSliceQueryToBI SettingsJobSliceQueryToBI,
	Areas []string) ([]JobsStatusGroup, error) {

	MapSwap := make(map[string]JobsStatussGroupSwap)

	for _, SliceQueryToBI := range SettingsJobSliceQueryToBI.SliceQueryToBI {

		stringSlice := strings.Split(SliceQueryToBI.Query[0].Area, ",")
		for _, Area := range stringSlice {
			AreaFinde := false
			for _, AreaProblem := range Areas {
				if Area == AreaProblem {
					AreaFinde = true
					break
				}
			}
			if !AreaFinde {
				continue
			}
			var GetJobsStatusProblemsGroupSwap JobsStatussGroupSwap
			GetJobsStatusProblemsGroupSwap.Base = SliceQueryToBI.Query[0].Base
			GetJobsStatusProblemsGroupSwap.ExchangeJobID = SliceQueryToBI.Query[0].ExchangeJobID
			MapSwap[Area] = GetJobsStatusProblemsGroupSwap
		}

		// fmt.Println("Query len", len(SliceQueryToBI.Query), " - ", SliceQueryToBI.Query[0].Area, " - ",
		// 	SliceQueryToBI.Query[0].Base, SliceQueryToBI.Query[0].ExchangeJobID, " - ")
	}

	ResultMap := make(map[JobsStatussGroupSwap][]string)
	for key, value := range MapSwap {

		Records, ok := ResultMap[value]
		if ok != true {
			var NewRecord []string
			NewRecord = append(NewRecord, key)
			ResultMap[value] = NewRecord
		} else {
			Records = append(Records, key)
			ResultMap[value] = Records
		}

	}

	var JobsStatusGroupSlice []JobsStatusGroup
	for key, value := range ResultMap {
		resultArea := strings.Join(value, ",")
		JobsStatusGroup := JobsStatusGroup{
			Areas:         resultArea,
			Base:          key.Base,
			ExchangeJobID: key.ExchangeJobID,
			AmountAreas:   len(value),
		}
		//fmt.Printf("Base: %s, ExchangeJobID %s, areas: %s \n", key.Base, key.ExchangeJobID, resultArea)
		JobsStatusGroupSlice = append(JobsStatusGroupSlice, JobsStatusGroup)
	}

	return JobsStatusGroupSlice, nil

}
