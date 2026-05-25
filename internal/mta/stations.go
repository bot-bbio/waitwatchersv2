package mta

import (
	"fmt"
	"strings"
)

// stationMap maps station names to their MTA Stop IDs (Complex/Parent IDs).
// Some stations house multiple lines across different IDs in the same complex.
var stationMap = map[string][]string{
	// 1, 2, 3 (Broadway-7th Ave)
	"Van Cortlandt Park-242 St [1]": {"101"},
	"238 St [1]":                    {"103"},
	"231 St [1]":                    {"104"},
	"Marble Hill-225 St [1]":        {"106"},
	"215 St [1]":                    {"107"},
	"207 St [1]":                    {"108"},
	"Dyckman St [1]":                {"109"},
	"191 St [1]":                    {"110"},
	"181 St [1]":                    {"111"},
	"168 St-Washington Hts [1,A,C]": {"112", "A09"},
	"157 St [1]":                    {"113"},
	"145 St [1]":                    {"114"},
	"137 St-City College [1]":       {"115"},
	"125 St [1]":                    {"116"},
	"116 St-Columbia University [1]": {"117"},
	"Cathedral Pkwy (110 St) [1]":   {"118"},
	"103 St [1]":                    {"119"},
	"96 St [1,2,3]":                 {"120"},
	"86 St [1]":                     {"121"},
	"79 St [1]":                     {"122"},
	"72 St [1,2,3]":                 {"123"},
	"66 St-Lincoln Center [1]":      {"124"},
	"59 St-Columbus Circle [1,A,B,C,D]": {"125", "A24"},
	"50 St [1,C,E]":                 {"126", "A25"},
	"42 St-Times Sq / Port Authority [1,2,3,7,A,C,E,N,Q,R,W,S]": {"127", "R16", "725", "A27"},
	"34 St-Penn Station [1,2,3,A,C,E]": {"128", "A28"},
	"28 St [1]":                     {"129"},
	"23 St [1]":                     {"130"},
	"18 St [1]":                     {"131"},
	"14 St / 8 Av [1,2,3,A,C,E,L]":   {"132", "A31", "L01"},
	"14 St / 6 Av [F,M,L]":          {"L02", "D19"},
	"Christopher St [1]":            {"133"},
	"Houston St [1]":                {"134"},
	"Canal St [6,J,Z,N,Q,R,W,A,C,E]": {"135", "R23", "A34", "M20", "639"},
	"Franklin St [1]":               {"136"},
	"Chambers St [1,2,3,A,C,E]":     {"137", "A36"},
	"WTC Cortlandt [1]":             {"138"},
	"South Ferry / Whitehall St [1,N,R,W]": {"142", "R27"},

	// 2, 3 (Lenox / Eastern Pkwy)
	"Franklin Av / Botanic Garden [S,2,3,4,5]": {"139", "S04", "239"},
	"Central Park North (110 St) [2,3]":        {"227"},
	"125 St 23 [2,3]":                          {"224"},
	"Fulton St [2,3,4,5,A,C,J,Z]":              {"229", "A38", "418", "M22"},
	"Wall St [2,3,4,5]":                        {"230", "419"},
	"Clark St [2,3]":                           {"231"},
	"Borough Hall [2,3,4,5]":                   {"232", "423"},
	"Nevins St [2,3,4,5]":                      {"234"},
	"Atlantic Av-Barclays Ctr [2,3,4,5,B,D,N,Q,R]": {"235", "R31", "D24"},
	"Grand Army Plaza [2,3,4,5]":               {"237"},
	"Crown Hts-Utica Av [3,4,5]":               {"250"},

	// 4, 5, 6 (Lexington Ave)
	"Woodlawn [4]":                  {"401"},
	"Mosholu Pkwy [4]":              {"402"},
	"Bedford Park Blvd-Lehman College [4]": {"405"},
	"Kingsbridge Rd [4]":            {"406"},
	"Fordham Rd [4]":                {"407"},
	"183 St [4]":                    {"408"},
	"Burnside Av [4]":               {"409"},
	"176 St [4]":                    {"410"},
	"Mt Eden Av [4]":                {"411"},
	"170 St [4]":                    {"412"},
	"167 St [4]":                    {"413"},
	"161 St-Yankee Stadium [4,D]":   {"414", "D11"},
	"149 St-Grand Concourse [2,4,5]": {"415", "222"},
	"138 St-Grand Concourse [4,5]":  {"416"},
	"125 St 456 [4,5,6]":            {"621"},
	"116 St 6 [6]":                  {"622"},
	"110 St 6 [6]":                  {"623"},
	"103 St 6 [6]":                  {"624"},
	"96 St 6 [6]":                   {"625"},
	"86 St 6 [4,5,6]":               {"626"},
	"77 St [6]":                     {"627"},
	"68 St-Hunter College [6]":      {"628"},
	"59 St / Lexington Av [4,5,6,N,R,W]": {"629", "R11"},
	"51 St / Lexington Av-53 St [6,E,M]": {"630"},
	"Grand Central-42 St [4,5,6,7,S]":   {"631", "723"},
	"33 St [6]":                     {"632"},
	"28 St 6 [6]":                   {"633"},
	"23 St 6 [6]":                   {"634"},
	"14 St-Union Sq [4,5,6,L,N,Q,R,W]": {"635", "R20", "L03"},
	"Astor Pl [6]":                  {"636"},
	"Bleecker St / Broadway-Lafayette [6,B,D,F,M]": {"637", "D21", "636"},
	"Spring St 6 [6]":               {"638"},
	"Brooklyn Bridge-City Hall / Chambers St [4,5,6,J,Z]": {"640"},

	// A, C, E (8th Ave)
	"Inwood-207 St [A]":             {"A02"},
	"Dyckman St A [A]":              {"A03"},
	"190 St [A]":                    {"A05"},
	"181 St A [A]":                  {"A06"},
	"175 St [A]":                  {"A07"},
	"163 St-Amsterdam Av [C]":       {"A10"},
	"155 St [C]":                    {"A11"},
	"145 St A [A,B,C,D]":            {"A12"},
	"135 St A [B,C]":                {"A14"},
	"125 St A [A,B,C,D]":            {"A15"},
	"116 St A [B,C]":                {"A16"},
	"Cathedral Pkwy (110 St) A [B,C]": {"A17"},
	"103 St A [B,C]":                {"A18"},
	"96 St A [B,C]":                 {"A19"},
	"86 St A [B,C]":                 {"A20"},
	"81 St-Museum of Natural History [B,C]": {"A21"},
	"72 St A [B,C]":                 {"A22"},
	"23 St A [C,E]":                 {"A30"},
	"W 4 St-Wash Sq [A,C,E,B,D,F,M]": {"A32", "D20"},
	"High St [A,C]":                 {"A40"},
	"Jay St-MetroTech [A,C,F,R]":    {"A41", "R29"},
	"Hoyt-Schermerhorn Sts [A,C,G]": {"A42"},
	"Nostrand Av [A,C]":             {"A46"},
	"Utica Av [A,C]":                {"A48"},
	"Euclid Av [A,C]":               {"A55"},
	"Howard Beach-JFK Airport [A]":  {"H03"},

	// B, D, F, M (6th Ave)
	"21 St-Queensbridge [F]":        {"B04"},
	"Roosevelt Island [F]":          {"B05"},
	"Lexington Av/63 St [F,Q]":        {"B08"},
	"57 St [F]":                     {"B10"},
	"47-50 Sts-Rockefeller Ctr [B,D,F,M]": {"D15"},
	"42 St-Bryant Pk / 5 Av [B,D,F,M,7]":   {"D16", "724"},
	"34 St-Herald Sq [B,D,F,M,N,Q,R,W]":   {"D17", "R17"},
	"23 St F [F,M]":                 {"D18"},
	"14 St F [F,M]":                 {"D19"},
	"Grand St [B,D]":                  {"D22"},
	"Second Av [F]":                 {"F14"},
	"Delancey St-Essex St [F,J,M,Z]":      {"F15", "M18"},
	"East Broadway [F]":             {"F16"},
	"York St [F]":                   {"F18"},
	"West 8 St-NY Aquarium [F,Q]":     {"F38", "D42"},
	"Coney Island-Stillwell Av [D,F,N,Q]": {"D43", "F39", "N12"},

	// N, Q, R, W (Broadway / 2nd Ave)
	"Astoria-Ditmars Blvd [N,W]":      {"R01"},
	"Queensboro Plaza [7,N,W]":          {"R09", "718"},
	"5 Av/59 St [N,R,W]":                {"R13"},
	"57 St-7 Av [N,Q,R,W]":              {"R14"},
	"49 St [N,R,W]":                     {"R15"},
	"28 St B [R,W]":                   {"R18"},
	"23 St B [R,W]":                   {"R19"},
	"8 St-NYU [R,W]":                  {"R21"},
	"Prince St [R,W]":                 {"R22"},
	"City Hall [R,W]":                 {"R24"},
	"Cortlandt St [R,W]":              {"R25"},
	"Rector St B [R,W]":               {"R26"},
	"DeKalb Av [B,Q,R]":                 {"R30"},
	"96 St Q [Q]":                 {"Q05"},
	"86 St Q [Q]":                 {"Q04"},
	"72 St Q [Q]":                 {"Q03"},

	// L (Canarsie)
	"1 Av [L]":                      {"L06"},
	"Bedford Av [L]":                {"L08"},
	"Lorimer St / Metropolitan Av [L,G]": {"L10", "G32", "G29", "M01"},
	"Graham Av [L]":                {"L11"},
	"Grand St [L]":                {"L12"},
	"Canarsie-Rockaway Pkwy [L]":    {"L29"},

	// 7 (Flushing)
	"Flushing-Main St [7]":          {"701"},
	"Mets-Willets Point [7]":        {"702"},
	"Junction Blvd [7]":             {"707"},
	"Jackson Hts-Roosevelt Av / 74 St [7,E,F,M,R]": {"710", "R14", "G14"},
	"Court Sq [7,E,G,M]":                           {"719", "G22"},
	"Hunters Point Av [7]":          {"720"},
	"Presidents St [7]":             {"721"}, // renaming to avoid duplicate Whitehall
	"34 St-Hudson Yards [7]":        {"726"},

	// G (Crosstown)
	"21 St [G]":                   {"G24"},
	"Greenpoint Av [G]":             {"G26"},
	"Nassau Av [G]":                 {"G28"},
	"Broadway [G]":                {"G30"},

	// J, Z (Nassau)
	"Jamaica Center-Parsons/Archer [E,J,Z]": {"G05", "E09"},
	"Sutphin Blvd-Archer Av-JFK [E,J,Z]":    {"G06", "E06"},
	"Marcy Av [J,M,Z]":                      {"M16"},
	"Bowery [J,Z]":                        {"M19"},
	"Broad St [J,Z]":                      {"M23"},

	// S (Shuttles)
	"Prospect Park [B,Q,S]":           {"S06", "D26"},
}

// ResolveStation converts a station name to a slice of MTA Stop IDs.
func ResolveStation(name string) ([]string, error) {
	// First check: try exact match (case-insensitive)
	for k, v := range stationMap {
		if strings.EqualFold(k, name) {
			return v, nil
		}
	}

	// Second check: strip bracketed line indicators from input and keys
	inputBase := name
	if idx := strings.Index(name, "["); idx >= 0 {
		inputBase = name[:idx]
	}
	inputBase = strings.TrimSpace(inputBase)

	for k, v := range stationMap {
		keyBase := k
		if idx := strings.Index(k, "["); idx >= 0 {
			keyBase = k[:idx]
		}
		keyBase = strings.TrimSpace(keyBase)

		if strings.EqualFold(keyBase, inputBase) {
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
