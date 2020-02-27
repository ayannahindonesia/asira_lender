package irate

import (
	"math"
)

func PMT(rate float64, nper float64, pv float64, fv float64) float64 {
	if rate == 0 {
		return -(pv + fv) / nper
	}

	pvif := math.Pow(1+rate, nper)

	return rate / (pvif - 1) * -(pv*pvif + fv)
}

func IPMT(pv float64, pmt float64, rate float64, per float64) float64 {
	tmp := math.Pow(1+rate, per)
	return 0 - (pv*tmp*rate + pmt*(tmp-1))
}

func PPMT(rate float64, per float64, nper float64, pv float64, fv float64) float64 {
	if per < 1 || per >= nper+1 {
		return 0
	}

	pmt := PMT(rate, nper, pv, fv)

	return pmt - IPMT(pv, pmt, rate, per-1)
}
