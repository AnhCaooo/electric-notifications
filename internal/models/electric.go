package models

// Represents a struct of data that will received from RabbitMQ producer.
type PricesMessage struct {
	Data      TodayTomorrowPrice `json:"data"`      // Data represents the price of today and tomorrow
	TimeStamp string             `json:"timestamp"` // TimeStamp represents the time when the message is produced. This will help `notification-service` to decide whether to push notifications or not.
}

// Represents a struct of today and tomorrow exchange price
type TodayTomorrowPrice struct {
	Today    DailyPrice `json:"today"`
	Tomorrow DailyPrice `json:"tomorrow"`
}

// Represents a struct of daily price and bool flag to indicate does tomorrow's price available or not
type DailyPrice struct {
	Available bool        `json:"available"`
	Prices    PriceSeries `json:"prices"`
}

// Represents a series of electric data with the name of unit (ex: c/kwh)
type PriceSeries struct {
	Name string `json:"name" example:"c/kWh"` // unit of electric price
	Data []Data `json:"data"`
}

// Represents single electric data at specific time
type Data struct {
	TimeUTC      string  `json:"time_utc" example:"2024-12-08 22:00:00"`  // timestamp in UTC format
	OriginalTime string  `json:"orig_time" example:"2024-12-09 00:00:00"` // the current time where server is located
	Time         string  `json:"time" example:"2024-12-09 00:00:00"`      // the current time.
	Price        float64 `json:"price" example:"2.47"`                    // the price of specified time range
	VatFactor    float64 `json:"vat_factor" example:"1.255"`              // amount of VAT that applies to electric price.
	IsToday      bool    `json:"isToday" example:"false"`                 // IsToday indicates whether the current time is today or not
	IncludeVat   string  `json:"includeVat" example:"1" enums:"0,1"`      // IncludeVat is legacy property that return string value and value "0" means no VAT included and string "1" is included
}
