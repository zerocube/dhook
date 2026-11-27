package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var webhookURL string

const webhookURLEnvKey = "DHOOK_URL"

// Change this for a release. The -X linker flag cannot change a const.
const version = "0.2.0"

func main() {
	// Flags for the main payload - Message and username
	message := flag.String("message", "", "The message to send")
	flag.StringVar(message, "msg", *message, "alias for -message")

	username := flag.String("username", "dhook", "The username to send the message as")
	flag.StringVar(username, "user", *username, "alias for -username")

	avatarURL := flag.String("avatar-url", "", "The URL of the avatar to use")
	flag.StringVar(avatarURL, "avatar", *avatarURL, "alias for -avatar-url")

	webhookURLFlag := flag.String("webhook-url", "", "The webhook URL to send the message to")
	flag.StringVar(webhookURLFlag, "url", *webhookURLFlag, "alias for -webhook-url")

	// Flags for the main embed struct - URL, Title, and Description
	embedURL := flag.String("embed-url", "", "The URL for the embed title")

	embedTitle := flag.String("embed-title", "", "The embed title")
	flag.StringVar(embedTitle, "title", *embedTitle, "alias for -embed-title")
	embedDescription := flag.String("embed-description", "", "The embed description")
	flag.StringVar(embedDescription, "description", *embedDescription, "alias for -embed-description")
	embedColour := flag.Int("embed-colour", 0, "The embed colour")
	flag.IntVar(embedColour, "colour", *embedColour, "alias for -embed-colour")
	flag.IntVar(embedColour, "embed-color", *embedColour, "alias for -embed-colour")
	flag.IntVar(embedColour, "color", *embedColour, "alias for -embed-colour")

	// Flags for the embed footer struct - Footer (text), and Icon (URL)
	embedFooterText := flag.String("embed-footer-text", "", "The embed footer text")
	flag.StringVar(embedFooterText, "footer", *embedFooterText, "alias for -embed-footer-text")

	embedFooterIcon := flag.String("embed-footer-icon", "", "The URL for the footer icon")
	flag.StringVar(embedFooterIcon, "footer-icon", *embedFooterIcon, "alias for -embed-footer-icon")

	showVersion := flag.Bool("version", false, "Shows version information and exits")
	debug := flag.Bool("debug", false, "Enables debug logging")
	flag.Bool("verbose", false, "Has no effect (kept for compatibility)")

	// Work the flag magic
	flag.Parse()

	if *showVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	// Precedence: the command-line flag, then the package variable (set at
	// build time), then the environment variable.
	url, source := *webhookURLFlag, "command-line flag"
	if url == "" {
		url, source = webhookURL, "built-in package variable (set with -ldflags)"
	}
	if url == "" {
		url, source = os.Getenv(webhookURLEnvKey), webhookURLEnvKey+" environment variable"
	}
	if url == "" {
		fmt.Fprintln(os.Stderr, "Please provide a webhook URL to send the message to using one of the options below.")
		fmt.Fprintln(os.Stderr, "  Environment variable:", webhookURLEnvKey)
		fmt.Fprintln(os.Stderr, "  Command Line parameter: -webhook-url or -url")
		fmt.Fprintln(os.Stderr, "  Build time: -ldflags \"-X main.webhookURL=...\"")
		os.Exit(1)
	}
	if *debug {
		log.Println("Using the webhook URL from the", source)
	}

	// Leave the footer out unless a footer option is set, so we never send "footer":{}.
	var embedFooter *discordgo.MessageEmbedFooter
	if *embedFooterText != "" || *embedFooterIcon != "" {
		embedFooter = &discordgo.MessageEmbedFooter{
			Text:    *embedFooterText,
			IconURL: *embedFooterIcon,
		}
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
				Footer:      embedFooter,
				Color:       *embedColour,
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

	if strings.TrimSpace(*message) == "" && webhookPayload.Embeds == nil {
		log.Fatalln("Error: Please provide a message (-message) or embed content (for example -title, -description or -footer).")
	}

	sendJSONPayload(&url, webhookPayload)
}
