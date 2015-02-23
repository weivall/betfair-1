package main

import (
	"betfair"
	"fmt"
	"os"
)

func certLogin() *betfair.Session {
	c, err := betfair.NewCredentials("<USERNAME>", "<PASS>", "UK",
		"client-2048.crt",
		"client-2048.key")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	session, err := betfair.NewSession(c)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := session.SetUsedApplication("<APPNAME>", true); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return session
}

func restLogin() *betfair.Session {
	c, err := betfair.NewCredentials("<USERNAME>", "<PASS>", "UK",
		"<APPKEY>")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	session, err := betfair.NewSession(c)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return session
}

func main() {

	session := certLogin()

	/*query := &betfair.Query{
		MarketFilter: &betfair.MarketFilter{
			EventIds: []string{"27327279"},
		},
		MaxResults: 1000,
		MarketProjection: []string{"COMPETITION", "EVENT",
			"RUNNER_DESCRIPTION", "RUNNER_METADATA"},
		Locale: "en",
	}

	_, err := session.ListMarketCatalogue(query, func(s *betfair.Session,
		q *betfair.Query, v interface{}) {
		x := *v.(*[]betfair.MarketCatalogue)
		for _, item := range x {
			fmt.Println("\t", item.Competition.Name, item.Event.Name, item.MarketId,
				item.MarketName, item.TotalMatched)
		}
	})*/

	acc, _ := session.GetAccountDetails()
	fmt.Println(acc)

	query := &betfair.Query{
		MaxResults: 1000,
		MarketIds:  []string{"1.116734361"},
		Locale:     "en",
		PriceProjection: &betfair.PriceProjection{
			PriceData: []string{"EX_ALL_OFFERS"},
		},
	}

	_, err := session.ListMarketBook(query, func(s *betfair.Session,
		q *betfair.Query, v interface{}) {
		x := *v.(*[]betfair.MarketBook)
		for _, item := range x {
			fmt.Println("\t", item.MarketId, item.IsMarketDataDelayed,
				item.Inplay, item.NumberOfActiveRunners, item.TotalMatched,
				item.TotalAvailable, item.Version)

			for _, r := range item.Runners {
				fmt.Println("\t\t", r.SelectionId, r.LastPriceTraded,
					r.TotalMatched, r.Ex.AvailableToLay)
			}
		}
	})
	if err != nil {
		fmt.Println(err)
		return
	}
}
