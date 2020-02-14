package utils

import (
	"../../../entity/pharmacy"
	"../../../storage/maskStorage"
	"math"
	"sort"
)

const (
	FilterAdult = "大人庫存大於0"
	FilterChild = "兒童庫存大於0"
	FilterAdultAndChild = "兩者皆有庫存"
	FilterZero = "兩者皆無庫存"
)

// 算法1
func EarthDistance(lat1, lng1, lat2, lng2 float64) float64 {
	//radius := 6378.137 // km
	//radius := 6378137.0
	radius := 6367000.0
	rad := math.Pi / 180.0
	lat1 = lat1 * rad
	lng1 = lng1 * rad
	lat2 = lat2 * rad
	lng2 = lng2 * rad
	theta := lng2 - lng1
	dist := math.Acos(math.Sin(lat1)*math.Sin(lat2) + math.Cos(lat1)*math.Cos(lat2)*math.Cos(theta))
	return dist * radius
}

// 算法2
func GetDistance(lat1, lng1, lat2, lng2 float64) float64 {

	earthRadius := 6367000.0 //approximate radius of earth in meters
	//  Convert these degrees to radians       to work with the formula
	lat1 = (lat1 * math.Pi) / 180
	lng1 = (lng1 * math.Pi) / 180
	lat2 = (lat2 * math.Pi) / 180
	lng2 = (lng2 * math.Pi) / 180
	// Using the       Haversine formula       http://en.wikipedia.org/wiki/Haversine_formula       calculate the distance
	calcLongitude := lng2 - lng1
	calcLatitude := lat2 - lat1
	stepOne := math.Pow(math.Sin(calcLatitude/2), 2) + math.Cos(lat1)*math.Cos(lat2)*math.Pow(math.Sin(calcLongitude/2), 2)
	stepTwo := 2 * math.Asin(math.Min(1, math.Sqrt(stepOne)))
	calculatedDistance := earthRadius * stepTwo
	return calculatedDistance
}

type PharmacyDistance struct {
	Distance float64
	MaskData *pharmacy.Pharmacy
}

type PharmacyDistanceArray []PharmacyDistance

//Len()
func (da PharmacyDistanceArray) Len() int {
	return len(da)
}

//Less() 由小到大排序
func (da PharmacyDistanceArray) Less(i, j int) bool {
	if da[i].Distance < da[j].Distance {
		return true
	} else {
		return false
	}
}

//Swap()
func (da PharmacyDistanceArray) Swap(i, j int) {
	da[i], da[j] = da[j], da[i]
}

func CalcLocation(storage *maskStorage.Storage, topCnt int, filterData string, latitude, longitude float64) []PharmacyDistance {
	maskDataList := storage.GetAllList()
	newDataArray := PharmacyDistanceArray{}
	ListLoop:
	for idx := range maskDataList {
		maskData := maskDataList[idx]
		if maskData.UpdTime == nil || maskData.Disabled != 0 {
			continue ListLoop
		} else {
			switch filterData {
			case FilterAdult:
				if maskData.AdultCount <= 0 {
					continue ListLoop
				}
			case FilterChild:
				if maskData.ChildCount <= 0 {
					continue ListLoop
				}
			case FilterAdultAndChild:
				if maskData.AdultCount <= 0 || maskData.ChildCount <= 0 {
					continue ListLoop
				}
			case FilterZero:
				if maskData.AdultCount > 0 || maskData.ChildCount > 0 {
					continue ListLoop
				}
			default:
				if maskData.AdultCount <= 0 && maskData.ChildCount <= 0 {
					continue ListLoop
				}
			}
		}
		pharmacyDistance := PharmacyDistance{
			Distance: EarthDistance(latitude, longitude, maskData.Latitude, maskData.Longitude),
			MaskData: &maskData,
		}
		newDataArray = append(newDataArray, pharmacyDistance)
	}
	sort.Sort(newDataArray)
	if len(newDataArray) < topCnt {
		topCnt = len(newDataArray)
	}
	return newDataArray[0:topCnt]
}
