# dhook

`dhook` is a small utility to send messages to Discord via a webhook endpoint.

## Usage:

`/path/to/dhook [options]`

### Command Line Parameters

|Parameter|Alias(es)|Description|
|---------|---------|-----------|
|-message|-msg|The message content to send. Required if you do not set an embed option.|
|-webhook-url|-url|Required [1]: The webhook URL to send the message to.|
|-username|-user|The username to send the webhook as.|
|-avatar-url|-avatar|URL of an image to use as the avatar.|
|-embed-url||The URL for the embed title.|
|-embed-title|-title|The title of the embed.|
|-embed-description|-description|The description text in the embed object.|
|-embed-colour|-colour|The colour of the embed [2].|
|-embed-footer-text|-footer|The text of the footer in the embed object.|
|-embed-footer-icon|-footer-icon|A link to the icon to be used in the footer of the embed object. Discord requires `-embed-footer-text` when you set this.|
|-version||Print the version and exit.|
|-debug||Write debug messages to stderr. The messages do not include the webhook URL.|

[1] `-webhook-url` may be set in one of the following ways - List below is in order of preference:

  1. Command-line parameter
  1. Golang Package Variable: `main.webhookURL`
  1. Environment variable: `DHOOK_URL`

[2] `-embed-colour`:
  * Also has the following aliases: `-embed-color` and `-color`
  * This is the decimal value, e.g. `-embed-colour 34047`
  * A good tool to use is https://www.spycolor.com/.

### Exit status

`dhook` exits with status 0 when the server accepts the message. It exits with a non-zero status when it cannot send the message, for example when the webhook URL is missing, the request fails, or the server returns an error.

## Docker

The release workflow publishes an image to `ghcr.io/zerocube/dhook`, with the tags `latest`, the version number (for example `0.2.0`) and the commit SHA. The image runs `dhook`, so give the options after the image name:

```bash
docker run --rm -e DHOOK_URL="https://discord.com/api/webhooks/..." ghcr.io/zerocube/dhook:latest -msg "Hello"
```

To build an image that contains a webhook URL, set the `WEBHOOK_URL` build argument. The URL contains the webhook token, so do not publish that image.

```bash
docker build --build-arg WEBHOOK_URL="https://discord.com/api/webhooks/..." -t dhook .
```
