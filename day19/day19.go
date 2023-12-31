package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"golang.org/x/exp/maps"
)

type State struct {
	timeRemaining int
	ore           int
	clay          int
	obsidian      int
	geode         int
	oreBots       int
	clayBots      int
	obsidianBots  int
	geodeBots     int
}

type Blueprint struct {
	id                int
	oreOreCost        int
	clayOreCost       int
	obsidianOreCost   int
	obsidianClayCost  int
	geodeOreCost      int
	geodeObsidianCost int

	maxOreBots      int
	maxClayBots     int
	maxObsidianBots int
}

var blueprints []Blueprint

var cacheHit int
var cacheMiss int
var stateBestResultCache map[State]int

func loadPuzzle(file *os.File) {
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		id := 0
		oreOreCost := 0
		clayOreCost := 0
		obsidianOreCost := 0
		obsidianClayCost := 0
		geodeOreCost := 0
		geodeObsidianCost := 0

		fmt.Sscanf(
			line,
			"Blueprint %d: Each ore robot costs %d ore.  Each clay robot costs %d ore.  Each obsidian robot costs %d ore and %d clay.  Each geode robot costs %d ore and %d obsidian.",
			&id, &oreOreCost, &clayOreCost, &obsidianOreCost, &obsidianClayCost, &geodeOreCost, &geodeObsidianCost,
		)

		// we will never need more of a given bot than it takes to produce the most expensive bot that uses that resource
		// for example, consider oreBots, where clay bots cost 2 ore, obsidian bots cost 3 ore, and geode bots cost 2 ore
		// if we have 3 ore bots, then we can make any of those bots on a given turn.  by the time we do that and another
		// turn comes along, we'll still have 3 ore bots, so we can make any of those bots again.  so we'll never need more

		maxOreBotsAllowed := oreOreCost
		if clayOreCost > maxOreBotsAllowed {
			maxOreBotsAllowed = clayOreCost
		}
		if obsidianOreCost > maxOreBotsAllowed {
			maxOreBotsAllowed = obsidianOreCost
		}
		if geodeOreCost > maxOreBotsAllowed {
			maxOreBotsAllowed = geodeOreCost
		}

		maxClayBotsAllowed := clayOreCost
		if obsidianClayCost > maxClayBotsAllowed {
			maxClayBotsAllowed = obsidianClayCost
		}
		if geodeObsidianCost > maxClayBotsAllowed {
			maxClayBotsAllowed = geodeObsidianCost
		}

		maxObsidianBotsAllowed := obsidianClayCost
		if geodeObsidianCost > maxObsidianBotsAllowed {
			maxObsidianBotsAllowed = geodeObsidianCost
		}

		blueprint := Blueprint{
			id:              id,
			oreOreCost:      oreOreCost,
			clayOreCost:     clayOreCost,
			obsidianOreCost: obsidianOreCost, obsidianClayCost: obsidianClayCost,
			geodeOreCost: geodeOreCost, geodeObsidianCost: geodeObsidianCost,
			maxOreBots:      maxOreBotsAllowed,
			maxClayBots:     maxClayBotsAllowed,
			maxObsidianBots: maxObsidianBotsAllowed,
		}
		blueprints = append(blueprints, blueprint)
	}
}

func copyState(state *State) *State {
	return &State{
		timeRemaining: state.timeRemaining,
		ore:           state.ore,
		clay:          state.clay,
		obsidian:      state.obsidian,
		geode:         state.geode,
		oreBots:       state.oreBots,
		clayBots:      state.clayBots,
		obsidianBots:  state.obsidianBots,
		geodeBots:     state.geodeBots,
	}
}

func deriveChildStates(state *State, blueprint *Blueprint) []*State {
	// baseline outcome, no action but time passes and bots produce
	idleState := copyState(state)
	idleState.timeRemaining--
	idleState.ore += idleState.oreBots
	idleState.clay += idleState.clayBots
	idleState.obsidian += idleState.obsidianBots
	idleState.geode += idleState.geodeBots

	buildWillFinishAndProduce := state.timeRemaining >= 2
	if !buildWillFinishAndProduce {
		return []*State{idleState}
	}

	canMakeOreBot := state.ore >= blueprint.oreOreCost
	canMakeClayBot := state.ore >= blueprint.clayOreCost
	canMakeObsidianBot := (state.ore >= blueprint.obsidianOreCost) && (state.clay >= blueprint.obsidianClayCost)
	canMakeGeodeBot := (state.ore >= blueprint.geodeOreCost) && (state.obsidian >= blueprint.geodeObsidianCost)

	// if we can make a geode bot, that's the only thing we should do
	if canMakeGeodeBot {
		makeGeodeBotState := copyState(idleState)
		makeGeodeBotState.geodeBots++
		makeGeodeBotState.ore -= blueprint.geodeOreCost
		makeGeodeBotState.obsidian -= blueprint.geodeObsidianCost
		return []*State{makeGeodeBotState}
	}

	// if we can already make every type of bot, then there is no point in idling.  we should make a bot of some sort
	idleWillEnableGeodeBot := !canMakeGeodeBot && (state.ore < blueprint.geodeOreCost) && (state.obsidian < blueprint.geodeObsidianCost) && (state.oreBots > 0) && (state.obsidianBots > 0)
	idleWillEnableObsidianBot := !canMakeObsidianBot && (state.ore < blueprint.obsidianOreCost) && (state.clay < blueprint.obsidianClayCost) && (state.oreBots > 0) && (state.clayBots > 0)
	idleWillEnableClayBot := !canMakeClayBot && (state.ore < blueprint.clayOreCost) && (state.oreBots > 0)
	idleWillEnableOreBot := !canMakeOreBot && (state.ore < blueprint.oreOreCost)
	idleAllowed := idleWillEnableOreBot || idleWillEnableClayBot || idleWillEnableObsidianBot || idleWillEnableGeodeBot

	states := make([]*State, 0)
	if idleAllowed {
		states = append(states, idleState)
	}

	if canMakeObsidianBot && state.obsidianBots < blueprint.maxObsidianBots {
		makeObsidianBotState := copyState(idleState)
		makeObsidianBotState.obsidianBots++
		makeObsidianBotState.ore -= blueprint.obsidianOreCost
		makeObsidianBotState.clay -= blueprint.obsidianClayCost
		states = append(states, makeObsidianBotState)
	}

	if canMakeClayBot && state.clayBots < blueprint.maxClayBots {
		makeClayBotState := copyState(idleState)
		makeClayBotState.clayBots++
		makeClayBotState.ore -= blueprint.clayOreCost
		states = append(states, makeClayBotState)
	}

	if canMakeOreBot && state.oreBots < blueprint.maxOreBots {
		makeOreBotState := copyState(idleState)
		makeOreBotState.oreBots++
		makeOreBotState.ore -= blueprint.oreOreCost
		states = append(states, makeOreBotState)
	}

	return states
}

func calculateMaxGeodes(state *State, blueprint *Blueprint) int {
	if state.timeRemaining <= 0 {
		return state.geode
	}

	cachedValue, ok := stateBestResultCache[*state]
	if ok {
		cacheHit++
		return cachedValue
	}
	cacheMiss++

	childStates := deriveChildStates(state, blueprint)

	maxChildGeodes := 0
	for _, childState := range childStates {
		childGeodes := calculateMaxGeodes(childState, blueprint)
		if childGeodes > maxChildGeodes {
			maxChildGeodes = childGeodes
		}
	}

	stateBestResultCache[*state] = maxChildGeodes
	return maxChildGeodes
}

func (state *State) print() {
	fmt.Printf("State: %d time remaining, %d ore, %d clay, %d obsidian, %d geode, %d ore bots, %d clay bots, %d obsidian bots, %d geode bots\n",
		state.timeRemaining, state.ore, state.clay, state.obsidian, state.geode, state.oreBots, state.clayBots, state.obsidianBots, state.geodeBots)
}

func (blueprint *Blueprint) print() {
	fmt.Printf("Blueprint %d: Each ore robot costs %d ore.  Each clay robot costs %d ore.  Each obsidian robot costs %d ore and %d clay.  Each geode robot costs %d ore and %d obsidian.\n",
		blueprint.id, blueprint.oreOreCost, blueprint.clayOreCost, blueprint.obsidianOreCost, blueprint.obsidianClayCost, blueprint.geodeOreCost, blueprint.geodeObsidianCost)
}

func part1() {
	start := time.Now()
	timeAlloted := 24
	totalQuality := 0
	stateBestResultCache = make(map[State]int)

	for _, blueprint := range blueprints {
		blueprintStart := time.Now()
		cacheHit = 0
		cacheMiss = 0
		maps.Clear(stateBestResultCache)
		blueprint.print()

		initialState := State{timeRemaining: timeAlloted, oreBots: 1}
		maxGeodes := calculateMaxGeodes(&initialState, &blueprint)
		fmt.Printf("Blueprint %d: Max geodes: %d. Cache hit: %d, miss: %d\n", blueprint.id, maxGeodes, cacheHit, cacheMiss)
		fmt.Printf("Blueprint time: %s\n", time.Since(blueprintStart))
		totalQuality += maxGeodes * blueprint.id
	}

	println(totalQuality)
	fmt.Printf("Time: %s\n", time.Since(start))
}

func part2() {
	start := time.Now()
	timeAlloted := 32
	blueprints = blueprints[0:3]
	outputProduct := 1
	stateBestResultCache = make(map[State]int)

	for _, blueprint := range blueprints {
		blueprintStart := time.Now()
		cacheHit = 0
		cacheMiss = 0
		maps.Clear(stateBestResultCache)
		blueprint.print()

		initialState := State{timeRemaining: timeAlloted, oreBots: 1}
		maxGeodes := calculateMaxGeodes(&initialState, &blueprint)
		fmt.Printf("Blueprint %d: Max geodes: %d. Cache hit: %d, miss: %d\n", blueprint.id, maxGeodes, cacheHit, cacheMiss)
		fmt.Printf("Blueprint time: %s\n", time.Since(blueprintStart))
		outputProduct *= maxGeodes
	}

	println(outputProduct)
	fmt.Printf("Time: %s\n", time.Since(start))
}

func main() {
	filename := "input.txt"
	method := part1
	for _, v := range os.Args {
		if v == "part2" || v == "2" {
			method = part2
		}
		if strings.HasSuffix(v, ".txt") {
			filename = v
		}
	}

	file, _ := os.Open(filename)
	defer file.Close()
	loadPuzzle(file)
	method()
}
