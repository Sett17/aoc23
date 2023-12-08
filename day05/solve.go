package day05

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Seed int

func NewSeed(seed string) Seed {
	res, err := strconv.Atoi(seed)
	if err != nil {
		panic(err)
	}
	return Seed(res)
}

type SeedRange struct {
	start, length int
}

func NewSeedRange(start, length int) *SeedRange {
	return &SeedRange{
		start:  start,
		length: length,
	}
}

func (s *SeedRange) iterate() chan Seed {
	ch := make(chan Seed)
	go func() {
		for i := 0; i < s.length; i++ {
			ch <- Seed(s.start + i)
		}
		close(ch)
	}()
	return ch
}

func (s *SeedRange) iterateChunk(chunk_size int) chan Seed {
	ch := make(chan Seed)
	go func() {
		for i := 0; i < s.length; i += chunk_size {
			length := chunk_size
			if i+chunk_size > s.length {
				length = s.length - i
			}
			for j := 0; j < length; j++ {
				ch <- Seed(s.start + i + j)
			}
		}
		close(ch)
	}()
	return ch
}

func (s *SeedRange) IterateChunked(chunk_size int) []chan Seed {
	chunks := []chan Seed{}
	for i := 0; i < s.length; i += chunk_size {
		length := chunk_size
		if i+chunk_size > s.length {
			length = s.length - i
		}
		chunks = append(chunks, s.iterateChunk(length))
	}
	return chunks
}

type Fertilizer int
type Water int
type Light int
type Temperature int
type Humidity int
type Location int

type Range struct {
	src_start, dest_start, length int
	src_end, dest_end             int
}

func NewRange(src_start, dest_start, length int) *Range {
	return &Range{
		src_start:  src_start,
		dest_start: dest_start,
		length:     length,
		src_end:    src_start + length,
		dest_end:   dest_start + length,
	}
}

func (r *Range) IsInside(src int) bool {
	return src >= r.src_start && src < r.src_end
}

func (r *Range) Map(src int) int {
	return r.dest_start + (src - r.src_start)
}

type Interval struct {
	SrcStart  int
	DestStart int
	Length    int
}

type Mapping struct {
	intervals []Interval
}

func NewMapping() *Mapping {
	return &Mapping{intervals: []Interval{}}
}

func (m *Mapping) AddRangeMapping(srcStart, destStart, length int) {
	m.intervals = append(m.intervals, Interval{SrcStart: srcStart, DestStart: destStart, Length: length})
}

func (m *Mapping) Map(src int) int {
	for _, interval := range m.intervals {
		if src >= interval.SrcStart && src < interval.SrcStart+interval.Length {
			return interval.DestStart + (src - interval.SrcStart)
		}
	}
	return src
}

const SEED2SOIL = "seed-to-soil"
const SOIL2FERTILIZER = "soil-to-fertilizer"
const FERTILIZER2WATER = "fertilizer-to-water"
const WATER2LIGHT = "water-to-light"
const LIGHT2TEMPERATURE = "light-to-temperature"
const TEMPERATURE2HUMIDITY = "temperature-to-humidity"
const HUMIDITY2LOCATION = "humidity-to-location"

type Almanac struct {
	seed2soil            Mapping
	soil2fertilizer      Mapping
	fertilizer2water     Mapping
	water2light          Mapping
	light2temperature    Mapping
	temperature2humidity Mapping
	humidity2location    Mapping
}

func NewAlmanac() *Almanac {
	return &Almanac{
		seed2soil:            Mapping{intervals: []Interval{}},
		soil2fertilizer:      Mapping{intervals: []Interval{}},
		fertilizer2water:     Mapping{intervals: []Interval{}},
		water2light:          Mapping{intervals: []Interval{}},
		light2temperature:    Mapping{intervals: []Interval{}},
		temperature2humidity: Mapping{intervals: []Interval{}},
		humidity2location:    Mapping{intervals: []Interval{}},
	}
}

func (a *Almanac) AddSeed2SoilRange(seed_start, soil_start, length int) {
	fmt.Printf("adding seed2soil range: %v, %v, %v\n", seed_start, soil_start, length)
	a.seed2soil.AddRangeMapping(seed_start, soil_start, length)
}

func (a *Almanac) AddSoil2FertilizerRange(soil_start, fertilizer_start, length int) {
	fmt.Printf("adding soil2fertilizer range: %v, %v, %v\n", soil_start, fertilizer_start, length)
	a.soil2fertilizer.AddRangeMapping(soil_start, fertilizer_start, length)
}

func (a *Almanac) AddFertilizer2WaterRange(fertilizer_start, water_start, length int) {
	fmt.Printf("adding fertilizer2water range: %v, %v, %v\n", fertilizer_start, water_start, length)
	a.fertilizer2water.AddRangeMapping(fertilizer_start, water_start, length)
}

func (a *Almanac) AddWater2LightRange(water_start, light_start, length int) {
	fmt.Printf("adding water2light range: %v, %v, %v\n", water_start, light_start, length)
	a.water2light.AddRangeMapping(water_start, light_start, length)
}

func (a *Almanac) AddLight2TemperatureRange(light_start, temperature_start, length int) {
	fmt.Printf("adding light2temperature range: %v, %v, %v\n", light_start, temperature_start, length)
	a.light2temperature.AddRangeMapping(light_start, temperature_start, length)
}

func (a *Almanac) AddTemperature2HumidityRange(temperature_start, humidity_start, length int) {
	fmt.Printf("adding temperature2humidity range: %v, %v, %v\n", temperature_start, humidity_start, length)
	a.temperature2humidity.AddRangeMapping(temperature_start, humidity_start, length)
}

func (a *Almanac) AddHumidity2LocationRange(humidity_start, location_start, length int) {
	fmt.Printf("adding humidity2location range: %v, %v, %v\n", humidity_start, location_start, length)
	a.humidity2location.AddRangeMapping(humidity_start, location_start, length)
}

func Solve() (string, error) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	file, err := os.Open("input")
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	seedRanges := []*SeedRange{}
	almanac := NewAlmanac()

	scanner := bufio.NewScanner(file)
	fmt.Printf("reading data\n")
	firstLine := true
	activeMapping := ""
	for scanner.Scan() {
		line := scanner.Text()
		if firstLine {
			firstLine = false
			split := strings.Split(strings.TrimPrefix(line, "seeds: "), " ")
			for i := 0; i < len(split); i += 2 {
				seed_start, err := strconv.Atoi(split[i])
				if err != nil {
					return "", err
				}
				length, err := strconv.Atoi(split[i+1])
				if err != nil {
					return "", err
				}
				seedRanges = append(seedRanges, NewSeedRange(seed_start, length))
			}
			fmt.Printf("got %d seed ranges\n", len(seedRanges))
			continue
		}
		if len(strings.TrimSpace(line)) == 0 {
			continue
		}
		if strings.Index(line, ":") != -1 {
			if strings.HasPrefix(line, SEED2SOIL) {
				activeMapping = SEED2SOIL
			} else if strings.HasPrefix(line, SOIL2FERTILIZER) {
				activeMapping = SOIL2FERTILIZER
			} else if strings.HasPrefix(line, FERTILIZER2WATER) {
				activeMapping = FERTILIZER2WATER
			} else if strings.HasPrefix(line, WATER2LIGHT) {
				activeMapping = WATER2LIGHT
			} else if strings.HasPrefix(line, LIGHT2TEMPERATURE) {
				activeMapping = LIGHT2TEMPERATURE
			} else if strings.HasPrefix(line, TEMPERATURE2HUMIDITY) {
				activeMapping = TEMPERATURE2HUMIDITY
			} else if strings.HasPrefix(line, HUMIDITY2LOCATION) {
				activeMapping = HUMIDITY2LOCATION
			} else {
				return "", fmt.Errorf("unknown mapping: %v", line)
			}
			fmt.Printf("new mapping: %v\n", activeMapping)
			continue
		}
		numberStrings := strings.Split(line, " ")
		numbers := []int{}
		for _, numberString := range numberStrings {
			number, err := strconv.Atoi(numberString)
			if err != nil {
				return "", err
			}
			numbers = append(numbers, number)
		}
		switch activeMapping {
		case SEED2SOIL:
			almanac.AddSeed2SoilRange(numbers[1], numbers[0], numbers[2])
		case SOIL2FERTILIZER:
			almanac.AddSoil2FertilizerRange(numbers[1], numbers[0], numbers[2])
		case FERTILIZER2WATER:
			almanac.AddFertilizer2WaterRange(numbers[1], numbers[0], numbers[2])
		case WATER2LIGHT:
			almanac.AddWater2LightRange(numbers[1], numbers[0], numbers[2])
		case LIGHT2TEMPERATURE:
			almanac.AddLight2TemperatureRange(numbers[1], numbers[0], numbers[2])
		case TEMPERATURE2HUMIDITY:
			almanac.AddTemperature2HumidityRange(numbers[1], numbers[0], numbers[2])
		case HUMIDITY2LOCATION:
			almanac.AddHumidity2LocationRange(numbers[1], numbers[0], numbers[2])
		default:
			return "", fmt.Errorf("unknown mapping: %v", activeMapping)
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}

	fmt.Printf("got the data\n")

	locations := []Location{}
	seedTotal := 0

	var wg sync.WaitGroup
	out := make(chan Location, 1e3)
	progress := make(chan int, 1e3)
	progressCount := 0

	for _, seedRange := range seedRanges {
		seedTotal += seedRange.length
	}

	fmt.Printf("seedTotal: %v\n", seedTotal)
	fmt.Printf("\n\033[1A\033[1G")

	timerStart := time.Now()
	for i, seedRange := range seedRanges {
		channels := seedRange.IterateChunked(1e9)
		fmt.Printf("%d von %d ranges, channels: %d\n", i+1, len(seedRanges), len(channels))
		for _, channel := range channels {
			wg.Add(1)
			go func(channel chan Seed) {
				defer wg.Done()
				for seed := range channel {
					soil := almanac.seed2soil.Map(int(seed))
					fertilizer := almanac.soil2fertilizer.Map(soil)
					water := almanac.fertilizer2water.Map(fertilizer)
					light := almanac.water2light.Map(water)
					temperature := almanac.light2temperature.Map(light)
					humidity := almanac.temperature2humidity.Map(temperature)
					location := almanac.humidity2location.Map(humidity)
					out <- Location(location)
					progress <- 1
				}
			}(channel)
		}
	}
	go func(progress chan int) {
		for p := range progress {
			progressCount += p
			percentage_done := float64(progressCount) / float64(seedTotal) * 100
			elapsed := time.Since(timerStart)
			eta := elapsed.Seconds() / percentage_done * (100 - percentage_done)

			fmt.Printf("\033[1A\033[1M\033[1G")
			fmt.Printf("seed: %v/%v (%.4f%%)\teta: %v", progressCount, seedTotal, percentage_done, time.Duration(eta)*time.Second)
			fmt.Printf("\n")
		}
	}(progress)

	go func() {
		wg.Wait()
		close(out)
		close(progress)
	}()

	for location := range out {
		locations = append(locations, location)
	}

	// for _, seedRange := range seedRanges {
	// 	for seed := range seedRange.iterate() {
	// 		//clear line
	// 		fmt.Printf("\033[1M\033[1G")
	// 		percentage_done := float64(seedIndex) / float64(seedTotal) * 100
	// 		elapsed := time.Since(timerStart)
	// 		eta := elapsed.Seconds() / percentage_done * (100 - percentage_done)
	// 		//print percentage with 4 decimals
	// 		fmt.Printf("seed: %v/%v (%.4f%%)\teta: %v", seedIndex, seedTotal, percentage_done, time.Duration(eta)*time.Second)
	// 		soil := almanac.seed2soil.Map(int(seed))
	// 		fertilizer := almanac.soil2fertilizer.Map(soil)
	// 		water := almanac.ferilizer2water.Map(fertilizer)
	// 		light := almanac.water2light.Map(water)
	// 		temperature := almanac.light2temperature.Map(light)
	// 		humidity := almanac.temperature2humidity.Map(temperature)
	// 		location := almanac.humidity2location.Map(humidity)
	// 		locations = append(locations, Location(location))
	// 		seedIndex++
	// 	}
	// }

	fmt.Printf("locations: %v\n", locations)

	//smallest loc
	smallest := locations[0]
	for _, loc := range locations {
		if loc < smallest {
			smallest = loc
		}
	}

	fmt.Printf("smallest: %v\n", smallest)

	return "", nil
}
