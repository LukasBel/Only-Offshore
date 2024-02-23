package Handlers

import "github.com/LukasBel/Only-Offshore.git/Models"

func Percent(before, after *Models.Stats) []float64 {
	percents := []float64{}
	percents = append(percents, float64((after.FlatBench-before.FlatBench)/before.FlatBench))
	percents = append(percents, float64((after.InclineBench-before.InclineBench)/before.InclineBench))
	percents = append(percents, float64((after.Squat-before.Squat)/before.Squat))
	percents = append(percents, float64((after.PullUps-before.PullUps)/before.PullUps))
	percents = append(percents, float64((after.WeightedPullUpMax-before.WeightedPullUpMax)/before.WeightedPullUpMax))
	percents = append(percents, float64((after.BodyWeight-before.BodyWeight)/before.BodyWeight))

	return percents
}
