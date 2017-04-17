package main

import "regexp"

type staticDir struct {
	path   string
	regstr string
	regex  *regexp.Regexp
}
