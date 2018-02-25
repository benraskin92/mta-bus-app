package config

type Config struct {
	Email EmailUser `yaml:"email"`
	MTA   MTAInfo   `yaml:"mta"`
}

type EmailUser struct {
	Username string   `yaml:"username"`
	Password string   `yaml:"password"`
	Server   string   `yaml:"server"`
	Port     int      `yaml:"port"`
	SendTo   []string `yaml:"sendTo"`
}

type MTAInfo struct {
	Key       string `yaml:"api_key"`
	Line      string `yaml:"line"`
	Direction string `yaml:"direction"`
	StopCheck string `yaml:"stop_check"`
	BeginTime int    `yaml:"begin_time"`
	EndTime   int    `yaml:"end_time"`
	Weekends  bool   `yaml:"weekends"`
}

type MTAResponse struct {
	Siri struct {
		ServiceDelivery struct {
			ResponseTimestamp         string `json:"ResponseTimestamp"`
			VehicleMonitoringDelivery []struct {
				VehicleActivity []struct {
					MonitoredVehicleJourney struct {
						LineRef                 string `json:"LineRef"`
						DirectionRef            string `json:"DirectionRef"`
						FramedVehicleJourneyRef struct {
							DataFrameRef           string `json:"DataFrameRef"`
							DatedVehicleJourneyRef string `json:"DatedVehicleJourneyRef"`
						} `json:"FramedVehicleJourneyRef"`
						JourneyPatternRef string        `json:"JourneyPatternRef"`
						PublishedLineName string        `json:"PublishedLineName"`
						OperatorRef       string        `json:"OperatorRef"`
						OriginRef         string        `json:"OriginRef"`
						DestinationRef    string        `json:"DestinationRef"`
						DestinationName   string        `json:"DestinationName"`
						SituationRef      []interface{} `json:"SituationRef"`
						Monitored         bool          `json:"Monitored"`
						VehicleLocation   struct {
							Longitude float64 `json:"Longitude"`
							Latitude  float64 `json:"Latitude"`
						} `json:"VehicleLocation"`
						Bearing       float64 `json:"Bearing"`
						ProgressRate  string  `json:"ProgressRate"`
						BlockRef      string  `json:"BlockRef"`
						VehicleRef    string  `json:"VehicleRef"`
						MonitoredCall struct {
							ExpectedArrivalTime   string `json:"ExpectedArrivalTime"`
							ExpectedDepartureTime string `json:"ExpectedDepartureTime"`
							Extensions            struct {
								Distances struct {
									StopsFromCall          int     `json:"StopsFromCall"`
									PresentableDistance    string  `json:"PresentableDistance"`
									DistanceFromCall       float64 `json:"DistanceFromCall"`
									CallDistanceAlongRoute float64 `json:"CallDistanceAlongRoute"`
								} `json:"Distances"`
							} `json:"Extensions"`
							StopPointRef  string `json:"StopPointRef"`
							VisitNumber   int    `json:"VisitNumber"`
							StopPointName string `json:"StopPointName"`
						} `json:"MonitoredCall"`
						OnwardCalls struct {
						} `json:"OnwardCalls"`
					} `json:"MonitoredVehicleJourney"`
					RecordedAtTime string `json:"RecordedAtTime"`
				} `json:"VehicleActivity"`
				ResponseTimestamp string `json:"ResponseTimestamp"`
				ValidUntil        string `json:"ValidUntil"`
			} `json:"VehicleMonitoringDelivery"`
			SituationExchangeDelivery []interface{} `json:"SituationExchangeDelivery"`
		} `json:"ServiceDelivery"`
	} `json:"Siri"`
}
