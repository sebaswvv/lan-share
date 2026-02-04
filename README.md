# 🌐 LAN Share

> Share files instantly over your local network with QR codes - no installation, no cloud, just direct transfers

## 📦 Installation

### macOS / Linux (Homebrew)

```bash
brew tap sebaswvv/tap
brew install lanshare
```

### Windows (Scoop)

```powershell
scoop bucket add lanshare https://github.com/sebaswvv/scoop-bucket
scoop install lanshare
```

### Using Go

```bash
go install github.com/sebaswvv/lan-share@latest
```

## 🚀 Quick Start

### Share a file (with path)

```bash
lanshare share document.pdf
```

### Share a file (with file picker)

```bash
lanshare share
```

Navigate folders with arrow keys and press Enter to select!

```

## 🌟 How It Works

1. Run `lanshare share yourfile.pdf`
2. A local web server starts on your computer
3. QR code + URL are displayed
4. Share the URL or scan the QR code with your phone
5. Click the download button on the beautiful web page
6. File downloads directly from your computer!

**No uploads to cloud services. No third-party servers. Just direct peer-to-peer on your LAN.**

## 🤝 Contributing

Contributions welcome! Open issues or submit PRs.

## 📄 License

MIT License - see [LICENSE](LICENSE)

## 👤 Author

**Sebastiaan van Vliet**

---

⭐ **Found this useful? Give it a star!**
