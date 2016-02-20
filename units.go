package main

type unit struct {
	unitID     int
	unitHealth int
	typ        *unitType
}

type unitType struct {
	//basic stats go here
	maxHealth int
	movement  float64
}
