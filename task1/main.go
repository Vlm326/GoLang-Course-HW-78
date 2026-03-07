package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
	"bufio"
)


// В го нету функции map, ужас
func Map[T, U any](input []T, f func(T) U) []U {
	result := make([]U, len(input))
	for i, v := range input {
		result[i] = f(v)
	}
	return result
}

const githubAPIURLTemplate = "https://api.github.com/repos/%s/%s"

type repositoryInfo struct {
	Name            string `json:"name"`
	Description     string `json:"description"`
	StargazersCount int    `json:"stargazers_count"`
	ForksCount      int    `json:"forks_count"`
	CreatedAt       string `json:"created_at"`
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("GitHub Repository Info CLI")
	fmt.Println("Введите имя репозитория в формате owner/repo или два аргумента owner repo")
	fmt.Print("> ")
	input, err := reader.ReadString('\n')
	if err != nil {
		exitWithError(fmt.Errorf("ошибка чтения ввода: %w", err))
	}
	
	input = strings.TrimSpace(input)
	args := strings.Fields(input)
	owner, repo, err := parseRepositoryInput(Map(args, func(s string) string { return strings.ToLower(s) }))
	if err != nil {
		exitWithError(err)
	}

	info, err := fetchRepositoryInfo(owner, repo)
	if err != nil {
		exitWithError(err)
	}

	printRepositoryInfo(info)
}

func parseRepositoryInput(args []string) (string, string, error) {
	switch len(args) {
	case 1:
		parts := strings.Split(args[0], "/")
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return "", "", errors.New("некорректный ввод: используйте формат owner/repo или два аргумента owner repo")
		}
		return parts[0], parts[1], nil
	case 2:
		owner := strings.TrimSpace(args[0])
		repo := strings.TrimSpace(args[1])
		if owner == "" || repo == "" {
			return "", "", errors.New("некорректный ввод: owner и repo не должны быть пустыми")
		}
		return owner, repo, nil
	default:
		return "", "", errors.New("некорректный ввод")
	}
}

func fetchRepositoryInfo(owner, repo string) (*repositoryInfo, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	url := fmt.Sprintf(githubAPIURLTemplate, owner, repo)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания HTTP-запроса: %w", err)
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "github-repo-cli")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("сетевая ошибка: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.New("репозиторий не найден")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API вернул неожиданный статус: %s", resp.Status)
	}

	var info repositoryInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, fmt.Errorf("ошибка чтения JSON-ответа: %w", err)
	}

	return &info, nil
}

func printRepositoryInfo(info *repositoryInfo) {
	description := info.Description
	if strings.TrimSpace(description) == "" {
		description = "(нет описания)"
	}

	createdAt := info.CreatedAt
	if parsed, err := time.Parse(time.RFC3339, info.CreatedAt); err == nil {
		createdAt = parsed.Format("2006-01-02 15:04:05 MST")
	}
	
	fmt.Println("Информация о репозитории:")
	fmt.Printf("| Имя: %s\n", info.Name)
	fmt.Printf("| Описание: %s\n", description)
	fmt.Printf("| Звёзды: %d\n", info.StargazersCount)
	fmt.Printf("| Форки: %d\n", info.ForksCount)
	fmt.Printf("| Дата создания: %s\n", createdAt)
}

func exitWithError(err error) {
	fmt.Fprintf(os.Stderr, "Ошибка: %v\n", err)
	os.Exit(1)
}
