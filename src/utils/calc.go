package utils

func calculateUsageHourly(period float64, days int, usagePerHour float64) float64 {
	hours := days * 24 // convert days to hours
	totalUsage := period * float64(hours) * usagePerHour
	return totalUsage
}

func calculateUsageGB(bandwidth float64, days int, usagePerGB float64) float64 {
	totalUsage := bandwidth * float64(days) * usagePerGB
	return totalUsage
}
