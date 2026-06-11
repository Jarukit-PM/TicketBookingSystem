package main

import "github.com/Jarukit-PM/TicketBookingSystem/api/internal/catalog"

// Thai cinema chains represented in seed data (Major Cineplex + SF Cinema).
// Movie posters sourced from cinematic.asia (Bangkok cinema listings, June 2026).

type thaiCinemaSeed struct {
	Chain   string
	Name    string
	Address string
	Screens []string
}

type thaiMovieSeed struct {
	Title       string
	PosterURL   string
	DurationMin int
	Rating      string
	Synopsis    string
	Status      string
}

var thaiCinemas = []thaiCinemaSeed{
	{
		Chain:   "Major",
		Name:    "Major Cineplex Paragon",
		Address: "5F Siam Paragon, 991 Rama I Rd, Pathum Wan, Bangkok 10330",
		Screens: []string{"Theater 1", "Theater 2", "Gold Class"},
	},
	{
		Chain:   "Major",
		Name:    "Major Cineplex CentralWorld",
		Address: "7F CentralWorld, 4/5 Ratchadamri Rd, Pathum Wan, Bangkok 10330",
		Screens: []string{"Theater 5", "Theater 6", "IMAX Laser"},
	},
	{
		Chain:   "Major",
		Name:    "Major Cineplex Ratchayothin",
		Address: "Major Avenue Ratchayothin, Phahonyothin Rd, Chatuchak, Bangkok 10900",
		Screens: []string{"Theater 1", "Theater 2"},
	},
	{
		Chain:   "Major",
		Name:    "Major Cineplex Gateway Bang Sue",
		Address: "5F Gateway at Bang Sue, 235 Pracharat Sai 2 Rd, Bang Sue, Bangkok 10800",
		Screens: []string{"Theater 3", "Theater 4"},
	},
	{
		Chain:   "SF",
		Name:    "SF World Cinema CentralWorld",
		Address: "7F CentralWorld, 999/9 Rama I Rd, Pathum Wan, Bangkok 10330",
		Screens: []string{"World Screen 1", "World Screen 2", "Atmos Deluxe"},
	},
	{
		Chain:   "SF",
		Name:    "SF Cinema Terminal 21 Asok",
		Address: "6F Terminal 21 Asok, 88 Sukhumvit Soi 19, Watthana, Bangkok 10110",
		Screens: []string{"Hall A", "Hall B"},
	},
	{
		Chain:   "SF",
		Name:    "SF Cinema MBK Center",
		Address: "7F MBK Center, 444 Phayathai Rd, Pathum Wan, Bangkok 10330",
		Screens: []string{"Hall 1", "Hall 2"},
	},
}

var thaiMovies = []thaiMovieSeed{
	{
		Title:       "Pee Nak 5",
		PosterURL:   "https://cinematic.asia/p_posters/pee-nak-5.jpg",
		DurationMin: 113,
		Rating:      "15",
		Synopsis:    "The latest installment of Thailand's beloved ghost-comedy franchise returns monks and mischief to a haunted dormitory.",
		Status:      catalog.MovieStatusNowShowing,
	},
	{
		Title:       "The Devil Wears Prada 2",
		PosterURL:   "https://cinematic.asia/p_posters/devil-wears-prada-2.jpg",
		DurationMin: 108,
		Rating:      "13",
		Synopsis:    "Miranda Priestly and her former assistant navigate a new era of fashion media and power in this long-awaited sequel.",
		Status:      catalog.MovieStatusNowShowing,
	},
	{
		Title:       "Detective Conan the Movie 29",
		PosterURL:   "https://cinematic.asia/p_posters/detective-conan-the-movie-29.jpg",
		DurationMin: 110,
		Rating:      "G",
		Synopsis:    "Conan Edogawa faces a deadly mystery tied to a legendary gem and a shadowy syndicate in the latest anime blockbuster.",
		Status:      catalog.MovieStatusNowShowing,
	},
	{
		Title:       "Disclosure Day",
		PosterURL:   "https://cinematic.asia/p_posters/disclosure-day.jpg",
		DurationMin: 133,
		Rating:      "13",
		Synopsis:    "When extraterrestrial contact is confirmed on live television, world leaders scramble as humanity faces an uncertain future.",
		Status:      catalog.MovieStatusNowShowing,
	},
	{
		Title:       "Haunted Universities 4",
		PosterURL:   "https://cinematic.asia/p_posters/haunted-universities-4.jpg",
		DurationMin: 98,
		Rating:      "15",
		Synopsis:    "Four chilling campus legends from Thai universities are woven into one anthology of supernatural terror.",
		Status:      catalog.MovieStatusNowShowing,
	},
	{
		Title:       "Kijsada Paradise",
		PosterURL:   "https://cinematic.asia/p_posters/kijsada-paradise.jpg",
		DurationMin: 105,
		Rating:      "13",
		Synopsis:    "A Thai family drama about ambition, loyalty, and the price of chasing an idyllic life in Bangkok.",
		Status:      catalog.MovieStatusNowShowing,
	},
	{
		Title:       "Masters of the Universe",
		PosterURL:   "https://cinematic.asia/p_posters/masters-of-the-universe.jpg",
		DurationMin: 110,
		Rating:      "13",
		Synopsis:    "He-Man and Skeletor clash on Eternia and Earth in a live-action revival of the classic sword-and-sorcery saga.",
		Status:      catalog.MovieStatusNowShowing,
	},
	{
		Title:       "Scary Movie",
		PosterURL:   "https://cinematic.asia/p_posters/scary-movie.jpg",
		DurationMin: 88,
		Rating:      "15",
		Synopsis:    "The cult horror spoof returns to theaters with a new round of pop-culture parodies and jump-scare gags.",
		Status:      catalog.MovieStatusNowShowing,
	},
	{
		Title:       "The Amazing Digital Circus the Last Act",
		PosterURL:   "https://cinematic.asia/p_posters/amazing-digital-circus-the-last-act.jpg",
		DurationMin: 95,
		Rating:      "13",
		Synopsis:    "Pomni and the troupe confront the truth of their digital prison in a feature-length finale to the viral web series.",
		Status:      catalog.MovieStatusNowShowing,
	},
	{
		Title:       "Michael",
		PosterURL:   "https://cinematic.asia/p_posters/michael.jpg",
		DurationMin: 115,
		Rating:      "13",
		Synopsis:    "A biographical drama tracing the rise, artistry, and legacy of the King of Pop.",
		Status:      catalog.MovieStatusNowShowing,
	},
	{
		Title:       "Until We Meet Again",
		PosterURL:   "https://cinematic.asia/p_posters/until-we-meet-again.jpg",
		DurationMin: 112,
		Rating:      "13",
		Synopsis:    "Two lovers separated by fate search for each other across lifetimes in this Thai romantic fantasy.",
		Status:      catalog.MovieStatusNowShowing,
	},
	{
		Title:       "Assassination Classroom the Movie Our Time",
		PosterURL:   "https://cinematic.asia/p_posters/assassination-classroom-the-movie-our-time.jpg",
		DurationMin: 102,
		Rating:      "G",
		Synopsis:    "Class 3-E reunites for a final mission as Koro-sensei’s students face their most emotional challenge yet.",
		Status:      catalog.MovieStatusNowShowing,
	},
	{
		Title:       "Power Ballad",
		PosterURL:   "https://cinematic.asia/p_posters/power-ballad.jpg",
		DurationMin: 118,
		Rating:      "13",
		Synopsis:    "A washed-up metal frontman and a pop diva are forced to co-headline a comeback tour.",
		Status:      catalog.MovieStatusComingSoon,
	},
	{
		Title:       "Colony",
		PosterURL:   "https://cinematic.asia/p_posters/colony.jpg",
		DurationMin: 104,
		Rating:      "15",
		Synopsis:    "Survivors on a remote island discover their refuge hides a terrifying experiment.",
		Status:      catalog.MovieStatusComingSoon,
	},
}

// thaiPriceTiers reflects typical Bangkok chain pricing in satang (1 THB = 100 satang).
var thaiPriceTiers = catalog.PriceTiers{
	Standard:   22000, // 220 THB
	VIP:        32000, // 320 THB Gold Class / premium
	Wheelchair: 22000,
}

// dailyShowtimeHours are local Bangkok screening slots.
var dailyShowtimeHours = []int{13, 16, 19, 22}

const thaiShowtimeDays = 30
