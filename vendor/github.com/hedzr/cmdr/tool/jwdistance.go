package tool

import "math"

// Soundex returns the english word's soundex value, such as: 'tags' => 't322'
func Soundex(s string) (snd4 string) {
	return soundex(s)
}

func soundex(s string) (snd4 string) {
	// if len(s) == 0 {
	// 	return
	// }

	var src, tgt []rune
	src = []rune(s)

	i := 0
	for ; i < len(src); i++ {
		if !(src[i] == '-' || src[i] == '~' || src[i] == '+') {
			// first char
			tgt = append(tgt, src[i])
			break
		}
	}

	for ; i < len(src); i++ {
		ch := src[i]
		switch ch {
		case 'a', 'e', 'i', 'o', 'u', 'y', 'h', 'w': // do nothing to remove it
		case 'b', 'f', 'p', 'v':
			tgt = append(tgt, '1')
		case 'c', 'g', 'j', 'k', 'q', 's', 'x', 'z':
			tgt = append(tgt, '2')
		case 'd', 't':
			tgt = append(tgt, '3')
		case 'l':
			tgt = append(tgt, '4')
		case 'm', 'n':
			tgt = append(tgt, '5')
		case 'r':
			tgt = append(tgt, '6')
		}
	}

	snd4 = string(tgt)
	return
}

// StringMetricFactor for JaroWinklerDistance algorithm
const StringMetricFactor = 100000000000

type (
	// StringDistance is an interface for string metric.
	// A string metric is a metric that measures distance between two strings.
	// In most case, it means that the edit distance about those two strings.
	// This is saying, it is how many times are needed while you were
	// modifying string to another one, note that inserting, deleting,
	// substing one character means once.
	StringDistance interface {
		Calc(s1, s2 string, opts ...DistanceOption) (distance int)
	}

	// DistanceOption is a functional options prototype
	DistanceOption func(StringDistance)
)

// JaroWinklerDistance returns an calculator for two strings distance metric, with Jaro-Winkler algorithm.
func JaroWinklerDistance(opts ...DistanceOption) StringDistance {
	x := &jaroWinklerDistance{threshold: 0.7, factor: StringMetricFactor}
	for _, c := range opts {
		c(x)
	}
	return x
}

// JWWithThreshold sets the threshold for Jaro-Winkler algorithm.
func JWWithThreshold(threshold float64) DistanceOption {
	return func(distance StringDistance) {
		if v, ok := distance.(*jaroWinklerDistance); ok {
			v.threshold = threshold
		}
	}
}

type jaroWinklerDistance struct {
	threshold float64
	factor    float64

	matches        int
	maxLength      int
	transpositions int // transpositions is a double number here
	prefix         int

	distance float64
}

func (s *jaroWinklerDistance) Calc(src1, src2 string, opts ...DistanceOption) (distance int) {
	s1, s2 := []rune(src1), []rune(src2)
	lenMax, lenMin := len(s1), len(s2)

	var sMax, sMin []rune
	if lenMax > lenMin {
		sMax, sMin = s1, s2
	} else {
		sMax, sMin = s2, s1
		lenMax, lenMin = lenMin, lenMax
	}
	s.maxLength = lenMax

	iMatchIndexes, matchFlags := s.match(sMax, sMin, lenMax, lenMin)
	s.findTranspositions(sMax, sMin, lenMax, lenMin, iMatchIndexes, matchFlags)

	// println("  matches, transpositions, prefix: ", s.matches, s.transpositions, s.prefix)

	if s.matches == 0 {
		s.distance = 0
		return 0
	}

	m := float64(s.matches)
	jaroDistance := m/float64(lenMax) + m/float64(lenMin)
	jaroDistance += (m - float64(s.transpositions)/2) / m
	jaroDistance /= 3

	var jw float64
	if jaroDistance < s.threshold {
		jw = jaroDistance
	} else {
		jw = jaroDistance + math.Min(0.1, 1/float64(s.maxLength))*float64(s.prefix)*(1-jaroDistance)
	}

	// println("  jaro, jw: ", jaroDistance, jw)

	s.distance = jw * s.factor
	distance = int(math.Round(s.distance))
	return
}

func (s *jaroWinklerDistance) match(sMax, sMin []rune, lenMax, lenMin int) (iMatchIndexes []int, matchFlags []bool) {
	iRange := Max(lenMax/2-1, 0)
	iMatchIndexes = make([]int, lenMin)
	for i := 0; i < lenMin; i++ {
		iMatchIndexes[i] = -1
	}

	s.prefix, s.matches = 0, 0
	for mi := 0; mi < len(sMin); mi++ {
		if sMax[mi] == sMin[mi] {
			s.prefix++
		} else {
			break
		}
	}
	s.matches = s.prefix

	matchFlags = make([]bool, lenMax)

	for mi := s.prefix; mi < lenMin; mi++ {
		c1 := sMin[mi]
		xi, xn := Max(mi-iRange, s.prefix), lenMax // min(mi+iRange-1, lenMax)
		for ; xi < xn; xi++ {
			if !matchFlags[xi] && c1 == sMax[xi] {
				iMatchIndexes[mi] = xi
				matchFlags[xi] = true
				s.matches++
				break
			}
		}
	}
	return
}

func (s *jaroWinklerDistance) findTranspositions(sMax, sMin []rune, lenMax, lenMin int, iMatchIndexes []int, matchFlags []bool) {
	ms1, ms2 := make([]rune, s.matches), make([]rune, s.matches)
	for i, si := 0, 0; i < lenMin; i++ {
		if iMatchIndexes[i] != -1 {
			ms1[si] = sMin[i]
			si++
		}
	}
	for i, si := 0, 0; i < lenMax; i++ {
		if matchFlags[i] {
			ms2[si] = sMax[i]
			si++
		}
	}
	// fmt.Printf("iMatchIndexes, s1, s2: %v, %v, %v\n", iMatchIndexes, string(sMax), string(sMin))
	// println("     ms1, ms2: ", string(ms1), string(ms2))

	s.transpositions = 0
	for mi := 0; mi < len(ms1); mi++ {
		if ms1[mi] != ms2[mi] {
			s.transpositions++
		}
	}
}
