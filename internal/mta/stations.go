package mta

import (
	"fmt"
	"strings"
)

// stationMap maps station names to their MTA Stop IDs (Complex/Parent IDs).
// Some stations house multiple lines across different IDs in the same complex.
var stationMap = map[string][]string{
	// 1, 2, 3 (Broadway-7th Ave)
	"Van Cortlandt Park-242 St": {"101"},
	"238 St":                    {"103"},
	"231 St":                    {"104"},
	"Marble Hill-225 St":        {"106"},
	"215 St":                    {"107"},
	"207 St":                    {"108"},
	"Dyckman St":                {"109"},
	"191 St":                    {"110"},
	"181 St":                    {"111"},
	"168 St-Washington Hts":     {"112", "A09"},
	"157 St":                    {"113"},
	"145 St":                    {"114"},
	"137 St-City College":       {"115"},
	"125 St":                    {"116"},
	"116 St-Columbia University": {"117"},
	"Cathedral Pkwy (110 St)":   {"118"},
	"103 St":                    {"119"},
	"96 St":                     {"120"},
	"86 St":                     {"121"},
	"79 St":                     {"122"},
	"72 St":                     {"123"},
	"66 St-Lincoln Center":      {"124"},
	"59 St-Columbus Circle":     {"125", "A24"},
	"50 St":                     {"126", "A25"},
	"Times Sq-42 St":            {"127", "R16", "725"},
	"34 St-Penn Station":        {"128", "A28"},
	"28 St":                     {"129"},
	"23 St":                     {"130"},
	"18 St":                     {"131"},
	"14 St":                     {"132", "A31", "L03"},
	"Christopher St":            {"133"},
	"Houston St":                {"134"},
	"Canal St":                  {"135", "R23", "A34", "M20"},
	"Franklin St":               {"136"},
	"Chambers St":               {"137", "A36"},
	"WTC Cortlandt":             {"138"},
	"South Ferry":               {"142", "R27"},

	// 2, 3 (Lenox / Eastern Pkwy)
	"Franklin Av":               {"139"},
	"Central Park North (110 St)": {"227"},
	"125 St 23":                 {"224"},
	"Fulton St":                 {"229", "A38", "418", "M22"},
	"Wall St":                   {"230", "419"},
	"Clark St":                  {"231"},
	"Borough Hall":              {"232", "423"},
	"Nevins St":                 {"234"},
	"Atlantic Av-Barclays Ctr":  {"235", "R31", "D24"},
	"Grand Army Plaza":          {"237"},
	"Crown Hts-Utica Av":        {"250"},

	// 4, 5, 6 (Lexington Ave)
	"Woodlawn":                  {"401"},
	"Mosholu Pkwy":              {"402"},
	"Bedford Park Blvd-Lehman College": {"405"},
	"Kingsbridge Rd":            {"406"},
	"Fordham Rd":                {"407"},
	"183 St":                    {"408"},
	"Burnside Av":               {"409"},
	"176 St":                    {"410"},
	"Mt Eden Av":                {"411"},
	"170 St":                    {"412"},
	"167 St":                    {"413"},
	"161 St-Yankee Stadium":     {"414", "D11"},
	"149 St-Grand Concourse":    {"415", "222"},
	"138 St-Grand Concourse":    {"416"},
	"125 St 456":                {"621"},
	"116 St 6":                  {"622"},
	"110 St 6":                  {"623"},
	"103 St 6":                  {"624"},
	"96 St 6":                   {"625"},
	"86 St 6":                   {"626"},
	"77 St 6":                   {"627"},
	"68 St-Hunter College":      {"628"},
	"59 St 6":                   {"629", "R11"},
	"51 St":                     {"630"},
	"Grand Central-42 St":       {"631", "723"},
	"33 St 6":                   {"632"},
	"28 St 6":                   {"633"},
	"23 St 6":                   {"634"},
	"14 St-Union Sq":            {"635", "R20", "L03"},
	"Astor Pl":                  {"636"},
	"Bleecker St":               {"637", "D21"},
	"Spring St 6":               {"638"},
	"Canal St 6":                {"639"},
	"Brooklyn Bridge-City Hall": {"640"},

	// A, C, E (8th Ave)
	"Inwood-207 St":             {"A02"},
	"Dyckman St A":              {"A03"},
	"190 St":                    {"A05"},
	"181 St A":                  {"A06"},
	"175 St":                    {"A07"},
	"163 St-Amsterdam Av":       {"A10"},
	"155 St":                    {"A11"},
	"145 St A":                  {"A12"},
	"135 St A":                  {"A14"},
	"125 St A":                  {"A15"},
	"116 St A":                  {"A16"},
	"Cathedral Pkwy (110 St) A": {"A17"},
	"103 St A":                  {"A18"},
	"96 St A":                   {"A19"},
	"86 St A":                   {"A20"},
	"81 St-Museum of Natural History": {"A21"},
	"72 St A":                   {"A22"},
	"50 St E":                   {"A25"},
	"42 St-Port Authority Bus Terminal": {"A27"},
	"23 St A":                   {"A30"},
	"14 St A":                   {"A31"},
	"W 4 St-Wash Sq":            {"A32", "D20"},
	"High St":                   {"A40"},
	"Jay St-MetroTech":          {"A41", "R29"},
	"Hoyt-Schermerhorn Sts":     {"A42"},
	"Nostrand Av A":             {"A46"},
	"Utica Av A":                {"A48"},
	"Euclid Av":                 {"A55"},
	"Howard Beach-JFK Airport":  {"H03"},

	// B, D, F, M (6th Ave)
	"21 St-Queensbridge":        {"B04"},
	"Roosevelt Island":          {"B05"},
	"Lexington Av/63 St":        {"B08"},
	"57 St":                     {"B10"},
	"47-50 Sts-Rockefeller Ctr": {"D15"},
	"42 St-Bryant Pk":           {"D16", "724"},
	"34 St-Herald Sq":           {"D17", "R17"},
	"23 St F":                   {"D18"},
	"14 St F":                   {"D19"},
	"Broadway-Lafayette St":     {"D21", "636"},
	"Grand St":                  {"D22"},
	"Second Av":                 {"F14"},
	"Delancey St-Essex St":      {"F15", "M18"},
	"East Broadway":             {"F16"},
	"York St":                   {"F18"},
	"West 8 St-NY Aquarium":     {"F38", "D42"},
	"Coney Island-Stillwell Av": {"D43", "F39", "N12"},

	// N, Q, R, W (Broadway / 2nd Ave)
	"Astoria-Ditmars Blvd":      {"R01"},
	"Queensboro Plaza":          {"R09", "718"},
	"Lexington Av/59 St":        {"R11", "629"},
	"5 Av/59 St":                {"R13"},
	"57 St-7 Av":                {"R14"},
	"49 St":                     {"R15"},
	"28 St B":                   {"R18"},
	"23 St B":                   {"R19"},
	"8 St-NYU":                  {"R21"},
	"Prince St":                 {"R22"},
	"City Hall":                 {"R24"},
	"Cortlandt St":              {"R25"},
	"Rector St B":               {"R26"},
	"Whitehall St-South Ferry":  {"R27", "142"},
	"DeKalb Av":                 {"R30"},
	"96 St Q":                   {"Q05"},
	"86 St Q":                   {"Q04"},
	"72 St Q":                   {"Q03"},

	// L (Canarsie)
	"8 Av":                      {"L01", "A31"},
	"6 Av":                      {"L02", "D19"},
	"1 Av":                      {"L06"},
	"Bedford Av":                {"L08"},
	"Lorimer St":                {"L10", "G32"},
	"Graham Av":                 {"L11"},
	"Grand St L":                {"L12"},
	"Canarsie-Rockaway Pkwy":    {"L29"},

	// 7 (Flushing)
	"Flushing-Main St":          {"701"},
	"Mets-Willets Point":        {"702"},
	"Junction Blvd":             {"707"},
	"74 St-Broadway 7":          {"710", "R14", "G14"},
	"Court Sq 7":                {"719", "G22"},
	"Hunters Point Av":          {"720"},
	"Vernon Blvd-Jackson Av":    {"721"},
	"34 St-Hudson Yards":        {"726"},

	// G (Crosstown)
	"Court Sq G":                {"G22", "719"},
	"21 St G":                   {"G24"},
	"Greenpoint Av":             {"G26"},
	"Nassau Av":                 {"G28"},
	"Metropolitan Av":           {"G29", "M01"},
	"Broadway G":                {"G30"},

	// J, Z (Nassau)
	"Jamaica Center-Parsons/Archer": {"G05", "E09"},
	"Sutphin Blvd-Archer Av-JFK":    {"G06", "E06"},
	"Marcy Av":                      {"M16"},
	"Bowery":                        {"M19"},
	"Broad St":                      {"M23"},

	// S (Shuttles)
	"Franklin Av Shuttle":       {"S01"},
	"Botanic Garden":            {"S04", "239"},
	"Prospect Park S":           {"S06", "D26"},
}

// ResolveStation converts a station name to a slice of MTA Stop IDs.
func ResolveStation(name string) ([]string, error) {
	for k, v := range stationMap {
		if strings.EqualFold(k, name) {
			return v, nil
		}
	}
	return nil, fmt.Errorf("station not found: %s", name)
}

// GetStationNames returns a list of all station names for the UI.
func GetStationNames() []string {
	names := make([]string, 0, len(stationMap))
	for k := range stationMap {
		names = append(names, k)
	}
	return names
}
