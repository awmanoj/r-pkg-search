package main 

import(
	"github.com/robfig/cron/v3"
)

func main() {
	c := cron.New()
	c.AddFunc("@every 1m", RunJob)
	c.Start()	

	
	select{}
}
