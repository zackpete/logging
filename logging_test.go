package logging

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"
	"zack.wtf/lib/chrono"
)

func Test(*testing.T) {
	check("{number} is a number", 3)
	check("this doesn't have any args")
	check("this has {one} extra arg", 1, "extra")
	check("this has invalid stuff { } in brackets {end}", "(end)")
	check("this is {missing} an arg")
	check("this has {two} {{escaped}} braces", 2)
	check("{n|%.2f} has custom formatting", 42.235235)
	check("is is now {time}", time.Now())
	check("{result} is lazily evaluated", func() int { return 2 + 2 })
}

func check(template string, args ...interface{}) {
	message, details := convert(template, args)
	fmt.Println(message)
	fmt.Printf("%v\n\n", details)
}

func Benchmark(b *testing.B) {
	const size = 100

	levels := make([]Level, size)
	templates := make([]string, size)
	args := make([][]interface{}, size)

	for i := 0; i < size; i++ {
		levels[i], templates[i], args[i] = generateMessage()
	}

	log := &Log{
		Level: None,
		Sink:  SinkFunc(func(*Event) {}),
		Clock: chrono.SystemClock,
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		n := i % size
		log.Emit(levels[n], templates[n], args[n])
	}
}

func generateMessage() (Level, string, []interface{}) {
	numWords := 3 + rand.Intn(5)

	words := make([]string, 0, numWords)
	for i := 0; i < numWords; i++ {
		words = append(words, wordList[rand.Intn(len(wordList))])
	}

	numMetrics := rand.Intn(3)
	metrics := make([]interface{}, 0, numMetrics)
	for i := 0; i < numMetrics; i++ {
		metrics = append(metrics, metricList[rand.Intn(len(metricList))])
	}

	taken := make([]bool, numWords)

	for remaining := numMetrics; remaining > 0; {
		if put := rand.Intn(numWords); !taken[put] {
			words[put] = fmt.Sprintf("{%s}", words[put])
			taken[put] = true
			remaining--
		}
	}

	return Level(rand.Intn(int(Debug)) + 1), strings.Join(words, " "), metrics
}

var metricList = []interface{}{
	1,
	"three",
	4.44444,
	time.Now(),
	time.Millisecond,
}

var wordList = []string{
	"abandon",
	"ability",
	"able",
	"abortion",
	"about",
	"above",
	"abroad",
	"absence",
	"absolute",
	"absolutely",
	"absorb",
	"abuse",
	"academic",
	"accept",
	"access",
	"accident",
	"accompany",
	"accomplish",
	"according",
	"account",
	"accurate",
	"accuse",
	"achieve",
	"achievement",
	"acid",
	"acknowledge",
	"acquire",
	"across",
	"act",
	"action",
	"active",
	"activist",
	"activity",
	"actor",
	"actress",
	"actual",
	"actually",
	"ad",
	"adapt",
	"add",
	"addition",
	"additional",
	"address",
	"adequate",
	"adjust",
	"adjustment",
	"administration",
	"administrator",
	"admire",
	"admission",
	"admit",
	"adolescent",
	"adopt",
	"adult",
	"advance",
	"advanced",
	"advantage",
	"adventure",
	"advertising",
	"advice",
	"advise",
	"adviser",
	"advocate",
	"affair",
	"affect",
	"afford",
	"afraid",
	"African",
	"African-American",
	"after",
	"afternoon",
	"again",
	"against",
	"age",
	"agency",
	"agenda",
	"agent",
	"aggressive",
	"ago",
	"agree",
	"agreement",
	"agricultural",
	"ah",
	"ahead",
	"aid",
	"aide",
	"AIDS",
	"aim",
	"air",
	"aircraft",
	"airline",
	"airport",
	"album",
	"alcohol",
	"alive",
	"all",
	"alliance",
	"allow",
	"ally",
	"almost",
	"alone",
	"along",
	"already",
	"also",
	"alter",
	"alternative",
	"although",
	"always",
	"AM",
	"amazing",
	"American",
	"among",
	"amount",
	"analysis",
	"analyst",
	"analyze",
	"ancient",
	"and",
	"anger",
	"angle",
	"angry",
	"animal",
	"anniversary",
	"announce",
	"annual",
	"another",
	"answer",
	"anticipate",
	"anxiety",
	"any",
	"anybody",
	"anymore",
	"anyone",
	"anything",
	"anyway",
	"anywhere",
	"apart",
	"apartment",
	"apparent",
	"apparently",
	"appeal",
	"appear",
	"appearance",
	"apple",
	"application",
	"apply",
	"appoint",
	"appointment",
	"appreciate",
	"approach",
	"appropriate",
	"approval",
	"approve",
	"approximately",
	"Arab",
	"architect",
	"area",
	"argue",
	"argument",
	"arise",
	"arm",
	"armed",
	"army",
	"around",
	"arrange",
	"arrangement",
	"arrest",
	"arrival",
	"arrive",
	"art",
	"article",
	"artist",
	"artistic",
	"as",
	"Asian",
	"aside",
	"ask",
	"asleep",
	"aspect",
	"assault",
	"assert",
	"assess",
	"assessment",
	"asset",
	"assign",
	"assignment",
	"assist",
	"assistance",
	"assistant",
	"associate",
	"association",
	"assume",
	"assumption",
	"assure",
	"at",
	"athlete",
	"athletic",
	"atmosphere",
	"attach",
	"attack",
	"attempt",
	"attend",
	"attention",
	"attitude",
	"attorney",
	"attract",
	"attractive",
	"attribute",
	"audience",
	"author",
	"authority",
	"auto",
	"available",
	"average",
	"avoid",
	"award",
	"aware",
	"awareness",
	"away",
	"awful",
}
