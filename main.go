package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"

    "github.com/go-telegram-bot-api/telegram-bot-api"
)

const weatherAPIURL = "http://api.openweathermap.org/data/2.5/weather"

func main() {
    botToken := os.Getenv("TELEGRAM_TOKEN")
    weatherAPIKey := os.Getenv("WEATHER_API_KEY")

    bot, err := tgbotapi.NewBotAPI(botToken)
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Authorized on account %s", bot.Self.UserName)

    u := tgbotapi.NewUpdate(0)
    u.Timeout = 60

    updates, err := bot.GetUpdatesChan(u)

    for update := range updates {
        if update.Message == nil {
            continue
        }

        msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
        if update.Message.IsCommand() {
            switch update.Message.Command() {
            case "start":
                msg.Text = "Send me the name of a city to get the weather."
            case "weather":
                city := update.Message.CommandArguments()
                if city == "" {
                    msg.Text = "Please provide a city name."
                } else {
                    weather, err := getWeather(city, weatherAPIKey)
                    if err != nil {
                        msg.Text = "Error getting weather data."
                    } else {
                        msg.Text = weather
                    }
                }
            default:
                msg.Text = "I don't know that command."
            }
        } else {
            msg.Text = "Please use commands like /start or /weather <city>."
        }

        if _, err := bot.Send(msg); err != nil {
            log.Printf("Failed to send message: %v", err)
        }
    }
}

func getWeather(city, apiKey string) (string, error) {
    url := fmt.Sprintf("%s?q=%s&appid=%s&units=metric&lang=ru", weatherAPIURL, city, apiKey)
    resp, err := http.Get(url)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    var data map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
        return "", err
    }

    if data["cod"].(float64) != 200 {
        return "City not found", nil
    }

    weather := data["weather"].([]interface{})[0].(map[string]interface{})["description"].(string)
    temp := data["main"].(map[string]interface{})["temp"].(float64)
    return fmt.Sprintf("Weather: %s\nTemperature: %.1f°C", weather, temp), nil
}