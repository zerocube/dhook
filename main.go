package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var debug bool
var verbose bool
var webhookURL string

const webhookURLEnvKey = "DHOOK_URL"

func init() {

	debug = strings.ToLower(os.Getenv("DEBUG")) == "true"
	verbose = strings.ToLower(os.Getenv("VERBOSE")) == "true"

	if verbose {
		log.Println("Verbose mode enabled.")
	}
	if debug {
		log.Println("Debug mode enabled.")
	}

	if webhookURL == "" {
		if debug {
			log.Println("webhookURL package variable empty")
		}
		webhookURL = os.Getenv(webhookURLEnvKey)
		if debug && webhookURL == "" {
			log.Println(webhookURLEnvKey, "environment variable not set.")
		}
	} else {
		if debug {
			log.Println("webhookURL set via package variable:", webhookURL)
		}
	}
}

func main() {

	// Flags for the main payload - Message and username
	message := flag.String("message", "", "The message to send")
	flag.StringVar(message, "msg", *message, "alias for -message")

	username := flag.String("username", "dhook", "The username to send the message as")
	flag.StringVar(username, "user", *username, "alias for -username")

	avatarURL := flag.String("avatar-url", "", "The URL of the avatar to use")
	flag.StringVar(avatarURL, "avatar", *avatarURL, "alias for -avatar-url")

	webhookURLFlag := flag.String("webhook-url", webhookURL, "The webhook URL to send the message to")
	flag.StringVar(webhookURLFlag, "url", *webhookURLFlag, "alias for -webhook-url")

	// Flags for the main embed struct - URL, Title, and Description
	embedURL := flag.String("embed-url", "", "The URL for the embed title")

	embedTitle := flag.String("embed-title", "", "The embed title")
	flag.StringVar(embedTitle, "title", *embedTitle, "alias for -embed-title")
	embedDescription := flag.String("embed-description", "", "The embed description")
	flag.StringVar(embedDescription, "description", *embedDescription, "alias for -embed-description")

	// Flags for the embed footer struct - Footer (text), and Icon (URL)
	embedFooterText := flag.String("embed-footer-text", "", "The embed footer text")
	flag.StringVar(embedFooterText, "footer", *embedFooterText, "alias for -embed-footer-text")

	embedFooterIcon := flag.String("embed-footer-icon", "", "The URL for the footer icon")
	flag.StringVar(embedFooterIcon, "footer-icon", *embedFooterIcon, "alias for -embed-footer-icon")

	// Work the flag magic
	flag.Parse()

	// The command-line parameter for the webhook URL takes precedence over the package / environment variable
	// First check if it's provided in some form.
	if webhookURL == "" && *webhookURLFlag == "" {
		println("Please provide a webhook URL to send the message to using one of the options below.")
		println("  Environment variable:", webhookURLEnvKey)
		println("  Command Line parameter: -webhook-url or -url")
		os.Exit(1)
	}

	var coalescedWebhookURL *string
	if *webhookURLFlag != "" {
		coalescedWebhookURL = webhookURLFlag
	} else {
		coalescedWebhookURL = &webhookURL
	}

	// If it's still empty at this point, something is horribly wrong.
	if *coalescedWebhookURL == "" {
		log.Fatalln("Error: Unable to determine webhook URL.")
	}

	webhookPayload := discordgo.WebhookParams{
		Content:   *message,
		Username:  *username,
		AvatarURL: *avatarURL,
		Embeds: []*discordgo.MessageEmbed{
			{
				URL:         *embedURL,
				Title:       *embedTitle,
				Description: *embedDescription,
				Footer: &discordgo.MessageEmbedFooter{
					Text:    *embedFooterText,
					IconURL: *embedFooterIcon,
				},
			},
		},
	}

	// If there was no embed data, remove it so we don't send a bad request to the webhook URL.
	// We can send an embed with nothing but footer data, so we've got that going for us, which is nice.
	if allStringsAreEmpty([]*string{
		embedURL, embedTitle, embedDescription,
		embedFooterText, embedFooterIcon,
	}) {
		webhookPayload.Embeds = nil
	}

	sendJSONPayload(
		coalescedWebhookURL,
		webhookPayload,
	)
}
