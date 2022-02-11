package main

import (
	"flag"
	"fmt"
	"github.com/pierredavidbelanger/raftman/engine"
	"log"
	"net/url"
)

func main() {

	var frontendArgs URLValues
	var backendArgs URLValues

	flag.Var(&frontendArgs, "frontend", "Frontend URLs")
	flag.Var(&backendArgs, "backend", "Backend URL")

	flag.Parse()

	if len(backendArgs) == 0 {
		backendArgs = append(backendArgs, mustParseURL("sqlite:///Users/yumm/temp/logs.db"))
	} else if len(backendArgs) > 1 {
		log.Fatal("At most one backend must be defined")
	}

	if len(frontendArgs) == 0 {
		frontendArgs = append(frontendArgs, mustParseURL("syslog+udp://:5513"))
		frontendArgs = append(frontendArgs, mustParseURL("syslog+tcp://:5515"))
		frontendArgs = append(frontendArgs, mustParseURL("api+http://:18181/api/"))
		frontendArgs = append(frontendArgs, mustParseURL("ui+http://:18282/"))
	}

	e, err := engine.NewEngine(backendArgs[0], frontendArgs)
	if err != nil {
		log.Fatalf("Unable to create engine: %s", err)
	}

	if err = e.Start(); err != nil {
		log.Fatalf("Unable to start engine: %s", err)
	}
	defer e.Close()

	e.Wait()
}

type URLValues []*url.URL

func (s *URLValues) String() string {
	return fmt.Sprintf("%+v", *s)
}

func (s *URLValues) Set(value string) error {
	parsed, err := url.Parse(value)
	if err != nil {
		return err
	}
	*s = append(*s, parsed)
	return nil
}

func mustParseURL(value string) *url.URL {
	parsed, err := url.Parse(value)
	if err != nil {
		log.Fatalf("Unable to parse URL: %s", err)
	}
	return parsed
}
