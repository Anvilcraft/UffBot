package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Memes struct {
	Command string
	APIUrl  string
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

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("" + Token)
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

	UffMemes := []Memes{
		{Command: "uff", APIUrl: "https://jensmemes.tilera.xyz/api/random?category=uff"},
		{Command: "uffat", APIUrl: "https://jensmemes.tilera.xyz/api/meme?id=257"},
		{Command: "uffgo", APIUrl: "https://jensmemes.tilera.xyz/api/meme?id=568"},
		{Command: "hey", APIUrl: "https://jensmemes.tilera.xyz/api/meme?id=243"},
		{Command: "uffch", APIUrl: "https://jensmemes.tilera.xyz/api/meme?id=283"},
		{Command: "drogen", APIUrl: "https://jensmemes.tilera.xyz/api/meme?id=80"},
		{Command: "kappa", APIUrl: "https://jensmemes.tilera.xyz/api/meme?id=182"},
		{Command: "hendrik", APIUrl: "https://jensmemes.tilera.xyz/api/meme?id=248"},
		{Command: "ufflie", APIUrl: "https://jensmemes.tilera.xyz/api/meme?id=284"},
		{Command: "uffns", APIUrl: "https://jensmemes.tilera.xyz/api/meme?id=287"},
		{Command: "uffhs", APIUrl: "https://jensmemes.tilera.xyz/api/meme?id=331"},
		{Command: "uffde", APIUrl: "https://jensmemes.tilera.xyz/api/meme?id=144"},
		{Command: "uffhre", APIUrl: "https://jensmemes.tilera.xyz/api/meme?id=312"},
		{Command: "uffpy", APIUrl: "https://jensmemes.tilera.xyz/api/meme?id=477"},
		{Command: "itbyhf", APIUrl: "https://jensmemes.tilera.xyz/api/meme?id=314"},
		{Command: "tilera", APIUrl: "https://jensmemes.tilera.xyz/api/meme?id=316"},
		{Command: "lordmzte", APIUrl: "https://jensmemes.tilera.xyz/api/meme?id=315"},
		{Command: "realtox", APIUrl: "https://jensmemes.tilera.xyz/api/meme?id=168"},
		{Command: "jonasled", APIUrl: "https://jensmemes.tilera.xyz/api/meme?id=164"},
		{Command: "sklave", APIUrl: "https://jensmemes.tilera.xyz/api/meme?id=304"},
		{Command: "jens", APIUrl: "https://jensmemes.tilera.xyz/api/meme?id=343"},
		{Command: "fresse", APIUrl: "https://jensmemes.tilera.xyz/api/meme?id=375"},
		{Command: "bastard", APIUrl: "https://jensmemes.tilera.xyz/api/meme?id=375"},
		{Command: "uffsr", APIUrl: "https://jensmemes.tilera.xyz/api/meme?id=585"},
		{Command: "party", APIUrl: "https://jensmemes.tilera.xyz/api/meme?id=619"},
		{Command: "uffrs", APIUrl: "https://jensmemes.tilera.xyz/api/meme?id=616"},
		{Command: "uffjs", APIUrl: "https://jensmemes.tilera.xyz/api/meme?id=615"},
	}

	Help := "" +
		"**__Commands / UFFBot__**" +
		"```"
	for _, meme := range UffMemes {
		Help += meme.Command + "\n"
	}
	Help += "```"

	for _, meme := range UffMemes {
		if strings.Title(m.Content) == strings.Title("_help") {
			s.ChannelMessageSend(m.ChannelID, Help)
			break
		}
		if strings.Title(m.Content) == strings.Title(meme.Command) {
			s.ChannelMessageSend(m.ChannelID, getMemeURL(meme.APIUrl))
		}
	}

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

	memeob := response{}
	jsonErr := json.Unmarshal(body, &memeob)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	return memeob.Meme.URL
}
