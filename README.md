# 🌐 LAN Share

> Share files instantly over your local network with QR codes - no installation, no cloud, just direct transfers

## 📦 Installation

> **Note:** macOS support is currently in development and not yet functional.

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

### Share a single file (with path)

```bash
lanshare share document.pdf
```

### Share multiple files

```bash
lanshare share document.pdf image.png video.mp4
```

### Share an entire folder

```bash
lanshare share /path/to/folder
```

### Share mixed files and folders

```bash
lanshare share file1.txt folder1/ file2.pdf folder2/
```

### Share using the interactive picker

```bash
lanshare share
```

Navigate folders with arrow keys and press Enter to select!
- When you select a folder, you can choose to share it or navigate into it
- Use the "Share this entire folder" option to share your current directory
- Select individual files to share them

**Note:** When sharing multiple files or folders, they are automatically bundled into a convenient ZIP archive.

```

## 🌟 How It Works

1. Run `lanshare share yourfile.pdf` (or specify multiple files/folders)
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
