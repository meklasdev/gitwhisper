<div align="center">

# 🤫 GitWhisper

<img src="https://readme-typing-svg.demolab.com?font=Fira+Code&size=22&pause=1000&color=7AA2F7&center=true&vCenter=true&width=600&lines=AI-Powered+Commit+Messages;Privacy-First+%7C+Multi-Provider;Stop+Writing+Commits+Manually;Built+with+Go+%E2%9A%A1" alt="Typing SVG" />

**The CLI that whispers perfect commit messages using AI**

[![License](https://img.shields.io/badge/license-MIT-7AA2F7?style=for-the-badge&labelColor=1a1b26)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-1.21+-00ADD8?style=for-the-badge&logo=go&labelColor=1a1b26)](https://go.dev)
[![Stars](https://img.shields.io/github/stars/meklas/gitwhisper?style=for-the-badge&color=bb9af7&labelColor=1a1b26)](https://github.com/meklas/gitwhisper/stargazers)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-7dcfff?style=for-the-badge&labelColor=1a1b26)](CONTRIBUTING.md)

</div>

---

## ✨ Why GitWhisper?

<table>
<tr>
<td width="50%">

### 🔒 **Privacy First**
Run **Ollama** locally. Your code never leaves your machine. Perfect for corporate environments.

### 🌐 **Universal**
Works with **OpenAI**, **Gemini**, **Ollama**, and any **OpenAI-compatible** API (Groq, DeepSeek, etc.).

</td>
<td width="50%">

### 🎨 **Beautiful UX**
Modern terminal UI with colors, spinners, and smooth interactions powered by Charm libraries.

### ⚡ **Lightning Fast**
Written in Go. Single binary. Zero dependencies. Instant startup.

</td>
</tr>
</table>

---

## 🚀 Quick Start

```bash
# Clone the repo
git clone https://github.com/meklas/gitwhisper.git
cd gitwhisper

# Build (requires Go 1.21+)
go build -o gitwhisper .

# Stage your changes
git add .

# Let AI write your commit
./gitwhisper generate
```

---

## 🛠️ Configuration

Create `~/.gitwhisper.yaml`:

```yaml
ai:
  provider: ollama  # openai | ollama | gemini | openai-compatible

  # Ollama (Local & Private)
  ollama_endpoint: "http://localhost:11434/api/generate"
  ollama_model: "mistral"  # or qwen, llama3, etc.

  # OpenAI
  openai_api_key: "sk-..."
  openai_model: "gpt-4o"

  # Google Gemini
  gemini_api_key: "AIza..."
  gemini_model: "gemini-pro"

  # Generic (DeepSeek, Groq, etc.)
  # openai_base_url: "https://api.deepseek.com/v1/chat/completions"
```

---

## � Features

<div align="center">

| Feature | Description |
|---------|-------------|
| 🔐 **Privacy** | Local models via Ollama - no data leaves your machine |
| 🤖 **Multi-AI** | OpenAI, Gemini, Ollama, or any compatible provider |
| 📝 **Conventional Commits** | Generates `feat:`, `fix:`, `docs:` format automatically |
| ✏️ **Interactive** | Review, edit, or reject before committing |
| 🎨 **Beautiful CLI** | Lipgloss styling + Bubbletea TUI |
| ⚡ **Fast** | Go binary - instant startup, no runtime needed |

</div>

---

## 🧠 Supported AI Providers

<div align="center">

| Provider | Privacy | Speed | Quality | Models |
|----------|---------|-------|---------|--------|
| **Ollama** | 🟢 Local | ⚡ Fast | ⭐⭐⭐⭐ | `mistral`, `llama3`, `qwen` |
| **OpenAI** | 🔴 Cloud | ⚡⚡ Very Fast | ⭐⭐⭐⭐⭐ | `gpt-4o`, `gpt-3.5-turbo` |
| **Gemini** | 🔴 Cloud | ⚡⚡⚡ Ultra Fast | ⭐⭐⭐⭐ | `gemini-pro` |
| **Compatible** | 🟡 Varies | 🟡 Varies | 🟡 Varies | Groq, DeepSeek, LocalAI |

</div>

---

## 📸 Demo

```bash
$ git add .
$ gitwhisper generate

⠋ Generating commit message...

╭─────────────────────────────────────────────────╮
│                                                 │
│  feat: add multi-provider AI engine support    │
│                                                 │
╰─────────────────────────────────────────────────╯

Commit with this message? (y/n/e[dit]): y
✓ Commit successful!
```

---

## 🏗️ Architecture

Built with industry-standard Go libraries:

- **[Cobra](https://github.com/spf13/cobra)** - CLI framework
- **[Viper](https://github.com/spf13/viper)** - Configuration management
- **[Lipgloss](https://github.com/charmbracelet/lipgloss)** - Terminal styling
- **[Bubbletea](https://github.com/charmbracelet/bubbletea)** - TUI framework

---

## 🤝 Contributing

Contributions are welcome! See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

1. Fork the repo
2. Create your feature branch (`git checkout -b feature/amazing`)
3. Commit using GitWhisper 😉 (`gitwhisper generate`)
4. Push and open a PR

---

## � License

MIT © [meklas](https://github.com/meklas)

---

<div align="center">

**[⭐ Star this repo](https://github.com/meklas/gitwhisper)** if you find it useful!

Made with ❤️ and Go

</div>
