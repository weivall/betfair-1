package betfair

import (
	"encoding/json"
	"errors"
	"strings"
	"time"
)

type TimeRange struct {
	From time.Time `json:"from,omitempty"`
	To   time.Time `json:"to,omitempty"`
}

type ExBestOffersOverrides struct {
	BestPricesDepth          int     `json:"bestPricesDepth,omitempty"`
	RollupModel              string  `json:"rollupModel,omitempty"`
	RollupLimit              int     `json:"rollupLimit,omitempty"`
	RollupLiabilityThreshold float64 `json:"rollupLiabilityThreshold,omitempty"`
	RollupLiabilityFactor    int     `json:"rollupLiabilityFactor,omitempty"`
}

type PriceProjection struct {
	PriceData             []string               `json:"priceData,omitempty"`
	ExBestOffersOverrides *ExBestOffersOverrides `json:"exBestOffersOverrides,omitempty"`
	Virtualise            bool                   `json:"virtualise,omitempty"`
	RolloverStakes        bool                   `json:"rolloverStakes,omitempty"`
}

type MarketFilter struct {
	TextQuery          string     `json:"textQuery,omitempty"`
	ExchangeIds        []string   `json:"exchangeIds,omitempty"`
	EventTypeIds       []string   `json:"eventTypeIds,omitempty"`
	EventIds           []string   `json:"eventIds,omitempty"`
	CompetitionIds     []string   `json:"competitionIds,omitempty"`
	MarketIds          []string   `json:"marketIds,omitempty"`
	Venues             []string   `json:"venues,omitempty"`
	BspOnly            bool       `json:"bspOnly,omitempty"`
	TurnInPlayEnabled  bool       `json:"turnInPlayEnabled,omitempty"`
	InPlayOnly         bool       `json:"inPlayOnly,omitempty"`
	MarketBettingTypes []string   `json:"marketBettingTypes,omitempty"`
	MarketCountries    []string   `json:"marketCountries,omitempty"`
	MarketTypeCodes    []string   `json:"marketTypeCodes,omitempty"`
	MarketStartTime    *TimeRange `json:"marketStartTime,omitempty"`
	WithOrders         []string   `json:"withOrders,omitempty"`
}

type Query struct {
	MarketFilter       *MarketFilter    `json:"filter,omitempty"`
	Locale             string           `json:"locale,omitempty"`
	MarketProjection   []string         `json:"marketProjection,omitempty"`
	MarketSort         []string         `json:"sort,omitempty"`
	MaxResults         uint16           `json:"maxResults,omitempty"`
	MarketIds          []string         `json:"marketIds,omitempty"`
	OrderProjection    string           `json:"orderProjection,omitempty"`
	MatchProjection    string           `json:"matchProjection,omitempty"`
	IncludeSettledBets bool             `json:"includeSettledBets,omitempty"`
	IncludeBspBets     bool             `json:"includeBspBets,omitempty"`
	NetOfCommission    bool             `json:"netOfCommission,omitempty"`
	PriceProjection    *PriceProjection `json:"priceProjection,omitempty"`
	CurrencyCode       string           `json:"currencyCode,omitempty"`
}

// Visitor Function type
/*
After reached result set, visitor function will be called for each item

s *Session - betfair Session pointer

q *Query - betfair Query pointer

v interface{} - Result set passed as interface
*/
type VisitorFunc func(s *Session, q *Query, v interface{})

// Event Type
type EventType struct {
	Id   string
	Name string
}

// Event Result Struct
type EventTypeResult struct {
	EventType   EventType
	MarketCount int
}

// Country Result
type CountryCodeResult struct {
	CountryCode string
	MarketCount int
}

// Event
type Event struct {
	Id          string
	Name        string
	CountryCode string
	Timezone    string
	Venue       string
	OpenDate    time.Time
}

// Event Result
type EventResult struct {
	Event       Event
	MarketCount int
}

// Competition
type Competition struct {
	Id   string
	Name string
}

// Competition Result
type CompetitionResult struct {
	Competition       Competition
	MarketCount       int
	CompetitionRegion string
}

// Market Type Result
type MarketTypeResult struct {
	MarketType  string
	MarketCount int
}

// Venue Result
type VenueResult struct {
	Venue       string
	MarketCount int
}

// MarketProfitAndLoss Result
type MarketProfitAndLoss struct {
	MarketId          string
	CommissionApplied float64
	ProfitAndLosses   []struct {
		SelectionId uint32
		IfWin       float64
		IfLose      float64
	}
}

type PriceSize struct {
	Price float64
	Size  float64
}

// Market Book
type MarketBook struct {
	MarketId              string
	IsMarketDataDelayed   bool
	Status                string
	BetDelay              int
	BspReconciled         bool
	Complete              bool
	Inplay                bool
	NumberOfWinners       int
	NumberOfRunners       int
	NumberOfActiveRunners int
	LastMatchTime         time.Time
	TotalMatched          float64
	TotalAvailable        float64
	CrossMatching         bool
	RunnersVoidable       bool
	Version               uint32
	Runners               []struct {
		SelectionId      uint32
		Handicap         float64
		Status           string
		AdjustmentFactor float64
		LastPriceTraded  float64
		TotalMatched     float64
		RemovalDate      time.Time
		Sp               struct {
			NearPrice         float64
			FarPrice          float64
			BackStakeTaken    []PriceSize
			layLiabilityTaken []PriceSize
			ActualSP          float64
		}
		Ex struct {
			AvailableToBack []PriceSize
			AvailableToLay  []PriceSize
			TradedVolume    []PriceSize
		}
		Orders []struct {
			PriceSize
			BetId           string
			OrderType       string
			Status          string
			PersistenceType string
			Side            string
			BspLiability    float64
			PlacedDate      time.Time
			AvgPriceMatched float64
			SizeMatched     float64
			SizeRemaining   float64
			SizeLapsed      float64
			SizeCancelled   float64
			SizeVoided      float64
		}
		Matches []struct {
			PriceSize
			BetId     string
			MatchId   string
			Side      string
			MatchDate time.Time
		}
	}
}

// Market Catalogue
type MarketCatalogue struct {
	MarketId        string
	MarketName      string
	MarketStartTime time.Time
	Description     struct{}
	TotalMatched    float64
	Runners         []struct {
		SelectionId  uint32
		RunnerName   string
		Handicap     float64
		SortPriority int
		Metadata     map[string]string
	}
	EventType   EventType
	Competition Competition
	Event       Event
}

// Returns event types as []EventResult or error if occured
func (s *Session) ListEventTypes(q *Query, fn ...VisitorFunc) ([]EventTypeResult,
	error) {
	var results []EventTypeResult
	if err := betRequest("listEventTypes", s, q, &results, fn...); err != nil {
		return nil, err
	}
	return results, nil
}

// Returns country list as string or error if occured
func (s *Session) ListCountries(q *Query, fn ...VisitorFunc) ([]CountryCodeResult,
	error) {
	var results []CountryCodeResult
	if err := betRequest("listCountries", s, q, &results, fn...); err != nil {
		return nil, err
	}
	return results, nil
}

// Returns events list as string or error if occured
func (s *Session) ListEvents(q *Query, fn ...VisitorFunc) ([]EventResult, error) {
	var results []EventResult
	if err := betRequest("listEvents", s, q, &results, fn...); err != nil {
		return nil, err
	}
	return results, nil
}

// Returns competitions list (ie. world cop) as string or error if occured
func (s *Session) ListCompetitions(q *Query, fn ...VisitorFunc) (
	[]CompetitionResult, error) {
	var results []CompetitionResult
	if err := betRequest("listCompetitions", s, q, &results, fn...); err != nil {
		return nil, err
	}
	return results, nil
}

// Returns a list of market types (i.e. MATCH_ODDS, NEXT_GOAL)
func (s *Session) ListMarketTypes(q *Query, fn ...VisitorFunc) (
	[]MarketTypeResult, error) {
	var results []MarketTypeResult
	if err := betRequest("listMarketTypes", s, q, &results, fn...); err != nil {
		return nil, err
	}
	return results, nil
}

// Returns a list of Venues (i.e. Cheltenham, Ascot)
func (s *Session) ListVenues(q *Query, fn ...VisitorFunc) ([]VenueResult,
	error) {
	var results []VenueResult
	if err := betRequest("listVenues", s, q, &results, fn...); err != nil {
		return nil, err
	}
	return results, nil
}

// Returns a list of information about published (ACTIVE/SUSPENDED) markets
func (s *Session) ListMarketCatalogue(q *Query, fn ...VisitorFunc) (
	[]MarketCatalogue, error) {
	var results []MarketCatalogue
	if err := betRequest("listMarketCatalogue", s, q, &results, fn...); err != nil {
		return nil, err
	}
	return results, nil
}

// Returns a list of dynamic data about markets
func (s *Session) ListMarketBook(q *Query, fn ...VisitorFunc) ([]MarketBook,
	error) {
	var results []MarketBook
	if err := betRequest("listMarketBook", s, q, &results, fn...); err != nil {
		return nil, err
	}
	return results, nil
}

// Retrieve profit and loss for a given list of markets
func (s *Session) ListMarketProfitAndLoss(q *Query, fn ...VisitorFunc) (
	[]MarketProfitAndLoss, error) {
	var results []MarketProfitAndLoss
	if err := betRequest("listMarketProfitAndLoss", s, q, &results, fn...); err != nil {
		return nil, err
	}
	return results, nil
}

// performs betting api requests
func betRequest(method string, s *Session, q *Query, r interface{},
	fn ...VisitorFunc) error {
	if q == nil {
		s.logger.Fatal("query parameter can not be nil")
		return errors.New("query parameter can not be nil")
	}

	p, err := json.Marshal(q)
	if err != nil {
		s.logger.Fatal(err)
		return err
	}

	payload := strings.NewReader(string(p))
	resp, err := doRequest(s, "betting", method, payload)
	if err != nil {
		s.logger.Fatal(method, err)
		return err
	}
	s.logger.Print(string(resp))

	if err := json.Unmarshal(resp, r); err != nil {
		s.logger.Fatal(method, err)
		return err
	}

	for _, f := range fn {
		f(s, q, r)
	}

	return nil
}
