package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"flag"

	"github.com/bwmarrin/discordgo"
	"os"
	"os/signal"
	"syscall"
)

type ConfigFile struct {
	HelpEnabled bool `json:"helpEnabled"`
	IsUserBot bool `json:"IsUserBot"`
	Memes       []struct {
		Command string `json:"Command"`
		APIURL  string `json:"api_url"`
	} `json:"memes"`
	BlockedUsers []struct {
		Username  string `json:"Username"`
		DiscordID string `json:"DiscordID"`
	} `json:"blocked_users"`
}
type meme struct {
	URL string `json:"link"`
}
type response struct {
	Meme meme `json:"meme"`
}

// Variables used for command line parameters
var (
	Token string
)

func init() {

	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {

	//print(getMemeURL(ReadConfig("./config.json","uff")))

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New(getBotType("./config.json"))
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	data, err := ioutil.ReadFile("./config.json")
	if err != nil {
		fmt.Print(err)
	}

	var obj ConfigFile
	err = json.Unmarshal(data, &obj)
	if err != nil {
		fmt.Println("error:", err)
	}

	for _, element := range obj.BlockedUsers {
		if element.DiscordID == m.Author.ID {
			println(m.Author.Username + " (blocked)")
			return
		}
	}

	if helpEnabled("./config.json") && m.Content == "_help" {
		Help := "" +
			"**__Commands / UFFBot__**" +
			"```"
		for _, element := range obj.Memes {
			Help += element.Command + "\n"
		}
		Help += "```"
		s.ChannelMessageSend(m.ChannelID, Help)
		println(m.Author.Username + " issued Command " + m.Content)
		return
	}
	for _, element := range obj.Memes {
		if strings.EqualFold(m.Content, element.Command) {
			s.ChannelMessageSend(m.ChannelID, getMemeURL(ReadConfig("./config.json", element.Command)))
			println(m.Author.Username + " issued Command " + m.Content)
			return
		}
	}

}

func ReadConfig(fileURL string, command string) string {
	data, err := ioutil.ReadFile(fileURL)
	if err != nil {
		fmt.Print(err)
	}

	var obj ConfigFile

	err = json.Unmarshal(data, &obj)
	if err != nil {
		fmt.Println("error:", err)
	}

	for _, element := range obj.Memes {
		if element.Command == command {
			return element.APIURL
		}
	}
	return ""
}
func helpEnabled(fileURL string) bool {
	data, err := ioutil.ReadFile(fileURL)
	if err != nil {
		fmt.Print(err)
	}

	var obj ConfigFile

	err = json.Unmarshal(data, &obj)
	if err != nil {
		fmt.Println("error:", err)
	}

	if obj.HelpEnabled {
		return true
	}
	return false
}

func getMemeURL(url string) string {
	spaceClient := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "UFF-BOT")

	res, getErr := spaceClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	MemeOjb := response{}
	jsonErr := json.Unmarshal(body, &MemeOjb)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	return MemeOjb.Meme.URL
}
func getBotType(fileURL string) string {
	data, err := ioutil.ReadFile(fileURL)
	if err != nil {
		fmt.Print(err)
	}

	var obj ConfigFile

	err = json.Unmarshal(data, &obj)
	if err != nil {
		fmt.Println("error:", err)
	}

	if obj.IsUserBot {
		println("Bot-Type: Userbot")
		return "" + Token
	}
	println("Bot-Type: Standard")
	return "Bot " + Token
}
