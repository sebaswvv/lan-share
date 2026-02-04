# 🌐 LAN Share

> Share files instantly over your local network with QR codes - no installation, no cloud, just direct transfers

## 📦 Installation (No Go Required!)

### Windows

1. Go to [Releases](https://github.com/yourusername/lan-share/releases/latest)
2. Download `lanshare_Windows_x86_64.zip`
3. Extract the zip file
4. Double-click `lanshare.exe` or run from PowerShell/CMD

```powershell
# After extracting, run:
.\lanshare.exe share myfile.txt
```

### macOS

1. Go to [Releases](https://github.com/yourusername/lan-share/releases/latest)
2. Download `lanshare_Darwin_x86_64.tar.gz` (Intel) or `lanshare_Darwin_arm64.tar.gz` (Apple Silicon)
3. Extract and move to PATH:

```bash
tar -xzf lanshare_Darwin_*.tar.gz
sudo mv lanshare /usr/local/bin/
```

### Linux

1. Go to [Releases](https://github.com/yourusername/lan-share/releases/latest)
2. Download `lanshare_Linux_x86_64.tar.gz`
3. Extract and install:

```bash
tar -xzf lanshare_Linux_x86_64.tar.gz
sudo mv lanshare /usr/local/bin/
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
