package utils

import (
	"strings"

	my1562client "github.com/my1562/client"
)

func FormatServiceMessage(status *my1562client.GetStatusAPIResponse) string {

	if !status.HasMessage {
		return ""
	}

	var out strings.Builder

	for i, message := range status.Messages {
		out.WriteString(message.Title)
		out.WriteString("\n")
		out.WriteString(message.Description)
		if i < len(status.Messages)-1 {
			out.WriteString("\n")
			out.WriteString("\n")
		}
	}

	return out.String()
}
