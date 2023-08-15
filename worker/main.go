package main

import (
	"log"
	"os"

	"github.com/surajkumarsinha/go-temporal/app"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

var (
	stripeKey     = os.Getenv("STRIPE_PRIVATE_KEY")
	mailgunDomain = os.Getenv("MAILGUN_DOMAIN")
	mailgunKey    = os.Getenv("MAILGUN_PRIVATE_KEY")
)

func main() {
	// Create the client object just once per process
	c, err := client.NewLazyClient(client.Options{})
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}
	defer c.Close()
	// This worker hosts both Worker and Activity functions
	w := worker.New(c, "CART_TASK_QUEUE", worker.Options{})

	if stripeKey == "" {
		log.Fatalln("Must set STRIPE_PRIVATE_KEY environment variable")
	}
	if mailgunDomain == "" {
		log.Fatalln("Must set MAILGUN_DOMAIN environment variable")
	}
	if mailgunKey == "" {
		log.Fatalln("Must set MAILGUN_PRIVATE_KEY environment variable")
	}

	a := &app.Activities{
		StripeKey:     stripeKey,
		MailgunDomain: mailgunDomain,
		MailgunKey:    mailgunKey,
	}

	w.RegisterActivity(a.CreateStripeCharge)
	w.RegisterActivity(a.SendAbandonedCartEmail)
	w.RegisterWorkflow(app.CartWorkflow)
	// Start listening to the Task Queue
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("unable to start Worker", err)
	}
}
