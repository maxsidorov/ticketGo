package service

import (
	"errors"
	"strings"
	"unicode"
)

func ValidateName(name string) (error, string) {
	if name == "admin" {
		return nil, name
	}
	if len(name) < 4 || len(name) > 50 {
		return errors.New("Имя должно содержать не менее 4 и не более 50 символов!"), name
	}
	for _, r := range name {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != ' ' && r != '\'' && r != '-' {
			return errors.New("Имя содержит недопустимые символы!"), name
		}
	}
	if strings.Contains(name, "  ") {
		return errors.New("Имя не может содержать несколько пробелов подряд!"), name
	}
	if strings.TrimSpace(name) != name {
		return errors.New("Имя не должно начинаться или заканчиваться пробелом!"), name
	}

	return nil, name
}

func ValidatePassword(name string) (error, string) {
	if (len(name) < 8 || len(name) > 50) && name != "admin" {
		return errors.New("Пароль должен содержать не менее 8 и не более 50 символов!"), name
	}
	for _, r := range name {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != ' ' && r != '\'' && r != '-' {
			return errors.New("Пароль содержит недопустимые символы!"), name
		}
	}
	if strings.Contains(name, "  ") {
		return errors.New("Пароль не может содержать несколько пробелов подряд!"), name
	}
	if strings.TrimSpace(name) != name {
		return errors.New("Пароль не должен начинаться или заканчиваться пробелом!"), name
	}
	return nil, name
}
