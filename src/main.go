package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

type KeyboardButton struct {
	Text            string `json:"text"`
	RequestContact  bool   `json:"request_contact"`
	RequestLocation bool   `json:"request_location"`
}

type InlineKeyboardButton struct {
	Text         string `json:"text"`
	CallbackData string `json:"callback_data"`
	Url          string `json:"url"`
}

type ReplyMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
	Keyboard       [][]KeyboardButton       `json:"keyboard"`
	ResizeKeyboard bool                     `json:"resize_keyboard"`
	OnTimeKeyboard bool                     `json:"one_time_keyboard"`
	Selective      bool                     `json:"selective"`
}

type SendMessage struct {
	ChatID      interface{} `json:"chat_id"`
	Text        string      `json:"text"`
	ParseMode   string      `json:"parse_mode"`
	ReplyMarkup interface{} `json:"reply_markup"`
}

func ReadSendMessageRequest(fileName string) (*SendMessage, error) {
	fileContent, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
	}
	defer fileContent.Close()

	byteValue, _ := ioutil.ReadAll(fileContent)
	var replyMarkup ReplyMarkup
	var message SendMessage
	json.Unmarshal(byteValue, &message)

	for key, value := range message.ReplyMarkup.(map[string]interface{}) {
		if key == "inline_keyboard" {
			var inlineKeyboardss [][]InlineKeyboardButton
			for _, x := range value.([]interface{}) {
				var inlineKeyboards []InlineKeyboardButton
				for _, y := range x.([]interface{}) {
					var inlineKeyboard InlineKeyboardButton
					for a, z := range y.(map[string]interface{}) {
						if a == "text" {
							inlineKeyboard.Text = z.(string)
						}
						if a == "callback_data" {
							inlineKeyboard.CallbackData = z.(string)
						}
						if a == "url" {
							inlineKeyboard.Url = z.(string)
						}
					}
					inlineKeyboards = append(inlineKeyboards, inlineKeyboard)
				}
				inlineKeyboardss = append(inlineKeyboardss, inlineKeyboards)
			}
			replyMarkup.InlineKeyboard = inlineKeyboardss
		} else {
			var keyboardss [][]KeyboardButton
			for _, x := range value.([]interface{}) {
				var keyboards []KeyboardButton
				for _, y := range x.([]interface{}) {
					var keyboard KeyboardButton
					for a, z := range y.(map[string]interface{}) {
						if a == "text" {
							keyboard.Text = z.(string)
						}
						if a == "request_contact" {
							keyboard.RequestContact = z.(bool)
						}
						if a == "request_location" {
							keyboard.RequestLocation = z.(bool)
						}
					}
					keyboards = append(keyboards, keyboard)
				}
				keyboardss = append(keyboardss, keyboards)
			}
			replyMarkup.Keyboard = keyboardss
		}
	}
	message.ReplyMarkup = replyMarkup
	if message.ChatID == nil {
		return nil, errors.New("chat_id is empty")
	}
	if len(message.Text) == 0 {
		return nil, errors.New("text is empty")
	}
	return &message, nil
}
