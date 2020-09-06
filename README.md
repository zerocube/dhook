# dhook

`dhook` is a small utility to send messages to Discord via a webhook endpoint.

## Usage:

`/path/to/dhook [options]`

### Command Line Parameters

|Parameter|Alias(es)|Description|
|---------|---------|-----------|
|-message|-msg|Required: The message content to send.|
|-webhook-url|-url|Required [1]: The webhook URL to send the message to.|
|-username|-user|The username to send the webhook as.|
|-avatar-url|-avatar|URL of an image to use as the avatar.|
|-embed-url||The URL for the embed title.|
|-embed-title|-title|The title of the embed.|
|-embed-description|-description|The description text in the embed object.|
|-embed-footer-text|-footer|The text of the footer in the embed object.|
|-embed-footer-icon|-footer-icon|A link to the icon to be used in the footer of the embed object.|

[1] `-webhook-url` may be set in one of the following ways - List below is in order of preference:

  1. Command-line parameter
  1. Golang Package Variable: `main.webhookURL`
  1. Environment variable: `DHOOK_URL`
