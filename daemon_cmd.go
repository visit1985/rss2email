//
// This is the daemon-subcommand.
//

package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/skx/rss2email/processor"
)

// Structure for our options and state.
type daemonCmd struct {

	// Should we be verbose in operation?
	verbose bool
}

// Info is part of the subcommand-API.
func (d *daemonCmd) Info() (string, string) {
	return "daemon", `Send emails for each new entry in our feed lists.

This sub-command polls all configured feeds, sending an email for
each item which is new.  Once the list of feeds has been processed
the command will pause for 15 minutes, before beginning again.

In terms of implementation this command follows everything documented
in the 'cron' sub-command.  The only difference is this one never
terminates - even if email-generation fails.


Example:

    $ rss2email daemon user1@example.com user2@example.com
`
}

// Arguments handles our flag-setup.
func (d *daemonCmd) Arguments(f *flag.FlagSet) {
	f.BoolVar(&d.verbose, "verbose", false, "Should we be extra verbose?")
}

//
// Entry-point
//
func (d *daemonCmd) Execute(args []string) int {

	// No argument?  That's a bug
	if len(args) == 0 {
		fmt.Printf("Usage: rss2email daemon email1@example.com .. emailN@example.com\n")
		return 1
	}

	// The list of addresses to which we should send our notices.
	recipients := []string{}

	// Save each argument away, checking it is fully-qualified.
	for _, email := range args {
		if strings.Contains(email, "@") {
			recipients = append(recipients, email)
		} else {
			fmt.Printf("Usage: rss2email daemon [flags] email1 .. emailN\n")
			return 1
		}
	}

	for true {

		// Create the helper
		p := processor.New()

		// Setup the state - note we ALWAYS send emails in this mode.
		p.SetVerbose(d.verbose)
		p.SetSendEmail(true)

		errors := p.ProcessFeeds(recipients)

		// If we found errors then show them.
		if len(errors) > 0 {
			for _, err := range errors {
				fmt.Fprintln(os.Stderr, err.Error())
			}
		}

		if d.verbose {
			fmt.Printf("sleeping for 15 minutes")
		}
		time.Sleep(60 * 15 * time.Second)
	}

	// All good.
	return 0
}
