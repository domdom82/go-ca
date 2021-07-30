package ca

import (
	"fmt"
	"math/rand"
)

type person struct {
	first string
	last  string
	email string
}

type location struct {
	country string
	city    string
	street  string
	number  int
	zip     int
}

type company struct {
	name    string
	domain  string
	address *location
}

var firstNames = []string{"Andrew", "Brian", "Chris", "Dylan", "Eric", "Frank", "Gordon", "Henry", "Ian", "John", "Kyle", "Larry", "Mark", "Nathan", "Oscar", "Paul",
	"Quentin", "Robert", "Steve", "Tyler", "Ulysses", "Vincent", "Walter", "Xavier", "Yusuf", "Zach",
	"Amelia", "Brittany", "Chloe", "Diana", "Elizabeth", "Florence", "Grace", "Haley", "Isabelle", "Jane", "Kelly", "Lindsey", "Madison", "Nicole", "Olivia", "Patricia",
	"Quinn", "Rachel", "Susan", "Taylor", "Uma", "Violet", "Whitney", "Xenia", "Yvette", "Zoey"}

var lastNames = []string{"Andersen", "Brown", "Cabbage", "Dawson", "Emerson", "Freeman", "Graham", "Hamilton", "Irving", "Jackson", "Keller", "Lewis", "Morgan", "Newton",
	"Owens", "Palmer", "Quigley", "Reynolds", "Smith", "Thompson", "Underwood", "Vaughn", "Williams", "Xenakis", "Yates", "Zimmerman"}

var domains = []string{"FootBow.com", "HabitWays.com", "TravelFlys.com", "HappyGangs.com", "GiftSavage.com", "MelodiesOfLife.com", "GameShowEvents.com", "FortuneHi.com",
	"MyFitnessSpot.com", "HotYogaNow.com", "CasaSupermarket.com", "OneTwoCare.com", "HolidayCheese.com", "BigTimeInternet.com", "OnetimeSales.com", "LifesGoodAgain.com",
	"MyTravelScore.com", "TheGoodFoodShop.com", "24Gaga.com", "9Lama.com", "GoatColor.com", "GoSoftTech.com", "TeamBlackwell.com", "CashIncomeOnline.com", "StandGood.com",
	"GymGyms.com", "RentToCars.com", "DesignByLogic.com", "VillagersWorld.com", "BigRedTag.com"}

var cities = []string{"Farmingside", "Bridgetown", "Baypool", "Aelborough", "Springville", "Great Foxville", "Casterfolk Hills", "Southcester", "Clamhampton", "Wingport",
	"Proford", "Angerville", "Bayford", "Castertown", "Jamestown", "Holtsness", "Backkarta", "Fauxtown", "Transness", "Medburgh", "South Hogfolk", "Fortburg", "Dayworth",
	"Redland", "Pailtown", "Hamview", "Highborough", "Cruxborough", "Lawborough", "Cloudtown", "West Bridgeley", "Sweetview", "Riverwich", "Fortcaster", "Parkgrad", "Bayburg",
	"Angerview", "Winterport", "Southingmouth", "Westdol", "Redport", "Costston"}

var streets = []string{"Fairview Road", "Shady Lane", "Madison Street", "Front Street", "Hilltop Road", "Locust Street", "Cooper Street", "Summit Street", "Mulberry Street",
	"Berkshire Drive", "Orchard Street", "Oak Avenue", "Grant Avenue", "Clark Street", "Taylor Street", "5th Street East", "Crescent Street", "Market Street", "State Street",
	"Adams Avenue", "Eagle Road", "5th Street", "Jones Street", "Hamilton Road", "Brook Lane", "Virginia Avenue", "York Road", "Valley Drive", "Grove Avenue", "Buckingham Drive",
	"Mill Street", "Mulberry Court", "Howard Street", "Ridge Avenue", "3rd Street", "Route 9", "Race Street", "Aspen Court", "Canterbury Road", "Brown Street", "Valley View Drive",
	"Euclid Avenue", "Marshall Street", "Oxford Road", "Pin Oak Drive", "Franklin Court", "Franklin Street", "Grand Avenue", "Aspen Drive", "Oak Street"}

var hostnames = []string{"www", "authentication", "account", "login", "api", "internal", "admin", "test", "dev", "prod", "alpha", "beta"}

func random(a []string) string {
	return a[rand.Intn(len(a))]
}

func randomLocation() *location {
	l := &location{
		country: "US",
		city:    random(cities),
		street:  random(streets),
		number:  rand.Intn(100) + 1,
		zip:     rand.Intn(89999) + 10000,
	}
	return l
}

func randomCompany() *company {
	c := &company{
		domain:  random(domains),
		address: randomLocation(),
	}

	c.name = c.domain

	return c
}
func randomPerson(comp *company) *person {
	p := &person{
		first: random(firstNames),
		last:  random(lastNames),
	}
	p.email = fmt.Sprintf("%s%s@%s", p.first[:1], p.last, comp.domain)

	return p
}

func randomHostname() string {
	return random(hostnames)
}
